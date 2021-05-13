package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	pkg_logger "tiniapp-backend-oauth-sample/pkg/logger"
	pkg_signature "tiniapp-backend-oauth-sample/pkg/signature"
)

type refreshAuthAccessTokenResult struct {
	Data *struct {
		AccessToken  string   `json:"access_token"`
		RefreshToken string   `json:"refresh_token"`
		ExpiresIn    int64    `json:"expires_in"`
		Scopes       []string `json:"scopes"`
	} `json:"data"`
	Error *struct {
		Code    int    `json:"code"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
	} `json:"error"`
}

type getAccessTokenFromRefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *Service) GetAccessTokenFromRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := pkg_logger.GetLogger().WithContext(ctx).WithPrefix("GetAccessTokenFromRefreshToken")

		var input getAccessTokenFromRefreshTokenInput
		err := c.ShouldBindJSON(&input)
		if err != nil {
			logger.Errorf("could not parse request data, error: %+v", err)
			c.JSON(
				http.StatusOK,
				map[string]interface{}{
					"error": map[string]interface{}{
						"status":  http.StatusBadRequest,
						"reason":  "bad_request",
						"message": "Bad Request",
					},
				},
			)
			return
		}

		timestamp := time.Now().UnixNano() / int64(time.Millisecond)
		data := map[string]interface{}{
			"refresh_token": input.RefreshToken,
			"client_id":     s.Config.ClientID,
			"timestamp":     timestamp,
		}

		signature, err := pkg_signature.SignWithData(s.Config.ClientSecret, data)
		if err != nil {
			logger.Errorf("could not calculate signature, error: %+v", err)
			c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"error": map[string]interface{}{
						"status":  http.StatusInternalServerError,
						"reason":  "internal_server_error",
						"message": err.Error(),
					},
				},
			)
			return
		}

		data["signature"] = signature
		logger.Debug(data)

		result := refreshAuthAccessTokenResult{}
		resp, err := s.HTTPClient.R().
			SetHeader("content-type", "application/json").
			SetHeader("x-request-id", uuid.NewString()).
			SetBody(data).
			SetResult(&result).
			Post(fmt.Sprintf("%s%s", s.Config.TiniAppServerAddress, "/api/v1/oauth/auth/token/refresh"))

		logger.Debug(resp, err)

		if err != nil {
			logger.Errorf("could not get response, error: %+v", err)
			c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"error": map[string]interface{}{
						"status":  http.StatusInternalServerError,
						"reason":  "internal_server_error",
						"message": err.Error(),
					},
				},
			)
			return
		}

		if result.Data == nil || result.Error != nil {
			logger.Errorf("got failure response, result: %+v", result)
			c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"error": map[string]interface{}{
						"status":  http.StatusBadRequest,
						"reason":  "bad_request",
						"message": "Bad Request",
					},
				},
			)
			return
		}
		logger.Infof("got success response, result: %+v", result)

		c.JSON(
			http.StatusOK,
			map[string]interface{}{
				"data": result.Data,
			},
		)
	}
}

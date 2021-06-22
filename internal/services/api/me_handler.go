package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	pkg_logger "tiniapp-backend-oauth-sample/pkg/logger"
	pkg_signature "tiniapp-backend-oauth-sample/pkg/signature"
)

type getMeResult struct {
	Data *struct {
		CustomerID   int64    `json:"customer_id"`
		CustomerName string   `json:"customer_name"`
		Scopes       []string `json:"scopes"`
	} `json:"data"`
	Error *struct {
		Code    int    `json:"code"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
	} `json:"error"`
}

type getMeFromAccessTokenInput struct {
	AccessToken string `json:"access_token" binding:"required"`
}

func (s *Service) GetMeFromAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := pkg_logger.GetLogger().WithContext(ctx).WithPrefix("GetMeFromAccessToken")

		var input getMeFromAccessTokenInput
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
			"access_token": input.AccessToken,
		}

		signingPayload, _ := pkg_signature.PreparePayload(timestamp, s.Config.ClientID, data)
		signature, err := pkg_signature.SignWithPayload(s.Config.ClientSecret, signingPayload)
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

		logger.Debug(data)

		result := getMeResult{}
		resp, err := s.HTTPClient.R().
			SetHeader("content-type", "application/json").
			SetHeader("x-request-id", uuid.NewString()).
			SetHeader("X-Tiniapp-Timestamp", strconv.FormatInt(timestamp, 10)).
			SetHeader("X-Tiniapp-Client-Id", s.Config.ClientID).
			SetHeader("X-Tiniapp-Signature", signature).
			SetBody(data).
			SetResult(&result).
			Post(fmt.Sprintf("%s%s", s.Config.TiniAppServerAddress, "/api/v1/oauth/me"))

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

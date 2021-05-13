package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	pkg_logger "tiniapp-backend-oauth-sample/pkg/logger"
)

func (s *Service) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := pkg_logger.GetLogger().WithContext(ctx).WithPrefix("Ping")

		data := map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"build":     os.Getenv("BUILD"),
		}

		logger.Infof("returning: %+v", data)
		c.JSON(
			http.StatusOK,
			data,
		)
	}
}

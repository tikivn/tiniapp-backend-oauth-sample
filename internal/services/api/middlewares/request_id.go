package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	pkg_ctxutil "tiniapp-backend-oauth-sample/pkg/ctxutil"
)

func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("x-request-id")
		if requestID == "" {
			requestID = uuid.NewString()
			c.Request.Header.Set("x-request-id", requestID)
			c.Writer.Header().Set("x-request-id", requestID)
		}

		ctx := c.Request.Context()
		newCtx := pkg_ctxutil.WithRequestID(ctx, requestID)
		c.Request = c.Request.WithContext(newCtx)

		c.Set("x-request-id", requestID)

		c.Next()
	}
}

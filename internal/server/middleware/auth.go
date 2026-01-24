package middleware

import (
	"net/http"
	"strings"

	"github.com/fengyuanluo/go-backend-template/internal/server/auth"
	"github.com/fengyuanluo/go-backend-template/internal/server/resp"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(c.GetHeader("Authorization"))
		if token == "" {
			resp.Error(c, http.StatusUnauthorized, resp.ErrUnauthorized)
			c.Abort()
			return
		}
		parts := strings.Fields(token)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			resp.Error(c, http.StatusUnauthorized, resp.ErrUnauthorized)
			c.Abort()
			return
		}
		if !auth.VerifyToken(parts[1]) {
			resp.Error(c, http.StatusUnauthorized, resp.ErrUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}

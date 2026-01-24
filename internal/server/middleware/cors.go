package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = false
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	if origins := parseAllowOrigins(os.Getenv("EXAMPLE_CORS_ALLOW_ORIGINS")); len(origins) > 0 {
		config.AllowAllOrigins = false
		config.AllowOrigins = origins
		config.AllowCredentials = true
	}
	return cors.New(config)
}

func parseAllowOrigins(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, item := range parts {
		origin := strings.TrimSpace(item)
		if origin == "" {
			continue
		}
		out = append(out, origin)
	}
	return out
}

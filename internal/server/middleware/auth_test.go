package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fengyuanluo/go-backend-template/internal/conf"
	"github.com/fengyuanluo/go-backend-template/internal/server/auth"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddlewareRequiresBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Auth())
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc")
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid scheme, got %d", recorder.Code)
	}
}

func TestAuthMiddlewareAllowsValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	originalSecret := conf.AppConfig.Auth.Secret
	t.Cleanup(func() {
		conf.AppConfig.Auth.Secret = originalSecret
	})
	conf.AppConfig.Auth.Secret = "secret"

	token, _, err := auth.GenerateToken(1)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	router := gin.New()
	router.Use(Auth())
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

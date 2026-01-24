package auth

import (
	"testing"

	"github.com/fengyuanluo/go-backend-template/internal/conf"
	"github.com/golang-jwt/jwt/v5"
)

func TestVerifyTokenRejectsWrongAlg(t *testing.T) {
	origProvider := secretProvider
	t.Cleanup(func() {
		secretProvider = origProvider
	})
	secretProvider = func() (string, error) { return "secret", nil }
	conf.AppConfig.Auth.Secret = "secret"

	claims := jwt.RegisteredClaims{}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte("secret"))
	if VerifyToken(token) {
		t.Fatalf("expected reject HS512")
	}
}

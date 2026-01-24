package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/fengyuanluo/go-backend-template/internal/conf"
	"github.com/golang-jwt/jwt/v5"
)

var secretProvider = func() (string, error) {
	secret := strings.TrimSpace(conf.AppConfig.Auth.Secret)
	if secret == "" {
		return "", errors.New("auth secret is empty")
	}
	return secret, nil
}

func GenerateToken(expiresMin int) (string, string, error) {
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    conf.APP_NAME,
	}
	if expiresMin == 0 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Duration(15) * time.Minute))
	} else if expiresMin > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Duration(expiresMin) * time.Minute))
	} else if expiresMin == -1 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Duration(30) * 24 * time.Hour))
	}
	secret, err := secretProvider()
	if err != nil {
		return "", "", err
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return token, claims.ExpiresAt.Format(time.RFC3339), nil
}

func VerifyToken(token string) bool {
	secret, err := secretProvider()
	if err != nil {
		return false
	}
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !jwtToken.Valid {
		return false
	}
	return true
}

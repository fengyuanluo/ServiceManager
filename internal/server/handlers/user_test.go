package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fengyuanluo/go-backend-template/internal/op"
	"github.com/gin-gonic/gin"
)

func TestChangePasswordInvalidOldPassword(t *testing.T) {
	original := userChangePassword
	t.Cleanup(func() {
		userChangePassword = original
	})
	userChangePassword = func(_, _ string) error {
		return fmt.Errorf("bad: %w", op.ErrInvalidOldPassword)
	}

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/change-password", strings.NewReader(`{"old_password":"old","new_password":"new"}`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	changePassword(c)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestChangeUsernameInvalidUsername(t *testing.T) {
	original := userChangeUsername
	t.Cleanup(func() {
		userChangeUsername = original
	})
	userChangeUsername = func(string) error {
		return fmt.Errorf("bad: %w", op.ErrInvalidUsername)
	}

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/change-username", strings.NewReader(`{"new_username":"admin"}`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	changeUsername(c)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

package cmd

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRouterPropagatesError(t *testing.T) {
	prev := registerRoutes
	registerRoutes = func(*gin.Engine) error { return errors.New("boom") }
	defer func() { registerRoutes = prev }()

	if _, err := setupRouter(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestInitializeAppStopsOnDBError(t *testing.T) {
	prevDB := initDB
	prevUser := initUser
	calledUser := false
	initDB = func() error { return errors.New("db") }
	initUser = func() error { calledUser = true; return nil }
	defer func() {
		initDB = prevDB
		initUser = prevUser
	}()

	if err := initializeApp(); err == nil {
		t.Fatalf("expected error")
	}
	if calledUser {
		t.Fatalf("initUser should not run on db error")
	}
}

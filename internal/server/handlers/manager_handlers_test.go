package handlers

import (
	"testing"

	"github.com/fengyuanluo/go-backend-template/internal/server/router"
)

func TestManagerHandlersListRoute(t *testing.T) {
	count := router.GetRouterCount()
	if count < 30 {
		t.Fatalf("expected manager routes, got %d", count)
	}
}

package handlers

import (
	"testing"

	"github.com/fengyuanluo/go-backend-template/internal/server/router"
)

func TestManagerRoutesRegistered(t *testing.T) {
	count := router.GetRouterCount()
	if count < 8 {
		t.Fatalf("expected manager routes, got %d", count)
	}
}

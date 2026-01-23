package db

import (
	"path/filepath"
	"testing"

	"github.com/fengyuanluo/go-backend-template/internal/conf"
)

func TestInitDBWithPureGoSQLite(t *testing.T) {
	tmp := t.TempDir()
	conf.AppConfig.Database.Type = "sqlite"
	conf.AppConfig.Database.Path = filepath.Join(tmp, "data.db")

	if err := InitDB(); err != nil {
		t.Fatalf("expected InitDB success with CGO disabled: %v", err)
	}
	t.Cleanup(func() {
		if err := Close(); err != nil {
			t.Fatalf("expected Close success: %v", err)
		}
	})
}

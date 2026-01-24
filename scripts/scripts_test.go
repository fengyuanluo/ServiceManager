package scripts

import (
	"os"
	"strings"
	"testing"
)

func TestBuildScriptFlags(t *testing.T) {
	data, err := os.ReadFile("build.sh")
	if err != nil {
		t.Fatalf("read build.sh: %v", err)
	}
	text := string(data)
	if !strings.Contains(text, "github.com/fengyuanluo/go-backend-template/internal/conf.Version") {
		t.Fatalf("missing correct ldflags path")
	}
	if !strings.Contains(text, `APP_NAME="service-manager"`) {
		t.Fatalf("APP_NAME not aligned")
	}
}

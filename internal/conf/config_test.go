package conf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadLoggingLevel(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.json")
	payload := `{"logging":{"level":"debug"},"server":{"host":"127.0.0.1","port":9001},"database":{"type":"sqlite","path":"` + filepath.Join(tmp, "data.db") + `"}}`
	if err := os.WriteFile(path, []byte(payload), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	viper.Reset()
	defer viper.Reset()

	if err := Load(path); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if AppConfig.Log.Level != "debug" {
		t.Fatalf("expected log level debug, got %q", AppConfig.Log.Level)
	}
}

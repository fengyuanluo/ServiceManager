package log

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestParseLevel(t *testing.T) {
	if got := parseLevel("debug"); got != zapcore.DebugLevel {
		t.Fatalf("debug -> %v", got)
	}
	if got := parseLevel("WARN"); got != zapcore.WarnLevel {
		t.Fatalf("WARN -> %v", got)
	}
	if got := parseLevel(""); got != zapcore.InfoLevel {
		t.Fatalf("empty -> %v", got)
	}
}

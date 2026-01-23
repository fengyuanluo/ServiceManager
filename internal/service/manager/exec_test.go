package manager

import (
	"errors"
	"testing"
)

func TestRunCommandError(t *testing.T) {
	_, err := RunCommand(ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		return nil, errors.New("boom")
	}}, "x")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolveCommandFallback(t *testing.T) {
	got, err := resolveCommand("systemctl", func(name string) (string, error) {
		return "", errors.New("not found")
	}, []string{"/bin/systemctl"})
	if err != nil || got != "/bin/systemctl" {
		t.Fatalf("expected fallback path")
	}
}

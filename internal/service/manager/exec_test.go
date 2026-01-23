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

package op

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitialCredentialsFromEnv(t *testing.T) {
	t.Setenv("EXAMPLE_INIT_USER", "rooter")
	t.Setenv("EXAMPLE_INIT_PASSWORD", "secret")
	initialCredsPath = filepath.Join(t.TempDir(), "init.json")
	randomPassword = func() (string, error) { return "rand", nil }

	user, pass, wrote, err := initialCredentials()
	if err != nil || wrote {
		t.Fatalf("unexpected result: %v", err)
	}
	if user != "rooter" || pass != "secret" {
		t.Fatalf("unexpected creds")
	}
}

func TestInitialCredentialsWritesFile(t *testing.T) {
	os.Unsetenv("EXAMPLE_INIT_USER")
	os.Unsetenv("EXAMPLE_INIT_PASSWORD")
	tmp := t.TempDir()
	initialCredsPath = filepath.Join(tmp, "init.json")
	randomPassword = func() (string, error) { return "randpass", nil }

	user, pass, wrote, err := initialCredentials()
	if err != nil || !wrote {
		t.Fatalf("expected file write: %v", err)
	}
	if user == "" || pass != "randpass" {
		t.Fatalf("unexpected creds")
	}
	if _, err := os.Stat(initialCredsPath); err != nil {
		t.Fatalf("expected creds file")
	}
}

package manager

import (
	"os"
	"testing"
)

func TestDetectEnvironment(t *testing.T) {
	src := EnvSource{
		ReadFile: func(path string) ([]byte, error) {
			switch path {
			case "/etc/os-release":
				return []byte("ID=ubuntu\n"), nil
			case "/proc/1/comm":
				return []byte("systemd\n"), nil
			default:
				return nil, os.ErrNotExist
			}
		},
		Exec: func(name string, args ...string) ([]byte, error) {
			return []byte("none\n"), nil
		},
	}
	info := DetectEnvironment(src)
	if info.Distro != "ubuntu" || info.Init != "systemd" {
		t.Fatalf("unexpected env: %+v", info)
	}
}

package manager

import (
	"context"
	"os"
	"os/exec"
)

func RunCommand(src ExecSource, name string, args ...string) ([]byte, error) {
	if src.Run != nil {
		return src.Run(name, args...)
	}
	resolved, err := resolveCommand(name, exec.LookPath, commandFallbacks(name))
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(context.Background(), resolved, args...)
	return cmd.CombinedOutput()
}

func resolveCommand(name string, look func(string) (string, error), fallbacks []string) (string, error) {
	if path, err := look(name); err == nil {
		return path, nil
	}
	for _, candidate := range fallbacks {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", errCommandNotFound
}

func commandFallbacks(name string) []string {
	switch name {
	case "systemctl":
		return []string{"/bin/systemctl", "/usr/bin/systemctl"}
	case "journalctl":
		return []string{"/bin/journalctl", "/usr/bin/journalctl"}
	default:
		return nil
	}
}

package manager

import (
	"context"
	"os/exec"
)

func RunCommand(src ExecSource, name string, args ...string) ([]byte, error) {
	if src.Run != nil {
		return src.Run(name, args...)
	}
	cmd := exec.CommandContext(context.Background(), name, args...)
	return cmd.CombinedOutput()
}

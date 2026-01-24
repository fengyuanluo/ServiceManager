package manager

import "strings"

func DetectManagers(env EnvInfo) []string {
	out := []string{}
	initName := strings.ToLower(strings.TrimSpace(env.Init))
	switch initName {
	case "systemd":
		out = append(out, "systemd")
	case "openrc":
		out = append(out, "openrc")
	case "init", "sysvinit":
		out = append(out, "sysv")
	}
	return out
}

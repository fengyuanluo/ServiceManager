package manager

import "strings"

func DetectEnvironment(src EnvSource) EnvInfo {
	distro := "unknown"
	init := "unknown"
	if src.ReadFile != nil {
		if data, err := src.ReadFile("/etc/os-release"); err == nil {
			distro = parseID(string(data))
		}
		if data, err := src.ReadFile("/proc/1/comm"); err == nil {
			init = strings.TrimSpace(string(data))
		}
	}
	return EnvInfo{Distro: distro, Init: init}
}

func parseID(content string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "ID=") {
			return strings.Trim(strings.TrimPrefix(line, "ID="), "\"'")
		}
	}
	return "unknown"
}

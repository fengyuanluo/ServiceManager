package manager

import "strings"

func NormalizeStatus(raw string) ServiceStatus {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" {
		return ServiceStatus("unknown")
	}
	switch {
	case strings.Contains(normalized, "inactive"):
		return ServiceStatus("stopped")
	case strings.Contains(normalized, "running"),
		strings.Contains(normalized, "active"),
		strings.Contains(normalized, "started"),
		strings.Contains(normalized, "online"):
		return ServiceStatus("running")
	case strings.Contains(normalized, "stopped"),
		strings.Contains(normalized, "dead"),
		strings.Contains(normalized, "inactive"),
		strings.Contains(normalized, "exited"),
		strings.Contains(normalized, "stop"):
		return ServiceStatus("stopped")
	case strings.Contains(normalized, "failed"),
		strings.Contains(normalized, "fatal"),
		strings.Contains(normalized, "error"):
		return ServiceStatus("failed")
	default:
		return ServiceStatus("unknown")
	}
}

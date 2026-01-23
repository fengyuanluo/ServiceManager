package manager

import "strings"

func ParseSupervisorList(out string) ([]ServiceSummary, error) {
	line := strings.TrimSpace(out)
	if line == "" {
		return []ServiceSummary{}, nil
	}
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return []ServiceSummary{}, nil
	}
	name := parts[0]
	status := ServiceStatus("unknown")
	if len(parts) > 1 {
		status = ServiceStatus(strings.ToLower(parts[1]))
	}
	return []ServiceSummary{{Name: name, Status: status}}, nil
}

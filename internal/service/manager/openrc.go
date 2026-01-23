package manager

import "strings"

func ParseOpenRCList(out string) ([]ServiceSummary, error) {
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
	if strings.Contains(line, "started") {
		status = ServiceStatus("started")
	}
	return []ServiceSummary{{Name: name, Status: status}}, nil
}

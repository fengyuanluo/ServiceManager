package manager

import "strings"

func ParseOpenRCList(out string) ([]ServiceSummary, error) {
	lines := strings.Split(out, "\n")
	items := make([]ServiceSummary, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		name := parts[0]
		status := ServiceStatus("unknown")
		if len(parts) > 2 {
			status = NormalizeStatus(parts[2])
		} else if strings.Contains(line, "started") {
			status = NormalizeStatus("started")
		}
		items = append(items, ServiceSummary{Name: name, Status: status})
	}
	return items, nil
}

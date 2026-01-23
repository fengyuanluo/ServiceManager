package manager

import "strings"

func ParseSupervisorList(out string) ([]ServiceSummary, error) {
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
		if len(parts) > 1 {
			status = NormalizeStatus(parts[1])
		}
		detail := ""
		if len(parts) > 2 {
			detail = strings.Join(parts[2:], " ")
		}
		items = append(items, ServiceSummary{Name: name, Status: status, Detail: detail})
	}
	return items, nil
}

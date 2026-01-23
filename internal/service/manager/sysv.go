package manager

import "strings"

func ParseSysVList(out string) ([]ServiceSummary, error) {
	lines := strings.Split(out, "\n")
	items := make([]ServiceSummary, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if len(fields) >= 4 && fields[0] == "[" {
			indicator := fields[1]
			name := fields[len(fields)-1]
			status := ServiceStatus("unknown")
			switch indicator {
			case "+":
				status = NormalizeStatus("running")
			case "-":
				status = NormalizeStatus("stopped")
			default:
				status = NormalizeStatus("unknown")
			}
			items = append(items, ServiceSummary{Name: name, Status: status})
			continue
		}
		name := fields[0]
		status := ServiceStatus("unknown")
		if len(fields) > 1 {
			status = NormalizeStatus(fields[1])
		}
		items = append(items, ServiceSummary{Name: name, Status: status})
	}
	return items, nil
}

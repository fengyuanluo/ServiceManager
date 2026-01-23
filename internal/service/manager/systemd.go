package manager

import "strings"

func ParseSystemdList(out string) ([]ServiceSummary, error) {
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
		name := strings.TrimSuffix(fields[0], ".service")
		status := ServiceStatus("unknown")
		if len(fields) > 3 {
			status = NormalizeStatus(fields[3])
		} else if len(fields) > 2 {
			status = NormalizeStatus(fields[2])
		}
		detail := ""
		if len(fields) > 4 {
			detail = strings.Join(fields[4:], " ")
		}
		items = append(items, ServiceSummary{Name: name, Status: status, Detail: detail})
	}
	return items, nil
}

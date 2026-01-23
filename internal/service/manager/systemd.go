package manager

import "strings"

func ParseSystemdList(out string) ([]ServiceSummary, error) {
	line := strings.TrimSpace(out)
	if line == "" {
		return []ServiceSummary{}, nil
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return []ServiceSummary{}, nil
	}
	name := strings.TrimSuffix(fields[0], ".service")
	status := ServiceStatus("unknown")
	if len(fields) > 2 {
		status = ServiceStatus(fields[2])
	}
	detail := ""
	if len(fields) > 3 {
		detail = strings.Join(fields[3:], " ")
	}
	return []ServiceSummary{{Name: name, Status: status, Detail: detail}}, nil
}

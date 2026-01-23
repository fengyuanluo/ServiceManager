package manager

import "strings"

func ParseSysVList(out string) ([]ServiceSummary, error) {
	line := strings.TrimSpace(out)
	if line == "" {
		return []ServiceSummary{}, nil
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return []ServiceSummary{}, nil
	}
	name := fields[0]
	status := ServiceStatus("unknown")
	if len(fields) > 1 {
		status = ServiceStatus(fields[1])
	}
	return []ServiceSummary{{Name: name, Status: status}}, nil
}

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

func ParseSysVDetail(out string, fallbackName string) ServiceDetail {
	detail := ServiceDetail{Name: fallbackName, Status: ServiceStatus("unknown")}
	lines := strings.Split(out, "\n")
	hasSystemdStyle := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "●") {
			hasSystemdStyle = true
			head := strings.TrimSpace(strings.TrimPrefix(line, "●"))
			parts := strings.SplitN(head, " - ", 2)
			if len(parts) > 0 {
				name := strings.TrimSuffix(strings.TrimSpace(parts[0]), ".service")
				if name != "" {
					detail.Name = name
				}
			}
			if len(parts) > 1 {
				detail.Description = strings.TrimSpace(parts[1])
			}
			continue
		}
		if strings.HasPrefix(line, "Active:") {
			hasSystemdStyle = true
			detail.Status = NormalizeStatus(line)
			if idx := strings.Index(line, " since "); idx != -1 {
				since := strings.TrimSpace(line[idx+7:])
				if cut := strings.Index(since, ";"); cut != -1 {
					since = strings.TrimSpace(since[:cut])
				}
				detail.Since = since
			}
			continue
		}
		if strings.HasPrefix(line, "Main PID:") {
			hasSystemdStyle = true
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				detail.PID = fields[2]
			}
			continue
		}
		if strings.HasPrefix(line, "Loaded:") {
			hasSystemdStyle = true
			if start := strings.Index(line, "("); start != -1 {
				rest := line[start+1:]
				end := strings.Index(rest, ";")
				if end == -1 {
					end = strings.Index(rest, ")")
				}
				if end != -1 {
					detail.UnitPath = strings.TrimSpace(rest[:end])
				}
			}
			continue
		}
	}
	if hasSystemdStyle {
		return detail
	}
	items, err := ParseSysVList(out)
	if err == nil && len(items) > 0 {
		return ServiceDetail{Name: items[0].Name, Status: items[0].Status, Description: items[0].Detail}
	}
	return detail
}

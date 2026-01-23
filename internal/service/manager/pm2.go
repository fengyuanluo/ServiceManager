package manager

import (
	"encoding/json"
	"strings"
)

type pm2Process struct {
	Name   string `json:"name"`
	PID    int    `json:"pid"`
	PM2Env struct {
		Status string `json:"status"`
	} `json:"pm2_env"`
}

func ParsePM2List(out string) ([]ServiceSummary, error) {
	line := strings.TrimSpace(out)
	if line == "" {
		return []ServiceSummary{}, nil
	}
	if strings.HasPrefix(line, "[") {
		var items []pm2Process
		if err := json.Unmarshal([]byte(line), &items); err == nil {
			results := make([]ServiceSummary, 0, len(items))
			for _, item := range items {
				results = append(results, ServiceSummary{
					Name:   item.Name,
					Status: NormalizeStatus(item.PM2Env.Status),
					Detail: "",
				})
			}
			return results, nil
		}
	}
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return []ServiceSummary{}, nil
	}
	name := parts[0]
	status := ServiceStatus("unknown")
	if len(parts) > 1 {
		status = NormalizeStatus(parts[1])
	}
	return []ServiceSummary{{Name: name, Status: status}}, nil
}

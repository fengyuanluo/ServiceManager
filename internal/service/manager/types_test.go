package manager

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestServiceSummaryJSON(t *testing.T) {
	data := ServiceSummary{Name: "ssh", Status: ServiceStatus("running"), Detail: "OpenSSH"}
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal failed")
	}
	text := string(raw)
	if !strings.Contains(text, "\"name\"") || strings.Contains(text, "\"Name\"") {
		t.Fatalf("unexpected json keys: %s", text)
	}
}

func TestServiceDetailJSON(t *testing.T) {
	data := ServiceDetail{Name: "ssh", Status: ServiceStatus("running"), PID: "1"}
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal failed")
	}
	text := string(raw)
	if !strings.Contains(text, "\"pid\"") || strings.Contains(text, "\"PID\"") {
		t.Fatalf("unexpected json keys: %s", text)
	}
}

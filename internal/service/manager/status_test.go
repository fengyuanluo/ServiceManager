package manager

import "testing"

func TestNormalizeStatus(t *testing.T) {
	cases := map[string]ServiceStatus{
		"running":  ServiceStatus("running"),
		"active":   ServiceStatus("running"),
		"started":  ServiceStatus("running"),
		"online":   ServiceStatus("running"),
		"stopped":  ServiceStatus("stopped"),
		"dead":     ServiceStatus("stopped"),
		"inactive": ServiceStatus("stopped"),
		"failed":   ServiceStatus("failed"),
		"fatal":    ServiceStatus("failed"),
		"unknown":  ServiceStatus("unknown"),
	}
	for input, expected := range cases {
		if got := NormalizeStatus(input); got != expected {
			t.Fatalf("unexpected status for %s: %s", input, got)
		}
	}
}

package manager

import "testing"

func TestParseSupervisorList(t *testing.T) {
	out := "sshd RUNNING"
	items, _ := ParseSupervisorList(out)
	if len(items) != 1 || items[0].Name != "sshd" {
		t.Fatalf("parse failed")
	}
}

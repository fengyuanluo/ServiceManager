package manager

import "testing"

func TestParseSupervisorList(t *testing.T) {
	out := "sshd RUNNING pid 1, uptime 0:00:01\ncron STOPPED"
	items, _ := ParseSupervisorList(out)
	if len(items) != 2 || items[0].Name != "sshd" || items[1].Name != "cron" {
		t.Fatalf("parse failed")
	}
}

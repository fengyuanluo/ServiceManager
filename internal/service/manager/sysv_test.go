package manager

import "testing"

func TestParseSysVList(t *testing.T) {
	out := "ssh start/running\ncron stop/waiting"
	items, _ := ParseSysVList(out)
	if len(items) != 2 || items[0].Name != "ssh" || items[1].Name != "cron" {
		t.Fatalf("parse failed")
	}
}

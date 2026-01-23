package manager

import "testing"

func TestParseSysVList(t *testing.T) {
	out := "ssh start/running"
	items, _ := ParseSysVList(out)
	if len(items) != 1 || items[0].Name != "ssh" {
		t.Fatalf("parse failed")
	}
}

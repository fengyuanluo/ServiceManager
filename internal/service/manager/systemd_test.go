package manager

import "testing"

func TestParseSystemdList(t *testing.T) {
	out := "ssh.service loaded active running OpenSSH"
	items, err := ParseSystemdList(out)
	if err != nil || len(items) != 1 || items[0].Name != "ssh" {
		t.Fatalf("parse failed")
	}
}

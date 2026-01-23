package manager

import "testing"

func TestSystemdListParseMulti(t *testing.T) {
	out := "ssh.service loaded active running OpenSSH\ncron.service loaded inactive dead"
	items, err := ParseSystemdList(out)
	if err != nil || len(items) != 2 {
		t.Fatalf("parse failed")
	}
	if items[0].Name != "ssh" || items[1].Name != "cron" {
		t.Fatalf("unexpected names")
	}
}

func TestParseSystemdDetailExecStart(t *testing.T) {
	out := "Id=cron.service\nExecStart=/usr/sbin/cron -f\nUnitFileState=enabled\nActiveState=active"
	detail := parseSystemdDetail(out, "cron")
	if detail.Extra["exec"] != "/usr/sbin/cron -f" || detail.Extra["autostart"] != "enabled" {
		t.Fatalf("missing exec/autostart")
	}
}

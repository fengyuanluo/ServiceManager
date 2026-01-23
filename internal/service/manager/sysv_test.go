package manager

import "testing"

func TestParseSysVList(t *testing.T) {
	out := "ssh start/running\ncron stop/waiting"
	items, _ := ParseSysVList(out)
	if len(items) != 2 || items[0].Name != "ssh" || items[1].Name != "cron" {
		t.Fatalf("parse failed")
	}
}

func TestParseSysVDetailSystemdStyle(t *testing.T) {
	out := "● cron.service - Regular background program\n Active: active (running) since Thu\n Main PID: 872 (cron)\n Loaded: loaded (/lib/systemd/system/cron.service; enabled)"
	detail := ParseSysVDetail(out, "cron")
	if detail.Status != ServiceStatus("running") || detail.PID != "872" {
		t.Fatalf("unexpected detail parse")
	}
}

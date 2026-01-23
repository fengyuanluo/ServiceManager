package manager

import (
	"os"
	"strings"
	"testing"
)

func TestCreateSystemdService(t *testing.T) {
	dir := t.TempDir()
	calls := []string{}
	src := ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		calls = append(calls, name+" "+strings.Join(args, " "))
		return []byte("ok"), nil
	}}
	fs := FileSource{
		WriteFile: os.WriteFile,
		MkdirAll:  os.MkdirAll,
		Chmod:     os.Chmod,
	}
	req := CreateRequest{
		Name:       "demo",
		Exec:       "/usr/bin/demo",
		TemplateID: "simple",
		Options:    CreateOptions{Start: true, Autostart: true},
	}
	err := CreateService("systemd", CreateSource{Exec: src, FS: fs, BaseDir: dir}, req)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	unitPath := dir + "/systemd/demo.service"
	content, err := os.ReadFile(unitPath)
	if err != nil {
		t.Fatalf("missing unit file: %v", err)
	}
	if !strings.Contains(string(content), "ExecStart=/usr/bin/demo") {
		t.Fatalf("expected ExecStart in unit")
	}
	if len(calls) < 3 {
		t.Fatalf("expected systemctl calls")
	}
	if calls[0] != "systemctl daemon-reload" {
		t.Fatalf("unexpected reload call: %s", calls[0])
	}
	if calls[1] != "systemctl enable demo" {
		t.Fatalf("unexpected enable call: %s", calls[1])
	}
	if calls[2] != "systemctl start demo" {
		t.Fatalf("unexpected start call: %s", calls[2])
	}
}

package manager

import (
	"errors"
	"strings"
	"testing"
)

func TestListServicesSystemd(t *testing.T) {
	out := "ssh.service loaded active running OpenSSH\ncron.service loaded inactive dead"
	var calledName string
	var calledArgs []string
	src := ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		calledName = name
		calledArgs = append([]string{}, args...)
		return []byte(out), nil
	}}
	items, err := ListServices("systemd", src, "")
	if err != nil || len(items) != 2 {
		t.Fatalf("unexpected list result")
	}
	if calledName != "systemctl" || len(calledArgs) == 0 || calledArgs[0] != "list-units" {
		t.Fatalf("unexpected command: %s %v", calledName, calledArgs)
	}
}

func TestGetServiceLogsFilter(t *testing.T) {
	out := "info line\nerror line\nwarn line"
	src := ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		return []byte(out), nil
	}}
	lines, err := GetServiceLogs("systemd", src, "ssh", LogQuery{Lines: 200, Keyword: "error"})
	if err != nil || len(lines) != 1 {
		t.Fatalf("unexpected logs result")
	}
	if !strings.Contains(lines[0], "error") {
		t.Fatalf("filter failed")
	}
}

func TestGetServiceDetailSystemd(t *testing.T) {
	out := "Id=ssh.service\nDescription=OpenSSH\nMainPID=123\nActiveState=active\nActiveEnterTimestamp=Mon\nFragmentPath=/lib/systemd/system/ssh.service\nUser=root"
	var calledName string
	var calledArgs []string
	src := ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		calledName = name
		calledArgs = append([]string{}, args...)
		return []byte(out), nil
	}}
	detail, err := GetServiceDetail("systemd", src, "ssh", "")
	if err != nil || detail.Name != "ssh" || detail.PID != "123" {
		t.Fatalf("unexpected detail result")
	}
	if calledName != "systemctl" || len(calledArgs) == 0 || calledArgs[0] != "show" {
		t.Fatalf("unexpected detail command: %s %v", calledName, calledArgs)
	}
}

func TestRunServiceActionSystemd(t *testing.T) {
	var calledName string
	var calledArgs []string
	src := ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		calledName = name
		calledArgs = append([]string{}, args...)
		return []byte("ok"), nil
	}}
	if err := RunServiceAction("systemd", src, ActionRequest{Name: "ssh", Action: "restart"}); err != nil {
		t.Fatalf("action failed")
	}
	if calledName != "systemctl" || len(calledArgs) < 2 || calledArgs[0] != "restart" {
		t.Fatalf("unexpected action command: %s %v", calledName, calledArgs)
	}
}

func TestRunServiceBatch(t *testing.T) {
	src := ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		if len(args) > 1 && args[1] == "fail" {
			return nil, errors.New("boom")
		}
		return []byte("ok"), nil
	}}
	items := []BatchItem{
		{Name: "ok", Action: BatchAction("restart")},
		{Name: "fail", Action: BatchAction("restart")},
	}
	results := RunServiceBatch("systemd", src, "", items)
	if len(results) != 2 || results[0].Success != true || results[1].Success != false {
		t.Fatalf("unexpected batch result")
	}
}

func TestInstallManagerApt(t *testing.T) {
	calls := []string{}
	src := ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		calls = append(calls, name+" "+strings.Join(args, " "))
		return []byte("ok"), nil
	}}
	res := InstallManager("supervisor", EnvInfo{Distro: "ubuntu"}, src)
	if !res.Success {
		t.Fatalf("install failed: %s", res.Message)
	}
	if len(calls) < 2 {
		t.Fatalf("missing install commands")
	}
	if calls[0] != "apt-get update" || !strings.Contains(calls[1], "apt-get install") {
		t.Fatalf("unexpected install commands: %v", calls)
	}
}

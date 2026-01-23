package manager

import (
	"reflect"
	"testing"
)

func TestBuildListCommand(t *testing.T) {
	cases := []struct {
		manager  string
		user     string
		wantName string
		wantArgs []string
	}{
		{manager: "systemd", wantName: "systemctl", wantArgs: []string{"list-units", "--type=service", "--no-legend", "--no-pager"}},
		{manager: "sysv", wantName: "service", wantArgs: []string{"--status-all"}},
		{manager: "openrc", wantName: "rc-service", wantArgs: []string{"--status-all"}},
		{manager: "supervisor", wantName: "supervisorctl", wantArgs: []string{"status"}},
		{manager: "pm2", wantName: "pm2", wantArgs: []string{"jlist"}},
		{manager: "pm2", user: "app", wantName: "su", wantArgs: []string{"-l", "app", "-c", "pm2 jlist"}},
	}
	for _, item := range cases {
		cmd, err := buildListCommand(item.manager, item.user)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", item.manager, err)
		}
		if cmd.Name != item.wantName {
			t.Fatalf("unexpected name for %s: %s", item.manager, cmd.Name)
		}
		if !reflect.DeepEqual(cmd.Args, item.wantArgs) {
			t.Fatalf("unexpected args for %s: %v", item.manager, cmd.Args)
		}
	}
	if _, err := buildListCommand("unknown", ""); err == nil {
		t.Fatalf("expected error for unknown manager")
	}
}

func TestBuildDetailCommand(t *testing.T) {
	cases := []struct {
		manager  string
		user     string
		wantName string
		wantArgs []string
	}{
		{manager: "systemd", wantName: "systemctl", wantArgs: []string{"show", "ssh", "--no-pager", "--property=Id,Description,MainPID,ActiveState,ActiveEnterTimestamp,FragmentPath,User"}},
		{manager: "sysv", wantName: "service", wantArgs: []string{"ssh", "status"}},
		{manager: "openrc", wantName: "rc-service", wantArgs: []string{"ssh", "status"}},
		{manager: "supervisor", wantName: "supervisorctl", wantArgs: []string{"status", "ssh"}},
		{manager: "pm2", user: "app", wantName: "su", wantArgs: []string{"-l", "app", "-c", "pm2 jlist"}},
	}
	for _, item := range cases {
		cmd, err := buildDetailCommand(item.manager, "ssh", item.user)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", item.manager, err)
		}
		if cmd.Name != item.wantName {
			t.Fatalf("unexpected name for %s: %s", item.manager, cmd.Name)
		}
		if !reflect.DeepEqual(cmd.Args, item.wantArgs) {
			t.Fatalf("unexpected args for %s: %v", item.manager, cmd.Args)
		}
	}
	if _, err := buildDetailCommand("unknown", "ssh", ""); err == nil {
		t.Fatalf("expected error for unknown manager")
	}
}

func TestBuildLogsCommand(t *testing.T) {
	cmd, err := buildLogsCommand("systemd", "ssh", LogQuery{Lines: 0})
	if err != nil || cmd.Name != "journalctl" {
		t.Fatalf("unexpected systemd logs command")
	}
	if !reflect.DeepEqual(cmd.Args, []string{"-u", "ssh", "-n", "200", "--no-pager"}) {
		t.Fatalf("unexpected systemd args: %v", cmd.Args)
	}

	cmd, err = buildLogsCommand("supervisor", "app", LogQuery{Lines: 50})
	if err != nil || cmd.Name != "supervisorctl" {
		t.Fatalf("unexpected supervisor logs command")
	}
	if !reflect.DeepEqual(cmd.Args, []string{"tail", "-n", "50", "app"}) {
		t.Fatalf("unexpected supervisor args: %v", cmd.Args)
	}

	cmd, err = buildLogsCommand("pm2", "api", LogQuery{Lines: 10, User: "node"})
	if err != nil || cmd.Name != "su" {
		t.Fatalf("unexpected pm2 logs command")
	}
	if !reflect.DeepEqual(cmd.Args, []string{"-l", "node", "-c", "pm2 logs api --lines 10 --nostream"}) {
		t.Fatalf("unexpected pm2 args: %v", cmd.Args)
	}

	if _, err := buildLogsCommand("openrc", "svc", LogQuery{Lines: 10}); err == nil {
		t.Fatalf("expected error for unsupported logs manager")
	}
}

func TestBuildActionCommands(t *testing.T) {
	cases := []struct {
		manager  string
		req      ActionRequest
		wantName string
		wantArgs []string
	}{
		{manager: "systemd", req: ActionRequest{Name: "ssh", Action: "start"}, wantName: "systemctl", wantArgs: []string{"start", "ssh"}},
		{manager: "sysv", req: ActionRequest{Name: "ssh", Action: "delete"}, wantName: "update-rc.d", wantArgs: []string{"-f", "ssh", "remove"}},
		{manager: "openrc", req: ActionRequest{Name: "ssh", Action: "disable"}, wantName: "rc-update", wantArgs: []string{"del", "ssh", "default"}},
		{manager: "supervisor", req: ActionRequest{Name: "app", Action: "delete"}, wantName: "supervisorctl", wantArgs: []string{"remove", "app"}},
		{manager: "pm2", req: ActionRequest{Name: "app", Action: "restart", User: "node"}, wantName: "su", wantArgs: []string{"-l", "node", "-c", "pm2 restart app"}},
		{manager: "pm2", req: ActionRequest{Action: "autostart", User: "node"}, wantName: "su", wantArgs: []string{"-l", "node", "-c", "pm2 save"}},
	}

	for _, item := range cases {
		cmds, err := buildActionCommands(item.manager, item.req)
		if err != nil || len(cmds) != 1 {
			t.Fatalf("unexpected action commands for %s", item.manager)
		}
		cmd := cmds[0]
		if cmd.Name != item.wantName {
			t.Fatalf("unexpected action name for %s: %s", item.manager, cmd.Name)
		}
		if !reflect.DeepEqual(cmd.Args, item.wantArgs) {
			t.Fatalf("unexpected action args for %s: %v", item.manager, cmd.Args)
		}
	}

	if _, err := buildActionCommands("systemd", ActionRequest{Name: "ssh"}); err == nil {
		t.Fatalf("expected error for missing action")
	}

	if _, err := buildActionCommands("supervisor", ActionRequest{Name: "app", Action: "autostart"}); err == nil {
		t.Fatalf("expected error for unsupported action")
	}
}

func TestBuildInstallCommands(t *testing.T) {
	cmds, err := buildInstallCommands("pm2", EnvInfo{Distro: "ubuntu"})
	if err != nil || len(cmds) < 3 {
		t.Fatalf("unexpected apt install commands")
	}
	if cmds[0].Name != "apt-get" || cmds[1].Name != "apt-get" || cmds[2].Name != "npm" {
		t.Fatalf("unexpected apt install sequence")
	}

	cmds, err = buildInstallCommands("openrc", EnvInfo{Distro: "alpine"})
	if err != nil || len(cmds) != 1 {
		t.Fatalf("unexpected apk install commands")
	}
	if cmds[0].Name != "apk" || len(cmds[0].Args) < 3 {
		t.Fatalf("unexpected apk install args")
	}

	if _, err := buildInstallCommands("systemd", EnvInfo{Distro: "alpine"}); err == nil {
		t.Fatalf("expected error for unsupported manager")
	}
	if _, err := buildInstallCommands("pm2", EnvInfo{Distro: "fedora"}); err == nil {
		t.Fatalf("expected error for unsupported distro")
	}

	if _, err := buildAptInstall("unknown"); err == nil {
		t.Fatalf("expected error for unsupported apt manager")
	}
	if _, err := buildApkInstall("unknown"); err == nil {
		t.Fatalf("expected error for unsupported apk manager")
	}
}

func TestParseListByManager(t *testing.T) {
	cases := []struct {
		manager string
		out     string
	}{
		{manager: "systemd", out: "ssh.service loaded active running OpenSSH"},
		{manager: "sysv", out: "ssh start/running"},
		{manager: "openrc", out: "ssh [ started ]"},
		{manager: "supervisor", out: "app RUNNING"},
		{manager: "pm2", out: `[{"name":"app","pid":1,"pm2_env":{"status":"online"}}]`},
	}

	for _, item := range cases {
		items, err := parseListByManager(item.manager, item.out)
		if err != nil || len(items) == 0 {
			t.Fatalf("unexpected parse result for %s", item.manager)
		}
	}

	if _, err := parseListByManager("unknown", ""); err == nil {
		t.Fatalf("expected error for unknown manager")
	}
}

func TestParsePM2DetailFromList(t *testing.T) {
	out := `[{"name":"api","pid":10,"pm2_env":{"status":"online"}}]`
	detail := parsePM2DetailFromList(out, "api")
	if detail.Name != "api" || detail.Status != ServiceStatus("running") {
		t.Fatalf("unexpected detail")
	}

	missing := parsePM2DetailFromList(out, "missing")
	if missing.Name != "missing" || missing.Status != ServiceStatus("unknown") {
		t.Fatalf("expected unknown detail")
	}
}

func TestGetServiceDetailSysV(t *testing.T) {
	src := ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		return []byte("cron start/running"), nil
	}}
	detail, err := GetServiceDetail("sysv", src, "cron", "")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if detail.Name != "cron" || detail.Status != ServiceStatus("running") {
		t.Fatalf("unexpected detail result")
	}
}

func TestRunCommandSuccess(t *testing.T) {
	output, err := RunCommand(ExecSource{}, "echo", "ok")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if string(output) == "" {
		t.Fatalf("expected output")
	}
}

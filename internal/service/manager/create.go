package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateService(manager string, src CreateSource, req CreateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("missing name")
	}
	content, err := renderCreateContent(manager, req)
	if err != nil {
		return err
	}
	dir, path, err := buildServicePath(manager, req.Name, src.BaseDir)
	if err != nil {
		return err
	}
	fs := normalizeFileSource(src.FS)
	if err := fs.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := fs.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	if manager == "sysv" || manager == "openrc" {
		if fs.Chmod != nil {
			if err := fs.Chmod(path, 0755); err != nil {
				return err
			}
		}
	}
	return runCreateCommands(manager, src.Exec, req, path)
}

func normalizeFileSource(fs FileSource) FileSource {
	if fs.WriteFile == nil {
		fs.WriteFile = func(path string, data []byte, perm os.FileMode) error {
			return os.WriteFile(path, data, perm)
		}
	}
	if fs.MkdirAll == nil {
		fs.MkdirAll = func(path string, perm os.FileMode) error {
			return os.MkdirAll(path, perm)
		}
	}
	if fs.Chmod == nil {
		fs.Chmod = func(path string, perm os.FileMode) error {
			return os.Chmod(path, perm)
		}
	}
	return fs
}

func buildServicePath(manager string, name string, baseDir string) (string, string, error) {
	dirName := managerDirName(manager)
	if dirName == "" {
		return "", "", ErrManagerUnavailable
	}
	dir := managerRoot(manager)
	if baseDir != "" {
		dir = filepath.Join(baseDir, dirName)
	}
	filename := name
	switch manager {
	case "systemd":
		filename = name + ".service"
	case "supervisor":
		filename = name + ".conf"
	case "pm2":
		filename = name + ".json"
	}
	return dir, filepath.Join(dir, filename), nil
}

func managerDirName(manager string) string {
	switch manager {
	case "systemd":
		return "systemd"
	case "sysv", "openrc":
		return "init.d"
	case "supervisor":
		return "supervisor"
	case "pm2":
		return "pm2"
	default:
		return ""
	}
}

func managerRoot(manager string) string {
	switch manager {
	case "systemd":
		return "/etc/systemd/system"
	case "sysv", "openrc":
		return "/etc/init.d"
	case "supervisor":
		return "/etc/supervisor/conf.d"
	case "pm2":
		return "/etc/pm2"
	default:
		return ""
	}
}

func renderCreateContent(manager string, req CreateRequest) (string, error) {
	raw := strings.TrimSpace(req.RawText)
	if raw != "" {
		return raw + "\n", nil
	}
	switch manager {
	case "systemd":
		return renderSystemd(req), nil
	case "sysv":
		return renderSysV(req), nil
	case "openrc":
		return renderOpenRC(req), nil
	case "supervisor":
		return renderSupervisor(req), nil
	case "pm2":
		return renderPM2(req)
	default:
		return "", ErrManagerUnavailable
	}
}

func renderSystemd(req CreateRequest) string {
	desc := req.Description
	if desc == "" {
		desc = req.Name
	}
	builder := strings.Builder{}
	builder.WriteString("[Unit]\n")
	builder.WriteString("Description=" + desc + "\n\n")
	builder.WriteString("[Service]\n")
	builder.WriteString("Type=simple\n")
	if strings.TrimSpace(req.Workdir) != "" {
		builder.WriteString("WorkingDirectory=" + strings.TrimSpace(req.Workdir) + "\n")
	}
	if strings.TrimSpace(req.Exec) != "" {
		builder.WriteString("ExecStart=" + strings.TrimSpace(req.Exec) + "\n")
	}
	if strings.TrimSpace(req.User) != "" {
		builder.WriteString("User=" + strings.TrimSpace(req.User) + "\n")
	}
	for _, env := range splitEnvLines(req.Env) {
		builder.WriteString("Environment=" + env + "\n")
	}
	builder.WriteString("\n[Install]\nWantedBy=multi-user.target\n")
	return builder.String()
}

func renderSysV(req CreateRequest) string {
	desc := req.Description
	if desc == "" {
		desc = req.Name
	}
	workdir := strings.TrimSpace(req.Workdir)
	if workdir == "" {
		workdir = "/"
	}
	cmd := strings.TrimSpace(req.Exec)
	if cmd == "" {
		cmd = req.Name
	}
	return fmt.Sprintf(`#!/bin/sh
### BEGIN INIT INFO
# Provides:          %s
# Required-Start:    $remote_fs $syslog
# Required-Stop:     $remote_fs $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Description:       %s
### END INIT INFO

PIDFILE=/var/run/%s.pid
DAEMON_CMD=%q
WORKDIR=%q

start() {
  start-stop-daemon --start --background --make-pidfile --pidfile "$PIDFILE" --chdir "$WORKDIR" --startas /bin/sh -- -c "$DAEMON_CMD"
}

stop() {
  start-stop-daemon --stop --pidfile "$PIDFILE" --retry 5
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart)
    stop
    start
    ;;
  status)
    if [ -f "$PIDFILE" ]; then
      echo "running"
    else
      echo "stopped"
    fi
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|status}"
    exit 1
    ;;
esac

exit 0
`, req.Name, desc, req.Name, cmd, workdir)
}

func renderOpenRC(req CreateRequest) string {
	desc := req.Description
	if desc == "" {
		desc = req.Name
	}
	workdir := strings.TrimSpace(req.Workdir)
	if workdir == "" {
		workdir = "/"
	}
	cmd := strings.TrimSpace(req.Exec)
	if cmd == "" {
		cmd = req.Name
	}
	escaped := strings.ReplaceAll(cmd, "\"", "\\\"")
	return fmt.Sprintf(`#!/sbin/openrc-run

name=%q
description=%q
command="/bin/sh"
command_args="-c \"%s\""
command_background="yes"
pidfile="/run/%s.pid"
directory=%q
`, req.Name, desc, escaped, req.Name, workdir)
}

func renderSupervisor(req CreateRequest) string {
	desc := req.Description
	if desc == "" {
		desc = req.Name
	}
	user := strings.TrimSpace(req.User)
	if user == "" {
		user = "root"
	}
	workdir := strings.TrimSpace(req.Workdir)
	if workdir == "" {
		workdir = "/"
	}
	cmd := strings.TrimSpace(req.Exec)
	if cmd == "" {
		cmd = req.Name
	}
	autostart := "false"
	if req.Options.Autostart {
		autostart = "true"
	}
	return fmt.Sprintf(`[program:%s]
command=%s
directory=%s
user=%s
autostart=%s
autorestart=true
stdout_logfile=/var/log/%s.log
stderr_logfile=/var/log/%s.err
description=%s
`, req.Name, cmd, workdir, user, autostart, req.Name, req.Name, desc)
}

func renderPM2(req CreateRequest) (string, error) {
	cmd := strings.TrimSpace(req.Exec)
	if cmd == "" {
		cmd = req.Name
	}
	workdir := strings.TrimSpace(req.Workdir)
	if workdir == "" {
		workdir = "/"
	}
	payload := map[string]any{
		"apps": []map[string]any{
			{
				"name":   req.Name,
				"script": cmd,
				"cwd":    workdir,
				"env":    parseEnvMap(req.Env),
			},
		},
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func splitEnvLines(raw string) []string {
	lines := strings.Split(raw, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		value := strings.TrimSpace(line)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}

func parseEnvMap(raw string) map[string]string {
	envs := map[string]string{}
	for _, line := range splitEnvLines(raw) {
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			envs[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return envs
}

func runCreateCommands(manager string, src ExecSource, req CreateRequest, path string) error {
	switch manager {
	case "systemd":
		if _, err := RunCommand(src, "systemctl", "daemon-reload"); err != nil {
			return err
		}
		if req.Options.Autostart {
			if _, err := RunCommand(src, "systemctl", "enable", req.Name); err != nil {
				return err
			}
		}
		if req.Options.Start {
			if _, err := RunCommand(src, "systemctl", "start", req.Name); err != nil {
				return err
			}
		}
	case "sysv":
		if req.Options.Autostart {
			if _, err := RunCommand(src, "update-rc.d", req.Name, "defaults"); err != nil {
				return err
			}
		}
		if req.Options.Start {
			if _, err := RunCommand(src, "service", req.Name, "start"); err != nil {
				return err
			}
		}
	case "openrc":
		if req.Options.Autostart {
			if _, err := RunCommand(src, "rc-update", "add", req.Name, "default"); err != nil {
				return err
			}
		}
		if req.Options.Start {
			if _, err := RunCommand(src, "rc-service", req.Name, "start"); err != nil {
				return err
			}
		}
	case "supervisor":
		if _, err := RunCommand(src, "supervisorctl", "reread"); err != nil {
			return err
		}
		if _, err := RunCommand(src, "supervisorctl", "update"); err != nil {
			return err
		}
		if req.Options.Start {
			if _, err := RunCommand(src, "supervisorctl", "start", req.Name); err != nil {
				return err
			}
		}
	case "pm2":
		if req.Options.Start {
			cmd, err := buildPM2Command(req.User, "start", path, "--name", req.Name)
			if err != nil {
				return err
			}
			if _, err := RunCommand(src, cmd.Name, cmd.Args...); err != nil {
				return err
			}
		}
		if req.Options.Autostart {
			cmd, err := buildPM2Command(req.User, "save")
			if err != nil {
				return err
			}
			if _, err := RunCommand(src, cmd.Name, cmd.Args...); err != nil {
				return err
			}
		}
	default:
		return ErrManagerUnavailable
	}
	return nil
}

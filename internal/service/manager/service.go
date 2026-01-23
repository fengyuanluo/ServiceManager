package manager

import (
	"errors"
	"fmt"
	"strings"
)

type commandSpec struct {
	Name string
	Args []string
}

func ListServices(manager string, src ExecSource, user string) ([]ServiceSummary, error) {
	cmd, err := buildListCommand(manager, user)
	if err != nil {
		return nil, err
	}
	out, err := RunCommand(src, cmd.Name, cmd.Args...)
	if err != nil {
		return nil, ErrManagerUnavailable
	}
	return parseListByManager(manager, string(out))
}

func GetServiceDetail(manager string, src ExecSource, name string, user string) (ServiceDetail, error) {
	cmd, err := buildDetailCommand(manager, name, user)
	if err != nil {
		return ServiceDetail{}, err
	}
	out, err := RunCommand(src, cmd.Name, cmd.Args...)
	if err != nil {
		return ServiceDetail{}, ErrManagerUnavailable
	}
	if manager == "systemd" {
		return parseSystemdDetail(string(out), name), nil
	}
	if manager == "pm2" {
		return parsePM2DetailFromList(string(out), name), nil
	}
	items, err := parseListByManager(manager, string(out))
	if err != nil || len(items) == 0 {
		return ServiceDetail{Name: name, Status: ServiceStatus("unknown")}, nil
	}
	return ServiceDetail{Name: items[0].Name, Status: items[0].Status, Description: items[0].Detail}, nil
}

func GetServiceLogs(manager string, src ExecSource, name string, query LogQuery) ([]string, error) {
	cmd, err := buildLogsCommand(manager, name, query)
	if err != nil {
		return nil, err
	}
	out, err := RunCommand(src, cmd.Name, cmd.Args...)
	if err != nil {
		return nil, ErrManagerUnavailable
	}
	lines := splitLines(string(out))
	if query.Keyword == "" {
		return lines, nil
	}
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(line, query.Keyword) {
			filtered = append(filtered, line)
		}
	}
	return filtered, nil
}

func RunServiceAction(manager string, src ExecSource, req ActionRequest) error {
	commands, err := buildActionCommands(manager, req)
	if err != nil {
		return err
	}
	for _, cmd := range commands {
		if _, err := RunCommand(src, cmd.Name, cmd.Args...); err != nil {
			return err
		}
	}
	return nil
}

func RunServiceBatch(manager string, src ExecSource, user string, items []BatchItem) []BatchResult {
	results := make([]BatchResult, 0, len(items))
	for _, item := range items {
		err := RunServiceAction(manager, src, ActionRequest{Name: item.Name, Action: string(item.Action), User: user})
		result := BatchResult{
			Name:    item.Name,
			Action:  item.Action,
			Success: err == nil,
		}
		if err != nil {
			result.Message = err.Error()
		}
		results = append(results, result)
	}
	return results
}

func InstallManager(manager string, env EnvInfo, src ExecSource) InstallResult {
	commands, err := buildInstallCommands(manager, env)
	if err != nil {
		return InstallResult{Manager: manager, Success: false, Message: err.Error()}
	}
	for _, cmd := range commands {
		if _, err := RunCommand(src, cmd.Name, cmd.Args...); err != nil {
			return InstallResult{Manager: manager, Success: false, Message: err.Error()}
		}
	}
	return InstallResult{Manager: manager, Success: true}
}

func buildListCommand(manager string, user string) (commandSpec, error) {
	switch manager {
	case "systemd":
		return commandSpec{Name: "systemctl", Args: []string{"list-units", "--type=service", "--no-legend", "--no-pager"}}, nil
	case "sysv":
		return commandSpec{Name: "service", Args: []string{"--status-all"}}, nil
	case "openrc":
		return commandSpec{Name: "rc-service", Args: []string{"--status-all"}}, nil
	case "supervisor":
		return commandSpec{Name: "supervisorctl", Args: []string{"status"}}, nil
	case "pm2":
		return buildPM2Command(user, "jlist")
	default:
		return commandSpec{}, ErrManagerUnavailable
	}
}

func buildDetailCommand(manager string, name string, user string) (commandSpec, error) {
	switch manager {
	case "systemd":
		return commandSpec{
			Name: "systemctl",
			Args: []string{"show", name, "--no-pager", "--property=Id,Description,MainPID,ActiveState,ActiveEnterTimestamp,FragmentPath,User"},
		}, nil
	case "sysv":
		return commandSpec{Name: "service", Args: []string{name, "status"}}, nil
	case "openrc":
		return commandSpec{Name: "rc-service", Args: []string{name, "status"}}, nil
	case "supervisor":
		return commandSpec{Name: "supervisorctl", Args: []string{"status", name}}, nil
	case "pm2":
		return buildPM2Command(user, "jlist")
	default:
		return commandSpec{}, ErrManagerUnavailable
	}
}

func buildLogsCommand(manager string, name string, query LogQuery) (commandSpec, error) {
	lines := query.Lines
	if lines <= 0 {
		lines = 200
	}
	switch manager {
	case "systemd":
		return commandSpec{Name: "journalctl", Args: []string{"-u", name, "-n", fmt.Sprintf("%d", lines), "--no-pager"}}, nil
	case "supervisor":
		return commandSpec{Name: "supervisorctl", Args: []string{"tail", "-n", fmt.Sprintf("%d", lines), name}}, nil
	case "pm2":
		return buildPM2Command(query.User, "logs", name, "--lines", fmt.Sprintf("%d", lines), "--nostream")
	default:
		return commandSpec{}, ErrManagerUnavailable
	}
}

func buildActionCommands(manager string, req ActionRequest) ([]commandSpec, error) {
	action := strings.ToLower(strings.TrimSpace(req.Action))
	if action == "" {
		return nil, errors.New("missing action")
	}
	switch manager {
	case "systemd":
		switch action {
		case "start", "stop", "restart":
			return []commandSpec{{Name: "systemctl", Args: []string{action, req.Name}}}, nil
		case "autostart", "enable":
			return []commandSpec{{Name: "systemctl", Args: []string{"enable", req.Name}}}, nil
		case "disable", "noautostart":
			return []commandSpec{{Name: "systemctl", Args: []string{"disable", req.Name}}}, nil
		case "delete":
			return []commandSpec{{Name: "systemctl", Args: []string{"disable", "--now", req.Name}}}, nil
		}
	case "sysv":
		switch action {
		case "start", "stop", "restart":
			return []commandSpec{{Name: "service", Args: []string{req.Name, action}}}, nil
		case "autostart", "enable":
			return []commandSpec{{Name: "update-rc.d", Args: []string{req.Name, "defaults"}}}, nil
		case "disable", "noautostart":
			return []commandSpec{{Name: "update-rc.d", Args: []string{req.Name, "disable"}}}, nil
		case "delete":
			return []commandSpec{{Name: "update-rc.d", Args: []string{"-f", req.Name, "remove"}}}, nil
		}
	case "openrc":
		switch action {
		case "start", "stop", "restart":
			return []commandSpec{{Name: "rc-service", Args: []string{req.Name, action}}}, nil
		case "autostart", "enable":
			return []commandSpec{{Name: "rc-update", Args: []string{"add", req.Name, "default"}}}, nil
		case "disable", "noautostart":
			return []commandSpec{{Name: "rc-update", Args: []string{"del", req.Name, "default"}}}, nil
		case "delete":
			return []commandSpec{{Name: "rc-update", Args: []string{"del", req.Name, "default"}}}, nil
		}
	case "supervisor":
		switch action {
		case "start", "stop", "restart":
			return []commandSpec{{Name: "supervisorctl", Args: []string{action, req.Name}}}, nil
		case "delete":
			return []commandSpec{{Name: "supervisorctl", Args: []string{"remove", req.Name}}}, nil
		}
	case "pm2":
		switch action {
		case "start", "stop", "restart", "delete":
			cmd, err := buildPM2Command(req.User, action, req.Name)
			if err != nil {
				return nil, err
			}
			return []commandSpec{cmd}, nil
		case "autostart", "enable":
			cmd, err := buildPM2Command(req.User, "save")
			if err != nil {
				return nil, err
			}
			return []commandSpec{cmd}, nil
		}
	}
	return nil, errors.New("unsupported action")
}

func buildInstallCommands(manager string, env EnvInfo) ([]commandSpec, error) {
	distro := strings.ToLower(env.Distro)
	switch distro {
	case "ubuntu", "debian":
		return buildAptInstall(manager)
	case "alpine":
		return buildApkInstall(manager)
	default:
		return nil, errors.New("unsupported distro")
	}
}

func buildAptInstall(manager string) ([]commandSpec, error) {
	switch manager {
	case "sysv":
		return []commandSpec{
			{Name: "apt-get", Args: []string{"update"}},
			{Name: "apt-get", Args: []string{"install", "-y", "sysvinit-utils"}},
		}, nil
	case "openrc":
		return []commandSpec{
			{Name: "apt-get", Args: []string{"update"}},
			{Name: "apt-get", Args: []string{"install", "-y", "openrc"}},
		}, nil
	case "supervisor":
		return []commandSpec{
			{Name: "apt-get", Args: []string{"update"}},
			{Name: "apt-get", Args: []string{"install", "-y", "supervisor"}},
		}, nil
	case "pm2":
		return []commandSpec{
			{Name: "apt-get", Args: []string{"update"}},
			{Name: "apt-get", Args: []string{"install", "-y", "npm"}},
			{Name: "npm", Args: []string{"install", "-g", "pm2"}},
		}, nil
	default:
		return nil, errors.New("unsupported manager")
	}
}

func buildApkInstall(manager string) ([]commandSpec, error) {
	switch manager {
	case "openrc":
		return []commandSpec{{Name: "apk", Args: []string{"add", "--no-cache", "openrc"}}}, nil
	case "supervisor":
		return []commandSpec{{Name: "apk", Args: []string{"add", "--no-cache", "supervisor"}}}, nil
	case "pm2":
		return []commandSpec{
			{Name: "apk", Args: []string{"add", "--no-cache", "nodejs", "npm"}},
			{Name: "npm", Args: []string{"install", "-g", "pm2"}},
		}, nil
	default:
		return nil, errors.New("unsupported manager")
	}
}

func buildPM2Command(user string, args ...string) (commandSpec, error) {
	if user == "" || user == "root" {
		return commandSpec{Name: "pm2", Args: args}, nil
	}
	cmd := "pm2 " + strings.Join(args, " ")
	return commandSpec{Name: "su", Args: []string{"-l", user, "-c", cmd}}, nil
}

func parseListByManager(manager string, out string) ([]ServiceSummary, error) {
	switch manager {
	case "systemd":
		return ParseSystemdList(out)
	case "sysv":
		return ParseSysVList(out)
	case "openrc":
		return ParseOpenRCList(out)
	case "supervisor":
		return ParseSupervisorList(out)
	case "pm2":
		return ParsePM2List(out)
	default:
		return nil, ErrManagerUnavailable
	}
}

func parseSystemdDetail(out string, fallbackName string) ServiceDetail {
	info := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		info[parts[0]] = parts[1]
	}
	name := strings.TrimSuffix(info["Id"], ".service")
	if name == "" {
		name = fallbackName
	}
	return ServiceDetail{
		Name:        name,
		Status:      NormalizeStatus(info["ActiveState"]),
		PID:         info["MainPID"],
		Since:       info["ActiveEnterTimestamp"],
		UnitPath:    info["FragmentPath"],
		Description: info["Description"],
		User:        info["User"],
	}
}

func parsePM2DetailFromList(out string, name string) ServiceDetail {
	items, err := ParsePM2List(out)
	if err != nil {
		return ServiceDetail{Name: name, Status: ServiceStatus("unknown")}
	}
	for _, item := range items {
		if item.Name == name {
			return ServiceDetail{Name: item.Name, Status: item.Status, Description: item.Detail}
		}
	}
	return ServiceDetail{Name: name, Status: ServiceStatus("unknown")}
}

func splitLines(out string) []string {
	lines := strings.Split(out, "\n")
	results := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		results = append(results, line)
	}
	return results
}

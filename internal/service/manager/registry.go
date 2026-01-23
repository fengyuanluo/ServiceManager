package manager

func DetectManagers(env EnvInfo) []string {
	out := []string{}
	if env.Init == "systemd" {
		out = append(out, "systemd")
	}
	if env.Init == "openrc" {
		out = append(out, "openrc")
	}
	if env.Distro != "unknown" {
		out = append(out, "sysv")
	}
	return out
}

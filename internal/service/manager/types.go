package manager

type ServiceStatus string

type ServiceSummary struct {
	Name   string
	Status ServiceStatus
	Detail string
}

type EnvInfo struct {
	Distro string
	Init   string
}

type EnvSource struct {
	ReadFile func(path string) ([]byte, error)
	Exec     func(name string, args ...string) ([]byte, error)
}

type ExecSource struct {
	Run func(name string, args ...string) ([]byte, error)
}

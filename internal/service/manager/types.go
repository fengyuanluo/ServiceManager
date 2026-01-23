package manager

type ServiceStatus string

type ServiceSummary struct {
	Name   string        `json:"name"`
	Status ServiceStatus `json:"status"`
	Detail string        `json:"detail,omitempty"`
}

type ServiceDetail struct {
	Name        string            `json:"name"`
	Status      ServiceStatus     `json:"status"`
	PID         string            `json:"pid,omitempty"`
	Since       string            `json:"since,omitempty"`
	UnitPath    string            `json:"unit_path,omitempty"`
	Description string            `json:"description,omitempty"`
	User        string            `json:"user,omitempty"`
	Extra       map[string]string `json:"extra,omitempty"`
}

type LogQuery struct {
	Lines   int    `json:"lines"`
	Keyword string `json:"keyword,omitempty"`
	User    string `json:"user,omitempty"`
}

type ActionRequest struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	User   string `json:"user,omitempty"`
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

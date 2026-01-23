package manager

type InstallRequest struct {
	Manager string `json:"manager"`
}

type InstallResult struct {
	Manager  string `json:"manager"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	ExitCode int    `json:"exit_code,omitempty"`
}

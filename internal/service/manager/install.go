package manager

type InstallRequest struct {
	Manager string
}

type InstallResult struct {
	Manager  string
	Success  bool
	Message  string
	ExitCode int
}

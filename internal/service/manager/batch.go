package manager

type BatchAction string

type BatchItem struct {
	Name   string
	Action BatchAction
}

type BatchResult struct {
	Name    string
	Action  BatchAction
	Success bool
	Message string
}

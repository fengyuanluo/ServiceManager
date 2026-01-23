package manager

type BatchAction string

type BatchItem struct {
	Name   string      `json:"name"`
	Action BatchAction `json:"action"`
}

type BatchResult struct {
	Name    string      `json:"name"`
	Action  BatchAction `json:"action"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
}

package manager

import "errors"

var (
	ErrManagerUnavailable = errors.New("manager unavailable")
	errCommandNotFound    = errors.New("command not found")
)

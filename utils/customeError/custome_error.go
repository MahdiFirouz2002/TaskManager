package customeError

import "errors"

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidID         = errors.New("invalid ID format")
	ErrInvalidTaskFormat = errors.New("invalid Task format")
)

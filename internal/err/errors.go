package err

import "errors"

var (
	ErrNoServerRegistered = errors.New("no server registered")
	ErrNoHealthyServer    = errors.New("no healthy server available")
	ErrEmptyServerName    = errors.New("service name must not be empty")
)

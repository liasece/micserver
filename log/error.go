package log

import (
	"errors"
)

var (
	ErrNilLogger       = errors.New("logger is nil")
	ErrUnknownLogLevel = errors.New("unknown log level")
)

package log

import (
	"errors"
)

// error
var (
	ErrNilLogger       = errors.New("logger is nil")
	ErrUnknownLogLevel = errors.New("unknown log level")
)

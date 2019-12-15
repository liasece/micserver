package log

import (
	"sync"
)

var (
	LEVEL_FLAGS = []string{
		"[S]",
		"[D]",
		"[I]",
		"[WARNING]",
		"[ERROR]",
		"[FATALERROR]",
	}
)

const (
	SYS = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

// default
var (
	default_logger *Logger
)

func Syslog(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(SYS, fmt, args...)
}

func Debug(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(FATAL, fmt, args...)
}

func SetDefaultLogger(l *Logger) {
	if l != default_logger {
		// default_logger.GetLogWriter().Close()
		default_logger = l
	}
}

func GetDefaultLogger() *Logger {
	return default_logger
}

func init() {
	default_logger = NewLogger(false, "")
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
	default_logger.SetLogName("log")
}

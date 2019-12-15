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

func IsSyslogEnable() bool {
	return default_logger.IsSyslogEnable()
}

func IsDebugEnable() bool {
	return default_logger.IsDebugEnable()
}

func IsInfoEnable() bool {
	return default_logger.IsInfoEnable()
}

func IsWarnEnable() bool {
	return default_logger.IsWarnEnable()
}

func IsErrorEnable() bool {
	return default_logger.IsErrorEnable()
}

func IsFatalEnable() bool {
	return default_logger.IsFatalEnable()
}

func SetLogLevel(lvl int32) {
	default_logger.SetLogLevel(lvl)
}

func SetLogLevelByStr(lvl string) error {
	return default_logger.SetLogLevelByStr(lvl)
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

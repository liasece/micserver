package log

import (
	"time"
)

type ILogger interface {
	Syslog(fmt string, args ...interface{})
	Debug(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})
	Fatal(fmt string, args ...interface{})
	IsSyslogEnable() bool
	IsDebugEnable() bool
	IsWarnEnable() bool
	IsInfoEnable() bool
	IsErrorEnable() bool
	IsFatalEnable() bool
	Clone() *Logger
	SetTopic(topic string)
	GetLogLevel() int32
	SetLogName(logname string)
	GetLogger() *Logger
}

type Writer interface {
	Init() error
	Write(*Record) error
	GetType() writerType
}

type Rotater interface {
	Rotate() error
	RotateByTime(*time.Time) error
	SetPathPattern(string, string) error
}

type Flusher interface {
	Flush() error
}

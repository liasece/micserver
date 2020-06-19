package log

import (
	"time"
)

// ILogger 日志系统实现的接口
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

// Writer 输出器实现的接口
type Writer interface {
	Write(*Record) error
}

// Rotater 转储器实现的接口
type Rotater interface {
	Rotate() error
	RotateByTime(*time.Time) error
}

// Flusher 刷新输出
type Flusher interface {
	Flush() error
}

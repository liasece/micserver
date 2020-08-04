// Package log micserver 中使用的 log 系统，支持 Syslog(系统信息)/Debug(调试)/Info(关心的)/
// Warning(警告)/Error(错误)/Fatal(致命) 日志级别，日志可按大小/小时转储。
// 建议：
// 系统运行的各个细节/详细步骤使用 Syslog 级别记录，为了减小 log 系统造成的运算负担，
// 应该先使用 IsSyslogEnable() 判断 Syslog 级别的日志是否开启，再调用 Syslog() 记录；
// 业务的主要环节/关键节点的调试使用 Debug 级别记录；
// 业务造成的结果，或需要在后续运营维护中查看用户信息变更，使用 Info 级别记录；
// 客户端/用户的输入错误，或者系统设计意外的条件不满足，使用 Warning 级别记录；
// 该分布式系统内部的错误值（与客户端等无关），但是可以恢复或者对业务逻辑没有影响时，
// 使用 Error 级别记录；
// 在 Error 级别的基础上，如果错误无法代码恢复或者对业务逻辑产生必要影响时，使用 Fatal
// 级别记录。
// 在生产环境/正式环境中，应该将日志等级至高调整至 Info ，在需要必要的调试信息时，
// 可调整至 Debug 。不要在正式环境中启用 Syslog 日志等级，你只应该在开发环境中使用它。
// Warning / Error / Fatal 日志级别，无论何时你都要谨慎关闭他们，如果你确定不关心
// 你业务的信息，可以调整至 Warning 级别，但是无论如何，关闭警告或者错误都是一个极具风险
// 的决定。
package log

import (
	"sync"
)

// 默认的日志记录器
var (
	defaultLogger *Logger
)

// Syslog 默认 Logger 异步输出一条 Syslog 级别的日志
func Syslog(fmt string, args ...interface{}) {
	defaultLogger.deliverRecordToWriter(SysLevel, fmt, args...)
}

// Debug 默认 Logger 异步输出一条 Debug 级别的日志
func Debug(fmt string, args ...interface{}) {
	defaultLogger.deliverRecordToWriter(DebugLevel, fmt, args...)
}

// Warn 默认 Logger 异步输出一条 Warn 级别的日志
func Warn(fmt string, args ...interface{}) {
	defaultLogger.deliverRecordToWriter(WarnLevel, fmt, args...)
}

// Info 默认 Logger 异步输出一条 Info 级别的日志
func Info(fmt string, args ...interface{}) {
	defaultLogger.deliverRecordToWriter(InfoLevel, fmt, args...)
}

// Error 默认 Logger 异步输出一条 Error 级别的日志
func Error(fmt string, args ...interface{}) {
	defaultLogger.deliverRecordToWriter(ErrorLevel, fmt, args...)
}

// DPanic 默认 Logger 异步输出一条 DPanic 级别的日志
func DPanic(fmt string, args ...interface{}) {
	defaultLogger.deliverRecordToWriter(DPanicLevel, fmt, args...)
}

// Panic 默认 Logger 异步输出一条 Panic 级别的日志
func Panic(fmt string, args ...interface{}) {
	defaultLogger.deliverRecordToWriter(PanicLevel, fmt, args...)
}

// Fatal 默认 Logger 异步输出一条 Fatal 级别的日志
func Fatal(fmt string, args ...interface{}) {
	defaultLogger.deliverRecordToWriter(FatalLevel, fmt, args...)
}

// IsSyslogEnable 默认 Logger 判断 Syslog 日志级别是否开启
func IsSyslogEnable() bool {
	return defaultLogger.IsSyslogEnable()
}

// IsDebugEnable 默认 Logger 判断 Debug 日志级别是否开启
func IsDebugEnable() bool {
	return defaultLogger.IsDebugEnable()
}

// IsInfoEnable 默认 Logger 判断 Warn 日志级别是否开启
func IsInfoEnable() bool {
	return defaultLogger.IsInfoEnable()
}

// IsWarnEnable 默认 Logger 判断 Info 日志级别是否开启
func IsWarnEnable() bool {
	return defaultLogger.IsWarnEnable()
}

// IsErrorEnable 默认 Logger 判断 Error 日志级别是否开启
func IsErrorEnable() bool {
	return defaultLogger.IsErrorEnable()
}

// IsPanicEnable 默认 Logger 判断 Panic 日志级别是否开启
func IsPanicEnable() bool {
	return defaultLogger.IsPanicEnable()
}

// IsDPanicEnable 默认 Logger 判断 DPanic 日志级别是否开启
func IsDPanicEnable() bool {
	return defaultLogger.IsDPanicEnable()
}

// IsFatalEnable 默认 Logger 判断 Fatal 日志级别是否开启
func IsFatalEnable() bool {
	return defaultLogger.IsFatalEnable()
}

// SetLogLevel 默认 Logger 设置日志等级
func SetLogLevel(lvl Level) {
	defaultLogger.SetLogLevel(lvl)
}

// SetLogLevelByStr 默认 Logger 使用等级名设置日志等级
func SetLogLevelByStr(lvl string) error {
	return defaultLogger.SetLogLevelByStr(lvl)
}

// SetDefaultLogger 设置默认 Logger
func SetDefaultLogger(l *Logger) {
	if l != defaultLogger {
		// defaultLogger.GetLogWriter().Close()
		defaultLogger = l
	}
}

// GetDefaultLogger 获取默认 Logger
func GetDefaultLogger() *Logger {
	return defaultLogger
}

// init 在程序启动时，初始化一个默认 Logger
func init() {
	defaultLogger = NewLogger(nil)
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
	defaultLogger.SetLogName("log")
}

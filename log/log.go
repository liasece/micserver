/*
micserver 中使用的 log 系统，支持 Syslog(系统信息)/Debug(调试)/Info(关心的)/
Warning(警告)/Error(错误)/Fatal(致命) 日志级别，日志可按大小/小时转储。
建议：
系统运行的各个细节/详细步骤使用 Syslog 级别记录，为了减小 log 系统造成的运算负担，
应该先使用 IsSyslogEnable() 判断 Syslog 级别的日志是否开启，再调用 Syslog() 记录；
业务的主要环节/关键节点的调试使用 Debug 级别记录；
业务造成的结果，或需要在后续运营维护中查看用户信息变更，使用 Info 级别记录；
客户端/用户的输入错误，或者系统设计意外的条件不满足，使用 Warning 级别记录；
该分布式系统内部的错误值（与客户端等无关），但是可以恢复或者对业务逻辑没有影响时，
使用 Error 级别记录；
在 Error 级别的基础上，如果错误无法代码恢复或者对业务逻辑产生必要影响时，使用 Fatal
级别记录。
在生产环境/正式环境中，应该将日志等级至高调整至 Info ，在需要必要的调试信息时，
可调整至 Debug 。不要在正式环境中启用 Syslog 日志等级，你只应该在开发环境中使用它。
Warning / Error / Fatal 日志级别，无论何时你都要谨慎关闭他们，如果你确定不关心
你业务的信息，可以调整至 Warning 级别，但是无论如何，关闭警告或者错误都是一个极具风险
的决定。
*/
package log

import (
	"sync"
)

// 各个日志等级在一条Log中的日志等级标题
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

// 日志等级
const (
	SYS = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

// 默认的日志记录器
var (
	default_logger *Logger
)

// 默认 Logger 异步输出一条 Syslog 级别的日志
func Syslog(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(SYS, fmt, args...)
}

// 默认 Logger 异步输出一条 Debug 级别的日志
func Debug(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(DEBUG, fmt, args...)
}

// 默认 Logger 异步输出一条 Warn 级别的日志
func Warn(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(WARNING, fmt, args...)
}

// 默认 Logger 异步输出一条 Info 级别的日志
func Info(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(INFO, fmt, args...)
}

// 默认 Logger 异步输出一条 Error 级别的日志
func Error(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(ERROR, fmt, args...)
}

// 默认 Logger 异步输出一条 Fatal 级别的日志
func Fatal(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(FATAL, fmt, args...)
}

// 默认 Logger 判断 Syslog 日志级别是否开启
func IsSyslogEnable() bool {
	return default_logger.IsSyslogEnable()
}

// 默认 Logger 判断 Debug 日志级别是否开启
func IsDebugEnable() bool {
	return default_logger.IsDebugEnable()
}

// 默认 Logger 判断 Warn 日志级别是否开启
func IsInfoEnable() bool {
	return default_logger.IsInfoEnable()
}

// 默认 Logger 判断 Info 日志级别是否开启
func IsWarnEnable() bool {
	return default_logger.IsWarnEnable()
}

// 默认 Logger 判断 Error 日志级别是否开启
func IsErrorEnable() bool {
	return default_logger.IsErrorEnable()
}

// 默认 Logger 判断 Fatal 日志级别是否开启
func IsFatalEnable() bool {
	return default_logger.IsFatalEnable()
}

// 默认 Logger 设置日志等级
func SetLogLevel(lvl int32) {
	default_logger.SetLogLevel(lvl)
}

// 默认 Logger 使用等级名设置日志等级
func SetLogLevelByStr(lvl string) error {
	return default_logger.SetLogLevelByStr(lvl)
}

// 设置默认 Logger
func SetDefaultLogger(l *Logger) {
	if l != default_logger {
		// default_logger.GetLogWriter().Close()
		default_logger = l
	}
}

// 获取默认 Logger
func GetDefaultLogger() *Logger {
	return default_logger
}

// 在程序启动时，初始化一个默认 Logger
func init() {
	default_logger = NewLogger(false, "")
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
	default_logger.SetLogName("log")
}

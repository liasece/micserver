package log

import (
	"strings"

	"github.com/liasece/micserver/log/core"
)

// Logger 日志实例
type Logger struct {
	options
	logWriter *recordWriter
}

// NewLogger 构造一个日志
func NewLogger(core core.Core, optps ...Option) *Logger {
	l := new(Logger)
	l.options = _defaultoptions
	for _, opt := range optps {
		opt.apply(l)
	}
	l.logWriter = &recordWriter{}
	l.logWriter.Init(&l.options)

	for _, path := range l.options.FilePaths {
		l.logWriter.AddLogFile(path, &l.options)
	}

	if !l.options.NoConsole {
		w := newConsoleWriter()
		w.SetColor(!l.options.NoConsoleColor)
		l.logWriter.registerLogWriter(w)
	}

	return l
}

// Clone 浅拷贝出一个 Logger ，他们具备相同的底层写入接口，
// 需要注意的是，克隆出来的logger对象在写入器中会依旧受到源拷贝对象的影响
func (l *Logger) Clone() *Logger {
	if l == nil && l != defaultLogger {
		return defaultLogger.Clone()
	}
	res := *l
	return &res
}

// SetLogName 设置日志名称，一般使用进程或者模块名字
func (l *Logger) SetLogName(logname string) {
	if l == nil && l != defaultLogger {
		defaultLogger.SetLogName(logname)
		return
	}
	l.options.Name = logname
}

// getLogLevel 获取当前日志等级
func (l *Logger) getLogLevel() Level {
	if l == nil && l != defaultLogger {
		return defaultLogger.getLogLevel()
	}
	return l.options.Level
}

// IsSyslogEnable 判断 Syslog 日志级别是否开启
func (l *Logger) IsSyslogEnable() bool {
	return l.getLogLevel() >= SYS
}

// IsDebugEnable 判断 Debug 日志级别是否开启
func (l *Logger) IsDebugEnable() bool {
	return l.getLogLevel() >= DEBUG
}

// IsInfoEnable 判断 Warn 日志级别是否开启
func (l *Logger) IsInfoEnable() bool {
	return l.getLogLevel() >= DEBUG
}

// IsWarnEnable 判断 Info 日志级别是否开启
func (l *Logger) IsWarnEnable() bool {
	return l.getLogLevel() >= WARNING
}

// IsErrorEnable 判断 Error 日志级别是否开启
func (l *Logger) IsErrorEnable() bool {
	return l.getLogLevel() >= ERROR
}

// IsFatalEnable 判断 Fatal 日志级别是否开启
func (l *Logger) IsFatalEnable() bool {
	return l.getLogLevel() >= FATAL
}

// SetLogLevel 设置日志等级
func (l *Logger) SetLogLevel(lvl Level) {
	if l == nil && l != defaultLogger {
		defaultLogger.SetLogLevel(lvl)
		return
	}
	l.options.Level = lvl
}

// SetLogLevelByStr 使用等级名设置日志等级
func (l *Logger) SetLogLevelByStr(lvl string) error {
	if l == nil && l != defaultLogger {
		return defaultLogger.SetLogLevelByStr(lvl)
	}
	lvlUpper := strings.ToUpper(lvl)
	switch lvlUpper {
	case "SYS":
		l.options.Level = SYS
	case "DEBUG":
		l.options.Level = DEBUG
	case "INFO":
		l.options.Level = INFO
	case "WARNING":
		l.options.Level = WARNING
	case "ERROR":
		l.options.Level = ERROR
	case "FATAL":
		l.options.Level = FATAL
	default:
		return ErrUnknownLogLevel
	}
	return nil
}

// SetTopic 设置该日志主题，一般设置为 SetLogName() 的下一级系统名称
func (l *Logger) SetTopic(topic string) error {
	if l == nil || l == defaultLogger {
		// nil or default logger can't set topic
		return ErrNilLogger
	}
	l.options.Topic = topic
	return nil
}

// Syslog 异步输出一条 Syslog 级别的日志
func (l *Logger) Syslog(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(SYS, fmt, args...)
}

// Debug 异步输出一条 Debug 级别的日志
func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

// Warn 异步输出一条 Warn 级别的日志
func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

// Info 异步输出一条 Info 级别的日志
func (l *Logger) Info(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

// Error 异步输出一条 Error 级别的日志
func (l *Logger) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

// Fatal 异步输出一条 Fatal 级别的日志
func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}

// Panic 异步输出一条 Panic 级别的日志
func (l *Logger) Panic(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(PANIC, fmt, args...)
}

func (l *Logger) deliverRecordToWriter(level Level, format string, args ...interface{}) {
	if l == nil && l != defaultLogger {
		defaultLogger.deliverRecordToWriter(level, format, args...)
		return
	}
	l.logWriter.deliverRecord(&l.options, level, format, args...)
}

// GetLogger 获取当前 Logger 的 Logger ，意义在于会进行接收器 Logger 是否为空的判断，
// 如果为空，底层默认会使用 defaultLogger 操作，因此返回 defaultLogger 。
func (l *Logger) GetLogger() *Logger {
	if l == nil && l != defaultLogger {
		return defaultLogger.GetLogger()
	}
	return l
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields ...Field) *Logger {
	if len(fields) == 0 {
		return l
	}
	res := l.Clone()
	// res.core = l.core.With(fields)
	return res
}

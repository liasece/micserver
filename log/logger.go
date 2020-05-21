package log

import (
	"fmt"
	"strings"
	"time"
)

// Logger 日志实例
type Logger struct {
	logWriter   *RecordWriter
	level       int32
	logname     string
	lastTime    int64
	lastTimeStr string
	layout      string
	topic       string
}

// NewLogger 构造一个日志
func NewLogger(isDaemon bool, logFilePath string) *Logger {
	l := new(Logger)
	l.level = SYS
	l.layout = "060102-15:04:05"
	l.logWriter = &RecordWriter{}
	l.logWriter.Init()

	if len(logFilePath) != 0 {
		if isDaemon {
			l.logWriter.AddLogFile(logFilePath, true)
			l.logWriter.RemoveConsoleLog()
		} else {
			// 默认走控制台
			l.logWriter.AddLogFile(logFilePath, false)
			w := NewConsoleWriter()
			w.SetColor(true)
			l.logWriter.registerLogWriter(w)
		}
	} else {
		w := NewConsoleWriter()
		w.SetColor(true)
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
	l.logname = logname
}

// GetLogWriter 获取日志写入器
func (l *Logger) GetLogWriter() *RecordWriter {
	if l == nil && l != defaultLogger {
		return defaultLogger.GetLogWriter()
	}
	return l.logWriter
}

// getLogLevel 获取当前日志等级
func (l *Logger) getLogLevel() int32 {
	if l == nil && l != defaultLogger {
		return defaultLogger.getLogLevel()
	}
	return l.level
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
func (l *Logger) SetLogLevel(lvl int32) {
	if l == nil && l != defaultLogger {
		defaultLogger.SetLogLevel(lvl)
		return
	}
	l.level = lvl
}

// SetLogLevelByStr 使用等级名设置日志等级
func (l *Logger) SetLogLevelByStr(lvl string) error {
	if l == nil && l != defaultLogger {
		return defaultLogger.SetLogLevelByStr(lvl)
	}
	lvlUpper := strings.ToUpper(lvl)
	switch lvlUpper {
	case "SYS":
		l.level = SYS
	case "DEBUG":
		l.level = DEBUG
	case "INFO":
		l.level = INFO
	case "WARNING":
		l.level = WARNING
	case "ERROR":
		l.level = ERROR
	case "FATAL":
		l.level = FATAL
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
	l.topic = topic
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

func (l *Logger) deliverRecordToWriter(level int32, format string, args ...interface{}) {
	if l == nil && l != defaultLogger {
		defaultLogger.deliverRecordToWriter(level, format, args...)
		return
	}
	var inf, code string
	// 检查日志等级有效性
	if level < l.level {
		return
	}
	// 连接主题
	if l.topic != "" {
		inf += l.topic + " "
	}
	// 连接格式化内容
	if format != "" {
		inf += fmt.Sprintf(format, args...)
	} else {
		inf += fmt.Sprint(args...)
	}
	// format time
	now := time.Now()
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}
	// record to recorder
	r := recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = l.lastTimeStr
	r.level = level
	r.name = l.logname
	r.timeUnix = l.lastTime

	l.logWriter.write(r)
}

// GetLogger 获取当前 Logger 的 Logger ，意义在于会进行接收器 Logger 是否为空的判断，
// 如果为空，底层默认会使用 defaultLogger 操作，因此返回 defaultLogger 。
func (l *Logger) GetLogger() *Logger {
	if l == nil && l != defaultLogger {
		return defaultLogger.GetLogger()
	}
	return l
}

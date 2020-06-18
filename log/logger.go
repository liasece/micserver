package log

import (
	"fmt"
	"strings"
	"time"
)

// Logger 日志实例
type Logger struct {
	logWriter   *recordWriter
	lastTime    int64
	lastTimeStr string
	opt         Options
}

// Options of log
type Options struct {
	// NoConsole while remove console out put, default false.
	NoConsole bool
	// NoConsoleColor whil disable console output color, default false.
	NoConsoleColor bool
	// FilePaths is the log output file path, default none log file.
	FilePaths []string
	// RecordTimeLayout use for (time.Time).Format(layout string) record time field, default "060102-15:04:05",
	// will not be empty.
	RecordTimeLayout string
	// Level log record level limit, only higer thie level log can be get reach Writer, default SYS.
	Level Level
	// Name is thie logger name, default "".
	Name string
	// Topic is thie logger topic, default "".
	Topic string
	// AsyncWrite while asynchronously output the log record to Write, it may be more performance,
	// but if you exit(e.g. os.Exit(1), main() return) this process, it may be loss some log record,
	// because they didn't have time to Write and flush to file.
	AsyncWrite bool
	// AsyncWriteDuration only effective when AsyncWrite is true, this is time duration of asynchronously
	// check log output to write default 100ms.
	AsyncWriteDuration time.Duration
	// RedirectError duplicate stderr to log file, it will be call syscall.Dup2 in linux or syscall.DuplicateHandle
	// in windows, default false.
	RedirectError bool
}

var defaultOptions = Options{
	RecordTimeLayout: "060102-15:04:05",
}

// Check options
func (o *Options) Check() {
	if o.RecordTimeLayout == "" {
		o.RecordTimeLayout = "060102-15:04:05"
	}
	if o.AsyncWrite && o.AsyncWriteDuration.Milliseconds() == 0 {
		o.AsyncWriteDuration = time.Millisecond * 100
	}
}

// NewLogger 构造一个日志
func NewLogger(optp *Options) *Logger {
	opt := defaultOptions
	if optp != nil {
		opt = *optp
		(&opt).Check()
	}
	l := new(Logger)
	l.opt = opt
	l.logWriter = &recordWriter{}
	l.logWriter.Init(&l.opt)

	for _, path := range l.opt.FilePaths {
		l.logWriter.AddLogFile(path)
	}

	if !l.opt.NoConsole {
		w := newConsoleWriter()
		w.SetColor(!l.opt.NoConsoleColor)
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
	l.opt.Name = logname
}

// getLogLevel 获取当前日志等级
func (l *Logger) getLogLevel() Level {
	if l == nil && l != defaultLogger {
		return defaultLogger.getLogLevel()
	}
	return l.opt.Level
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
	l.opt.Level = lvl
}

// SetLogLevelByStr 使用等级名设置日志等级
func (l *Logger) SetLogLevelByStr(lvl string) error {
	if l == nil && l != defaultLogger {
		return defaultLogger.SetLogLevelByStr(lvl)
	}
	lvlUpper := strings.ToUpper(lvl)
	switch lvlUpper {
	case "SYS":
		l.opt.Level = SYS
	case "DEBUG":
		l.opt.Level = DEBUG
	case "INFO":
		l.opt.Level = INFO
	case "WARNING":
		l.opt.Level = WARNING
	case "ERROR":
		l.opt.Level = ERROR
	case "FATAL":
		l.opt.Level = FATAL
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
	l.opt.Topic = topic
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

func (l *Logger) deliverRecordToWriter(level Level, format string, args ...interface{}) {
	if l == nil && l != defaultLogger {
		defaultLogger.deliverRecordToWriter(level, format, args...)
		return
	}
	var inf, code string
	// 检查日志等级有效性
	if level < l.opt.Level {
		return
	}
	// 连接主题
	if l.opt.Topic != "" {
		inf += l.opt.Topic + " "
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
		l.lastTimeStr = now.Format(l.opt.RecordTimeLayout)
	}
	// record to recorder
	r := recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = l.lastTimeStr
	r.level = level
	r.name = l.opt.Name
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

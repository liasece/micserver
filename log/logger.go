package log

import (
	"fmt"
	"strings"
	"time"
)

type Logger struct {
	logWriter   *LogWriter
	level       int32
	logname     string
	lastTime    int64
	lastTimeStr string
	layout      string
	topic       string
}

func NewLogger(isDaemon bool, logFilePath string) *Logger {
	l := new(Logger)
	l.level = SYS
	l.layout = "060102-15:04:05"
	l.logWriter = &LogWriter{}
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

// 浅拷贝出一个 Logger ，他们具备相同的底层写入接口
func (l *Logger) Clone() *Logger {
	if l == nil && l != default_logger {
		return default_logger.Clone()
	}
	res := *l
	return &res
}

func (l *Logger) SetLogName(logname string) {
	if l == nil && l != default_logger {
		default_logger.SetLogName(logname)
		return
	}
	l.logname = logname
}

func (l *Logger) GetLogWriter() *LogWriter {
	if l == nil && l != default_logger {
		return default_logger.GetLogWriter()
	}
	return l.logWriter
}

func (l *Logger) getLogLevel() int32 {
	if l == nil && l != default_logger {
		return default_logger.getLogLevel()
	}
	return l.level
}

func (l *Logger) IsSyslogEnable() bool {
	return l.getLogLevel() >= SYS
}

func (l *Logger) IsDebugEnable() bool {
	return l.getLogLevel() >= DEBUG
}

func (l *Logger) IsInfoEnable() bool {
	return l.getLogLevel() >= DEBUG
}

func (l *Logger) IsWarnEnable() bool {
	return l.getLogLevel() >= WARNING
}

func (l *Logger) IsErrorEnable() bool {
	return l.getLogLevel() >= ERROR
}

func (l *Logger) IsFatalEnable() bool {
	return l.getLogLevel() >= FATAL
}

func (l *Logger) SetLogLevel(lvl int32) {
	if l == nil && l != default_logger {
		default_logger.SetLogLevel(lvl)
		return
	}
	l.level = lvl
}

func (l *Logger) SetLogLevelByStr(lvl string) error {
	if l == nil && l != default_logger {
		return default_logger.SetLogLevelByStr(lvl)
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
		return fmt.Errorf("unknown log level '%s'", lvl)
	}
	return nil
}

func (l *Logger) SetTopic(topic string) error {
	if l == nil || l == default_logger {
		// nil or default logger can't set topic
		return ErrNilLogger
	}
	l.topic = topic
	return nil
}

func (l *Logger) Syslog(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(SYS, fmt, args...)
}

func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

func (l *Logger) Info(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

func (l *Logger) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}

func (l *Logger) deliverRecordToWriter(level int32, format string, args ...interface{}) {
	if l == nil && l != default_logger {
		default_logger.deliverRecordToWriter(level, format, args...)
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

func (l *Logger) GetLogger() *Logger {
	if l == nil && l != default_logger {
		return default_logger.GetLogger()
	}
	return l
}

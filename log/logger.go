package log

import (
	"fmt"
	"path/filepath"
	"time"
)

type Logger struct {
	logWriter   *logWriter
	level       int32
	logname     string
	lastTime    int64
	lastTimeStr string
	layout      string
	topic       string
}

func NewLogger(settings map[string]string) *Logger {
	l := new(Logger)
	l.level = DEBUG
	l.layout = "060102-15:04:05"
	l.logWriter = &logWriter{}
	l.logWriter.Init()

	isDaemon := false
	if v, ok := settings["isdaemon"]; ok {
		if v == "true" {
			isDaemon = true
		}
	}
	logfilename := ""
	if v, ok := settings["logfilename"]; ok {
		logfilename = v
	}
	if v, ok := settings["logpath"]; ok && len(logfilename) != 0 {
		logfile := filepath.Join(v, logfilename)
		if isDaemon {
			l.AddLogFile(logfile, true)
			l.RemoveConsoleLog()
		} else {
			// 默认走控制台
			l.AddLogFile(logfile, false)
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

// 添加这个 Logger 及其同一父节点 Logger 的日志文件
func (l *Logger) AddLogFile(filename string, redirecterr bool) {
	if l == nil && l != default_logger {
		default_logger.AddLogFile(filename, redirecterr)
		return
	}
	l.logWriter.addLogFile(filename, redirecterr)
}

// 修改这个 Logger 及其同一父节点 Logger 的日志文件
func (l *Logger) ChangeLogFile(filename string) {
	if l == nil && l != default_logger {
		default_logger.ChangeLogFile(filename)
		return
	}
	l.logWriter.changeLogFile(filename)
}

// 移除这个 Logger 及其同一父节点 Logger 的控制台log
func (l *Logger) RemoveConsoleLog() {
	if l == nil && l != default_logger {
		default_logger.RemoveConsoleLog()
		return
	}
	l.logWriter.removeConsoleLog()
}

func (l *Logger) SetLogName(logname string) {
	if l == nil && l != default_logger {
		default_logger.SetLogName(logname)
		return
	}
	l.logname = logname
}

func (l *Logger) SetLogLevelByStr(loglevel string) {
	switch loglevel {
	case "debug":
		l.SetLogLevel(DEBUG)
	case "info":
		l.SetLogLevel(INFO)
	case "warning":
		l.SetLogLevel(WARNING)
	case "error":
		l.SetLogLevel(ERROR)
	case "fatal":
		l.SetLogLevel(FATAL)
	default:
		//errors.New("Invalid log level")
	}
}

func (l *Logger) GetLogLevel() int32 {
	if l == nil && l != default_logger {
		return default_logger.GetLogLevel()
	}
	return l.level
}

func (l *Logger) SetLogLevel(lvl int32) {
	if l == nil && l != default_logger {
		default_logger.SetLogLevel(lvl)
		return
	}
	l.level = lvl
}

func (l *Logger) SetLogLayout(layout string) {
	if l == nil && l != default_logger {
		default_logger.SetLogLayout(layout)
		return
	}
	l.layout = layout
}

func (l *Logger) SetTopic(topic string) {
	if l == nil && l != default_logger {
		default_logger.SetTopic(topic)
		return
	}
	l.topic = topic
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

func (l *Logger) CloseLogger() {
	if l == nil && l != default_logger {
		default_logger.CloseLogger()
		return
	}
	l.logWriter.close()
}

func (l *Logger) deliverRecordToWriter(level int32, format string, args ...interface{}) {
	var inf, code string
	if l == nil && l != default_logger {
		default_logger.deliverRecordToWriter(level, format, args...)
		return
	}

	if level < l.level {
		return
	}

	if l.topic != "" {
		inf += "[" + l.topic + "] "
	}

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

	r := recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = l.lastTimeStr
	r.level = level
	r.name = l.logname

	l.logWriter.write(r)
}

func (l *Logger) GetLogger() *Logger {
	if l == nil && l != default_logger {
		return default_logger.GetLogger()
	}
	return l
}

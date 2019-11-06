package log

import (
	"fmt"
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
	l.level = DEBUG
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

// func (l *Logger) SetLogLayout(layout string) {
// 	if l == nil && l != default_logger {
// 		default_logger.SetLogLayout(layout)
// 		return
// 	}
// 	l.layout = layout
// }

func (l *Logger) SetTopic(topic string) {
	if l == nil || l == default_logger {
		// nil or default logger can't set topic
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

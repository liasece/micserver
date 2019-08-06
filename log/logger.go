package log

import (
	"fmt"
	syslog "log"
	"path/filepath"
	"reflect"
	"time"
)

type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int32
	logname     string
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
	stopchan    chan struct{}
}

func NewLogger(settings map[string]string) *Logger {
	l := new(Logger)
	l.writers = make([]Writer, 0, 2)
	l.tunnel = make(chan *Record, tunnel_size_default)
	l.stopchan = make(chan struct{})
	l.c = make(chan bool, 1)
	l.level = DEBUG
	l.layout = "060102-15:04:05"

	go l.boostrapLogWriter()

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
			l.RegisterLogWriter(w)
		}
	} else {
		w := NewConsoleWriter()
		w.SetColor(true)
		l.RegisterLogWriter(w)
	}

	return l
}

func (l *Logger) AddLogFile(filename string, redirecterr bool) {
	if l == nil && l != default_logger {
		default_logger.AddLogFile(filename, redirecterr)
		return
	}
	//	fmt.Printf("log filename,%s \n", filename)
	filebasename := filename
	filename += ".%Y%M%D-%H"
	w := NewFileWriter()
	if redirecterr {
		w.Redirecterr = true
	}
	err := w.SetPathPattern(filebasename, filename)
	if err != nil {
	}
	l.RegisterLogWriter(w)
}

func (l *Logger) ChangeLogFile(filename string) {
	if l == nil && l != default_logger {
		default_logger.ChangeLogFile(filename)
		return
	}
	filebasename := filename
	filename += ".%Y%M%D-%H"
	for i := 0; i < len(l.writers); i++ {
		w := l.writers[i]
		if reflect.TypeOf(w).String() == "*log.FileWriter" {
			if r, ok := w.(Rotater); ok {
				err := r.SetPathPattern(filebasename, filename)
				if err != nil {
				}
				err = r.Rotate()
				if err != nil {
				}
			}
		}
	}
}

func (l *Logger) RemoveConsoleLog() {
	if l == nil && l != default_logger {
		default_logger.RemoveConsoleLog()
		return
	}
	newlist := make([]Writer, 0, 2)
	for i := 0; i < len(l.writers); i++ {
		w := l.writers[i]
		//		Debug("start RemoveConsoleLog, %s", reflect.TypeOf(w).String())
		if reflect.TypeOf(w).String() != "*log.ConsoleWriter" {
			newlist = append(newlist, w)
		}
	}
	l.writers = newlist
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

func (l *Logger) RegisterLogWriter(w Writer) {
	if l == nil && l != default_logger {
		default_logger.RegisterLogWriter(w)
		return
	}
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
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
	select {
	case <-l.stopchan:
		return
	default:
		close(l.stopchan)
		close(l.tunnel)
		break
	}
	select {
	case <-l.c:
		break
	}

	for _, w := range l.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				syslog.Println(err)
			}
		}
	}
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

	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
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

	select {
	case <-l.stopchan:
		return
	default:
	}
	select {
	case <-l.stopchan:
		break
	case l.tunnel <- r:
		break
	}
}

func (l *Logger) boostrapLogWriter() {
	if l == nil {
		if l != default_logger {
			default_logger.Error("[Logger.boostrapLogWriter] l==nil")
		}
		return
	}
	var (
		r  *Record
		ok bool
	)

	if r, ok = <-l.tunnel; !ok {
		l.c <- true
		return
	}

	for _, w := range l.writers {
		if err := w.Write(r); err != nil {
			syslog.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Millisecond * 500)
	//	rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case r, ok = <-l.tunnel:
			if !ok {
				l.c <- true
				return
			}

			for _, w := range l.writers {
				if err := w.Write(r); err != nil {
					syslog.Println(err)
				}
			}

			recordPool.Put(r)

		case <-flushTimer.C:
			for _, w := range l.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						syslog.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 1000)

		case <-rotateTimer.C:
			//	fmt.Printf("start rotate file,actions, 1111\n")
			for _, w := range l.writers {
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						syslog.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}

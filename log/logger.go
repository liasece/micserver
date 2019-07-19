package log

import (
	"fmt"
	syslog "log"
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
}

func NewLogger(settings map[string]string) *Logger {
	l := new(Logger)
	l.SetLogName("log")
	l.writers = make([]Writer, 0, 2)
	l.tunnel = make(chan *Record, tunnel_size_default)
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
	if v, ok := settings["logpath"]; ok {
		if isDaemon {
			l.AddlogFile(v, true)
			l.RemoveConsoleLog()
			l.Debug("Logger is start as a daemon")
		} else {
			l.AddlogFile(v, false)
		}
	}

	return l
}

func (l *Logger) AddlogFile(filename string, redirecterr bool) {
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
	l.Register(w)
}

func (l *Logger) ChangelogFile(filename string) {
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
	Debug("start RemoveConsoleLog")
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
	l.logname = logname
}

func (l *Logger) SetLevelByStr(loglevel string) {
	switch loglevel {
	case "debug":
		l.SetLevel(DEBUG)
	case "info":
		l.SetLevel(INFO)
	case "warning":
		l.SetLevel(WARNING)
	case "error":
		l.SetLevel(ERROR)
	case "fatal":
		l.SetLevel(FATAL)
	default:
		//errors.New("Invalid log level")
	}
}

func (l *Logger) GetLevel() int32 {
	return l.level
}

func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func (l *Logger) SetLevel(lvl int32) {
	l.level = lvl
}

func (l *Logger) SetLayout(layout string) {
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

func (l *Logger) Close() {
	close(l.tunnel)
	<-l.c

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

	l.tunnel <- r
}

func (l *Logger) boostrapLogWriter() {
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

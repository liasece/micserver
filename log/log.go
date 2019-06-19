package log

import (
	"fmt"
	syslog "log"
	"reflect"
	//	"path"
	//	"runtime"
	//	"strconv"
	"sync"
	"time"
)

var (
	LEVEL_FLAGS = [...]string{"DEBUG", " INFO", " WARN", "ERROR", "FATALERROR"}
	recordPool  *sync.Pool
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const tunnel_size_default = 1024

type Record struct {
	time  string
	name  string
	code  string
	info  string
	level int32
}

func AutoConfig(logfilename string, logname string, daemon bool) {
	SetLogName(logname)
	if daemon {
		AddlogFile(logfilename, true)
		RemoveConsoleLog()
		Debug("[log] Program is start as a daemon")
	} else {
		AddlogFile(logfilename, false)
	}
}

func (r *Record) String() string {
	return fmt.Sprintf("%s [%s] %s: %s\n", r.time, r.name, LEVEL_FLAGS[r.level], r.info)
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Rotater interface {
	Rotate() error
	SetPathPattern(string, string) error
}

type Flusher interface {
	Flush() error
}

type log struct {
	writers     []Writer
	tunnel      chan *Record
	level       int32
	logname     string
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
}

func Newlog() *log {
	if log_default != nil && !takeup {
		takeup = true
		return log_default
	}

	l := new(log)
	l.writers = make([]Writer, 0, 2)
	l.tunnel = make(chan *Record, tunnel_size_default)
	l.c = make(chan bool, 1)
	l.level = DEBUG
	l.layout = "060102-15:04:05"

	go boostrapLogWriter(l)

	return l
}

func (l *log) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func (l *log) SetLevel(lvl int32) {
	l.level = lvl
}

func (l *log) SetLayout(layout string) {
	l.layout = layout
}

func (l *log) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

func (l *log) Warn(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

func (l *log) Info(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

func (l *log) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

func (l *log) Fatal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}

func (l *log) Close() {
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

func (l *log) deliverRecordToWriter(level int32, format string, args ...interface{}) {
	var inf, code string

	if level < l.level {
		return
	}

	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	/*
		// source code, file and line num
		_, file, line, ok := runtime.Caller(2)
		if ok {
			code = path.Base(file) + ":" + strconv.Itoa(line)
		}
	*/

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

func boostrapLogWriter(log *log) {
	if log == nil {
		panic("log is nil")
	}

	var (
		r  *Record
		ok bool
	)

	if r, ok = <-log.tunnel; !ok {
		log.c <- true
		return
	}

	for _, w := range log.writers {
		if err := w.Write(r); err != nil {
			syslog.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Millisecond * 500)
	//	rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case r, ok = <-log.tunnel:
			if !ok {
				log.c <- true
				return
			}

			for _, w := range log.writers {
				if err := w.Write(r); err != nil {
					syslog.Println(err)
				}
			}

			recordPool.Put(r)

		case <-flushTimer.C:
			for _, w := range log.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						syslog.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 1000)

		case <-rotateTimer.C:
			//	fmt.Printf("start rotate file,actions, 1111\n")
			for _, w := range log.writers {
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

// default
var (
	log_default *log
	takeup      = false
)

func SetLevel(lvl int32) {
	log_default.level = lvl
}

func GetLevel() int32 {
	return log_default.level
}

func SetLayout(layout string) {
	log_default.layout = layout
}

func Debug(fmt string, args ...interface{}) {
	log_default.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	log_default.deliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	log_default.deliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	log_default.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	log_default.deliverRecordToWriter(FATAL, fmt, args...)
}

func Register(w Writer) {
	log_default.Register(w)
}

func Close() {
	log_default.Close()
}

func init() {
	fmt.Printf("log init \n")
	log_default = Newlog()
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
	SetLogName("SS")
	// 默认走控制台
	w := NewConsoleWriter()
	w.SetColor(false)
	Register(w)
}
func SetLogName(logname string) {
	log_default.logname = logname
}

func SetLogLevel(loglevel string) {
	switch loglevel {
	case "debug":
		SetLevel(DEBUG)
	case "info":
		SetLevel(INFO)
	case "warning":
		SetLevel(WARNING)
	case "error":
		SetLevel(ERROR)
	case "fatal":
		SetLevel(FATAL)
	default:
		//errors.New("Invalid log level")
	}
}

func AddlogFile(filename string, redirecterr bool) {
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
	Register(w)
}
func ChangelogFile(filename string) {
	filebasename := filename
	filename += ".%Y%M%D-%H"
	for i := 0; i < len(log_default.writers); i++ {
		w := log_default.writers[i]
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

func RemoveConsoleLog() {
	Debug("start RemoveConsoleLog")
	newlist := make([]Writer, 0, 2)
	for i := 0; i < len(log_default.writers); i++ {
		w := log_default.writers[i]
		//		Debug("start RemoveConsoleLog, %s", reflect.TypeOf(w).String())
		if reflect.TypeOf(w).String() != "*log.ConsoleWriter" {
			newlist = append(newlist, w)
		}
	}
	log_default.writers = newlist
}

package log

import (
	"fmt"
	"sync"
	"time"
)

var (
	LEVEL_FLAGS = [...]string{"[S]", "[D]", "[I]", "[WARNING]", "[ERROR]", "[FATALERROR]"}
	recordPool  *sync.Pool
)

const (
	SYS = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

const tunnel_size_default = 1024

type Record struct {
	time     string
	name     string
	code     string
	info     string
	level    int32
	timeUnix int64
}

func (r *Record) String() string {
	return fmt.Sprintf("%s [%s] %s %s\n", r.time, r.name, LEVEL_FLAGS[r.level], r.info)
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Rotater interface {
	Rotate() error
	RotateByTime(*time.Time) error
	SetPathPattern(string, string) error
}

type Flusher interface {
	Flush() error
}

// default
var (
	default_logger *Logger
	takeup         = false
)

func Syslog(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(SYS, fmt, args...)
}

func Debug(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	default_logger.deliverRecordToWriter(FATAL, fmt, args...)
}

func SetDefaultLogger(l *Logger) {
	if l != default_logger {
		default_logger.GetLogWriter().Close()
		default_logger = l
	}
}

func GetDefaultLogger() *Logger {
	return default_logger
}

func init() {
	default_logger = NewLogger(false, "")
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
	default_logger.SetLogName("log")
}

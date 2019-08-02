package log

import (
	"fmt"
	"sync"
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

// default
var (
	default_logger *Logger
	takeup         = false
)

func SetLevel(lvl int32) {
	default_logger.level = lvl
}

func GetLevel() int32 {
	return default_logger.level
}

func SetLayout(layout string) {
	default_logger.layout = layout
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

func CloseLogger() {
	default_logger.CloseLogger()
}

func SetDefaultLogger(l *Logger) {
	if l != default_logger {
		default_logger.CloseLogger()
		default_logger = l
	}
}

func GetDefaultLogger() *Logger {
	return default_logger
}

func init() {
	default_logger = NewLogger(nil)
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
	default_logger.SetLogName("log")
}

func SetLogName(logname string) {
	default_logger.logname = logname
}

func SetLogLevel(loglevel string) {
	default_logger.SetLogLevelByStr(loglevel)
}

func AddLogFile(filename string, redirecterr bool) {
	default_logger.AddLogFile(filename, redirecterr)
}

func ChangeLogFile(filename string) {
	default_logger.ChangeLogFile(filename)
}

func RemoveConsoleLog() {
	default_logger.RemoveConsoleLog()
}

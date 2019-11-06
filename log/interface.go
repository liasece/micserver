package log

type ILogger interface {
	Debug(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})
	Fatal(fmt string, args ...interface{})
	Clone() *Logger
	SetTopic(topic string)
	SetLogLevel(lvl int32)
	GetLogLevel() int32
	SetLogName(logname string)
	GetLogger() *Logger
}

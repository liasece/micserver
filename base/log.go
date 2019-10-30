package base

type ILog interface {
	Debug(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})
	Fatal(fmt string, args ...interface{})
	Clone() ILog
	Close()
	SetLogName(logname string)
	SetTopic(topic string)
}

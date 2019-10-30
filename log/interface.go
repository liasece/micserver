package log

type ILogger interface {
	Debug(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})
	Fatal(fmt string, args ...interface{})
	Clone() *Logger
	CloseLogger()
	SetTopic(topic string)
	SetLogLayout(layout string)
	SetLogLevel(lvl int32)
	GetLogLevel() int32
	SetLogLevelByStr(loglevel string)
	SetLogName(logname string)
	RemoveConsoleLog()
	ChangeLogFile(filename string)
	AddLogFile(filename string, redirecterr bool)
	GetLogger() *Logger
}

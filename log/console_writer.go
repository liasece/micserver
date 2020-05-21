package log

import (
	"fmt"
	"os"
)

// 一条日志记录
type colorRecord Record

// String 获取该记录在控制台中携带颜色的格式化字符串
func (r *colorRecord) String() string {
	switch r.level {
	case SYS:
		return fmt.Sprintf("\033[36m%s\033[0m [%s] %s %s\n",
			r.time, r.name, LEVELFLAGS[r.level], r.info)
	case DEBUG:
		return fmt.Sprintf("\033[36m%s\033[0m [%s] \033[34m%s\033[0m %s\n",
			r.time, r.name, LEVELFLAGS[r.level], r.info)
	case INFO:
		return fmt.Sprintf("\033[36m%s\033[0m [%s] \033[32m%s\033[0m %s\n",
			r.time, r.name, LEVELFLAGS[r.level], r.info)
	case WARNING:
		return fmt.Sprintf("\033[36m%s\033[0m [%s] \033[33m%s\033[0m %s\n",
			r.time, r.name, LEVELFLAGS[r.level], r.info)
	case ERROR:
		return fmt.Sprintf("\033[36m%s\033[0m [%s] \033[31m%s\033[0m %s\n",
			r.time, r.name, LEVELFLAGS[r.level], r.info)
	case FATAL:
		return fmt.Sprintf("\033[36m%s\033[0m [%s] \033[35m%s\033[0m %s\n",
			r.time, r.name, LEVELFLAGS[r.level], r.info)
	default:
		return fmt.Sprintf("\033[36m%s\033[0m [%s] \033[35mUNKNOW\033[0m %s\n",
			r.time, r.name, r.info)
	}
}

// ConsoleWriter 控制台输出器
type ConsoleWriter struct {
	color bool
}

// NewConsoleWriter 构造一个控制台输出器
func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

// Write 写入一条日志记录到控制台
func (w *ConsoleWriter) Write(r *Record) error {
	if w.color {
		fmt.Fprint(os.Stdout, ((*colorRecord)(r)).String())
	} else {
		fmt.Fprint(os.Stdout, r.String())
	}
	return nil
}

// Init 初始化控制台输出器
func (w *ConsoleWriter) Init() error {
	return nil
}

// GetType 获取输出器的类型 返回控制台类型
func (w *ConsoleWriter) GetType() WriterType {
	return writerTypeConsole
}

// SetColor 设置该输出器在控制台中是否携带颜色
func (w *ConsoleWriter) SetColor(c bool) {
	w.color = c
}

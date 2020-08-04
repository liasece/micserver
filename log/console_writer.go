package log

import (
	"fmt"

	"github.com/liasece/micserver/log/core"
)

// 一条日志记录
type colorRecord Record

// String 获取该记录在控制台中携带颜色的格式化字符串
func (r *colorRecord) String() string {
	flag := _levelToCapitalColorString[SysLevel]
	if r.level >= SysLevel && r.level <= FatalLevel {
		flag = _levelToCapitalColorString[r.level]
	}

	if r.name == "" {
		return fmt.Sprintf("\033[36m%s\033[0m "+flag+"\033[0m %s",
			r.time, r.info)
	}
	return fmt.Sprintf("\033[36m%s\033[0m [%s] "+flag+"\033[0m %s",
		r.time, r.name, r.info)
}

// consoleWriter 控制台输出器
type consoleWriter struct {
	encoder core.Encoder
	color   bool
}

// newConsoleWriter 构造一个控制台输出器
func newConsoleWriter() *consoleWriter {
	return &consoleWriter{
		encoder: core.NewJSONEncoder(core.EncoderConfig{}),
	}
}

// Write 写入一条日志记录到控制台
func (w *consoleWriter) Write(r *Record) error {
	fieldStr := ""
	if len(r.fields) > 0 {
		buf, _ := w.encoder.EncodeEntry(nil, r.fields)
		fieldStr = fmt.Sprintf(" %s", buf.String())
	}

	if w.color {
		fmt.Println(((*colorRecord)(r)).String() + fieldStr)
	} else {
		fmt.Println(r.String() + fieldStr)
	}
	return nil
}

// Init 初始化控制台输出器
func (w *consoleWriter) Init() error {
	return nil
}

// SetColor 设置该输出器在控制台中是否携带颜色
func (w *consoleWriter) SetColor(c bool) {
	w.color = c
}

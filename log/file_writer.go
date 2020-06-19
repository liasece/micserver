package log

import (
	"bufio"
	"errors"
	"fmt"
	syslog "log"
	"os"
	"path"
	"time"
)

// fileWriter 文件输出器，将日志记录输出到文件中
type fileWriter struct {
	filePath         string
	lastRoateTail    string
	rotateTimeLayout string
	file             *os.File
	fileBufWriter    *bufio.Writer
	variables        []interface{}
	redirectError    bool // 是否重定向错误信息到日志文件
}

// newFileWriter 构造一个文件输出器
func newFileWriter(filePath string, rotateTimeLayout string) *fileWriter {
	return &fileWriter{
		filePath:         filePath,
		rotateTimeLayout: rotateTimeLayout,
	}
}

// Init 初始化文件输出器
func (w *fileWriter) Init() error {
	return w.Rotate()
}

// Write 写入一条日志记录
func (w *fileWriter) Write(r *Record) error {
	if w.fileBufWriter == nil {
		return errors.New("no opened file")
	}
	if _, err := w.fileBufWriter.WriteString(r.String()); err != nil {
		return err
	}
	return nil
}

// Rotate 尝试转储文件
func (w *fileWriter) Rotate() error {
	now := time.Now()
	return w.doRotate(&now)
}

// RotateByTime 尝试小时转储
func (w *fileWriter) RotateByTime(t *time.Time) error {
	return w.doRotate(t)
}

// doRotate 尝试转储文件，如果不需要进行转储，返回 nil
func (w *fileWriter) doRotate(t *time.Time) error {
	if w.rotateTimeLayout == "" {
		return w.initFile(w.filePath)
	}
	newRotateTail := t.Format(w.rotateTimeLayout)
	if newRotateTail == w.lastRoateTail {
		return nil
	}

	if w.fileBufWriter != nil {
		if err := w.fileBufWriter.Flush(); err != nil {
			return err
		}
	}

	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return err
		}
	}

	filePath := fmt.Sprint(w.filePath, ".", newRotateTail)
	if err := w.initFile(filePath); err != nil {
		syslog.Println(err)
		return err
	}
	w.lastRoateTail = newRotateTail
	return nil
}

func (w *fileWriter) initFile(filePath string) error {
	if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	// 这是真正的日志文件
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		syslog.Println(filePath, err)
		return err
	}
	w.file = file

	if w.redirectError {
		// 把错误重定向到日志文件来
		sysDup(int(w.file.Fd()))
	}

	if w.fileBufWriter = bufio.NewWriterSize(w.file, 81920); w.fileBufWriter == nil {
		return errors.New("new fileBufWriter failed")
	}

	// 创建一个软链接
	if filePath != w.filePath {
		// 检查文件是存在
		_, fileerr := os.Stat(w.filePath)
		if fileerr == nil {
			// this file already exist
			os.Remove(w.filePath)
		}
		// Create a symlink
		err := os.Symlink(path.Base(filePath), w.filePath)
		if err != nil {
			syslog.Println("os.Symlink error:", err.Error())
			// return err
		}
	}
	return nil
}

// Flush 将文件缓冲区中的内容
func (w *fileWriter) Flush() error {
	if w.fileBufWriter != nil {
		return w.fileBufWriter.Flush()
	}
	return nil
}

package log

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"
	"time"
)

var pathVariableTable map[byte]func(*time.Time) int

type FileWriter struct {
	filebasename  string
	pathFmt       string
	file          *os.File
	fileBufWriter *bufio.Writer
	actions       []func(*time.Time) int
	variables     []interface{}
	Redirecterr   bool // 是否重定向错误信息到日志文件
}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

func (w *FileWriter) Init() error {
	return w.Rotate()
}

func (w *FileWriter) SetPathPattern(filebasename string, pattern string) error {
	n := 0
	for _, c := range pattern {
		if c == '%' {
			n++
		}
	}

	if n == 0 {
		w.filebasename = filebasename
		w.pathFmt = pattern
		return nil
	}

	w.actions = make([]func(*time.Time) int, 0, n)
	w.variables = make([]interface{}, n)
	tmp := []byte(pattern)

	variable := 0
	for _, c := range tmp {
		if variable == 1 {
			act, ok := pathVariableTable[c]
			if !ok {
				return errors.New("Invalid rotate pattern (" + pattern + ")")
			}
			w.actions = append(w.actions, act)
			variable = 0
			continue
		}
		if c == '%' {
			variable = 1
		}
	}

	w.filebasename = filebasename
	w.pathFmt = convertPatternToFmt(tmp)

	return nil
}

func (w *FileWriter) Write(r *Record) error {
	if w.fileBufWriter == nil {
		return errors.New("no opened file")
	}
	if _, err := w.fileBufWriter.WriteString(r.String()); err != nil {
		return err
	}
	return nil
}

func (w *FileWriter) Rotate() error {
	now := time.Now()
	v := 0
	rotate := false

	for i, act := range w.actions {
		v = act(&now)
		if v != w.variables[i] {
			w.variables[i] = v
			rotate = true
		}
	}
	//	fmt.Printf("start rotate file,actions:%d,%d,%d,%d,%d\n", len(w.actions), w.variables[0], w.variables[1], w.variables[2], w.variables[3])

	if !rotate {
		return nil
	}
	//	fmt.Printf("start rotate file,actions:%d,%d,%d,%d,%d\n", len(w.actions), w.variables[0], w.variables[1], w.variables[2], w.variables[3])

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

	filePath := fmt.Sprintf(w.pathFmt, w.variables...)

	if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	// 这是真正的日志文件
	if file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return err
	} else {
		w.file = file
	}

	// 创建一个软链接
	// 检查文件是存在
	_, fileerr := os.Stat(w.filebasename)
	if fileerr != nil && os.IsNotExist(fileerr) { // 文件不存在
	}
	os.Remove(w.filebasename)
	// Create a symlink
	err := os.Symlink(filePath, w.filebasename)
	if err != nil {
		// return err
	}

	if w.Redirecterr {
		// 把错误重定向到日志文件来
		e1 := syscall.Dup2(int(w.file.Fd()), 1)
		if e1 != nil {
			// return e1
		}
		e2 := syscall.Dup2(int(w.file.Fd()), 2)
		if e2 != nil {
			// return e2
		}
	}

	if w.fileBufWriter = bufio.NewWriterSize(w.file, 81920); w.fileBufWriter == nil {
		return errors.New("new fileBufWriter failed.")
	}

	return nil
}

func (w *FileWriter) Flush() error {
	if w.fileBufWriter != nil {
		return w.fileBufWriter.Flush()
	}
	return nil
}

func getYear(now *time.Time) int {
	return now.Year() % 100
}

func getMonth(now *time.Time) int {
	return int(now.Month())
}

func getDay(now *time.Time) int {
	return now.Day()
}

func getHour(now *time.Time) int {
	return now.Hour()
}

func getMin(now *time.Time) int {
	return now.Minute()
}

func convertPatternToFmt(pattern []byte) string {
	pattern = bytes.Replace(pattern, []byte("%Y"), []byte("%d"), -1)
	pattern = bytes.Replace(pattern, []byte("%M"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%D"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%H"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%m"), []byte("%02d"), -1)
	return string(pattern)
}

func init() {
	pathVariableTable = make(map[byte]func(*time.Time) int, 5)
	pathVariableTable['Y'] = getYear
	pathVariableTable['M'] = getMonth
	pathVariableTable['D'] = getDay
	pathVariableTable['H'] = getHour
	pathVariableTable['m'] = getMin
}

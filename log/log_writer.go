package log

import (
	syslog "log"
	"time"
)

const tunnelSizeDefault = 1024

// WriterType writer type
type WriterType int

const (
	writerTypeConsole = 1
	writerTypeFile    = 2
)

// RecordWriter log写入器
type RecordWriter struct {
	writers  []Writer
	tunnel   chan *Record
	c        chan bool
	stopchan chan struct{}
}

// Init 初始化log写入器
func (recordWriter *RecordWriter) Init() {
	recordWriter.writers = make([]Writer, 0, 2)
	recordWriter.tunnel = make(chan *Record, tunnelSizeDefault)
	recordWriter.stopchan = make(chan struct{})
	recordWriter.c = make(chan bool, 1)

	go recordWriter.boostrapLogWriter()
}

// AddLogFile 增加一个文件输出器
func (recordWriter *RecordWriter) AddLogFile(filename string, redirecterr bool) {
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
	recordWriter.registerLogWriter(w)
}

// ChangeLogFile 修改所有的文件输出器的目标文件
func (recordWriter *RecordWriter) ChangeLogFile(filename string) {
	filebasename := filename
	filename += ".%Y%M%D-%H"
	for i := 0; i < len(recordWriter.writers); i++ {
		w := recordWriter.writers[i]
		if w.GetType() == writerTypeFile {
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

// RemoveConsoleLog 移除控制台输出器
func (recordWriter *RecordWriter) RemoveConsoleLog() {
	newlist := make([]Writer, 0, 2)
	for i := 0; i < len(recordWriter.writers); i++ {
		w := recordWriter.writers[i]
		if w.GetType() != writerTypeConsole {
			newlist = append(newlist, w)
		}
	}
	recordWriter.writers = newlist
}

// registerLogWriter 注册一个文件输出器到该 log 写入器中
func (recordWriter *RecordWriter) registerLogWriter(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	recordWriter.writers = append(recordWriter.writers, w)
}

// Close 关闭当前写入器的所有输出器
func (recordWriter *RecordWriter) Close() {
	select {
	case <-recordWriter.stopchan:
		return
	default:
		close(recordWriter.stopchan)
		// close(recordWriter.tunnel)
		break
	}
	select {
	case <-recordWriter.c:
		break
	}

	for _, w := range recordWriter.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				syslog.Println(err)
			}
		}
	}
}

// write 写入一条日志记录，等待后续异步处理
func (recordWriter *RecordWriter) write(r *Record) {
	select {
	case <-recordWriter.stopchan:
		return
	default:
	}
	select {
	case <-recordWriter.stopchan:
		break
	case recordWriter.tunnel <- r:
		break
	}
}

// boostrapLogWriter 日志写入线程
func (recordWriter *RecordWriter) boostrapLogWriter() {
	var (
		r  *Record
		ok bool
	)

	if r, ok = <-recordWriter.tunnel; !ok {
		recordWriter.c <- true
		return
	}

	for _, w := range recordWriter.writers {
		if err := w.Write(r); err != nil {
			syslog.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Millisecond * 100)
	//	rotateTimer := time.NewTimer(time.Second * 10)
	lastRecordTimeUnix60 := int64(0)
	lastRecordTime := time.Now()
	for {
		select {
		case r, ok = <-recordWriter.tunnel:
			if !ok {
				recordWriter.c <- true
				return
			}
			needTryRotate := false
			// 如果上一条记录的时间和这条记录的时间不是同一分钟，需要尝试增加log文件
			nowTimeUnix60 := r.timeUnix / 60
			if lastRecordTimeUnix60 != nowTimeUnix60 {
				needTryRotate = true
				lastRecordTime = time.Unix(r.timeUnix, 0)
				// 更新最后一条记录的时间
				lastRecordTimeUnix60 = nowTimeUnix60
			}
			for _, w := range recordWriter.writers {
				if needTryRotate {
					// 需要增加log文件
					if r, ok := w.(Rotater); ok {
						if err := r.RotateByTime(&lastRecordTime); err != nil {
							syslog.Println(err)
						}
					}
				}
				// 写入log
				if err := w.Write(r); err != nil {
					syslog.Println(err)
				}
			}
			recordPool.Put(r)
		case <-recordWriter.stopchan:
			recordWriter.c <- true
			return
		case <-flushTimer.C:
			for _, w := range recordWriter.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						syslog.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 500)
		case <-rotateTimer.C:
			//	fmt.Printf("start rotate file,actions, 1111\n")
			recordWriter.tryRotate()
			rotateTimer.Reset(time.Second * 60)
		}
	}
}

// tryRotate 尝试将文件输出器转储
func (recordWriter *RecordWriter) tryRotate() {
	for _, w := range recordWriter.writers {
		if r, ok := w.(Rotater); ok {
			if err := r.Rotate(); err != nil {
				syslog.Println(err)
			}
		}
	}
}

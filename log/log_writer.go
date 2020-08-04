package log

import (
	"fmt"
	syslog "log"
	"time"
)

const tunnelSizeDefault = 1024

const (
	writerTypeConsole = 1
	writerTypeFile    = 2
)

// recordWriter log写入器
type recordWriter struct {
	wm                   *writerManager
	tunnel               chan *Record
	c                    chan bool
	stopchan             chan struct{}
	lastRecordTimeUnix60 int64
	lastRecordTime       time.Time
	lastTime             int64
	lastTimeStr          string
}

// Init 初始化log写入器
func (rec *recordWriter) Init(opt *options) {
	rec.wm = &writerManager{}
	rec.tunnel = make(chan *Record, tunnelSizeDefault)
	rec.stopchan = make(chan struct{})
	rec.c = make(chan bool, 1)

	go rec.boostrapLogWriter(opt)
}

// AddLogFile 增加一个文件输出器
func (rec *recordWriter) AddLogFile(filePath string, opt *options) error {
	w := newFileWriter(filePath, opt.RotateTimeLayout)
	w.redirectError = opt.RedirectError
	err := w.Init()
	if err != nil {
		syslog.Println(err)
		return err
	}
	rec.registerLogWriter(w)
	return nil
}

// registerLogWriter 注册一个文件输出器到该 log 写入器中
func (rec *recordWriter) registerLogWriter(w Writer) {
	rec.wm.AddWriter(w)
}

// Close 关闭当前写入器的所有输出器
func (rec *recordWriter) Close() {
	select {
	case <-rec.stopchan:
		return
	default:
		close(rec.stopchan)
		break
	}
	select {
	case <-rec.c:
		break
	}

	if err := rec.wm.Flush(); err != nil {
		syslog.Println("Close f.Flush error", err)
	}
}

func (rec *recordWriter) deliverRecord(opt *options, level Level, format string, originArgs ...interface{}) {
	var inf, code string
	// 检查日志等级有效性
	if level < opt.Level {
		return
	}
	// 连接主题
	if opt.Topic != "" {
		inf += opt.Topic + " "
	}

	var fields []Field
	var args []interface{}
	for _, vi := range originArgs {
		if v, ok := vi.(Field); ok {
			fields = append(fields, v)
		} else {
			args = append(args, vi)
		}
	}

	// 连接格式化内容
	if format != "" {
		inf += fmt.Sprintf(format, args...)
	} else {
		inf += fmt.Sprint(args...)
	}
	// format time
	now := time.Now()
	if now.Unix() != rec.lastTime {
		rec.lastTime = now.Unix()
		rec.lastTimeStr = now.Format(opt.RecordTimeLayout)
	}
	// record to recorder
	r := recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = rec.lastTimeStr
	r.level = level
	r.name = opt.Name
	r.timeUnix = rec.lastTime
	r.fields = fields

	rec.write(r, opt)
	if level >= PANIC {
		panic(r)
	}
}

// write 写入一条日志记录，等待后续异步处理
func (rec *recordWriter) write(r *Record, opt *options) {
	select {
	case <-rec.stopchan:
		return
	default:
	}

	// no async
	if !opt.AsyncWrite {
		rec.doWriteRecord(r, opt)
		return
	}

	select {
	case <-rec.stopchan:
		break
	case rec.tunnel <- r:
		break
	}
}

// boostrapLogWriter 日志写入线程
func (rec *recordWriter) boostrapLogWriter(opt *options) {
	var (
		r  *Record
		ok bool
	)

	flushTimer := time.NewTimer(opt.AsyncWriteDuration)
	rotateTimer := time.NewTimer(time.Millisecond * 100)
	for {
		select {
		case r, ok = <-rec.tunnel:
			if !ok {
				rec.c <- true
				return
			}
			rec.doWriteRecord(r, opt)
		case <-rec.stopchan:
			rec.c <- true
			return
		case <-flushTimer.C:
			if opt.AsyncWrite {
				err := rec.wm.Flush()
				if err != nil {
					syslog.Println(err)
				}
				flushTimer.Reset(opt.AsyncWriteDuration)
			}
		case <-rotateTimer.C:
			rec.tryRotate()
			rotateTimer.Reset(time.Second * 60)
		}
	}
}

func (rec *recordWriter) doWriteRecord(r *Record, opt *options) error {
	needTryRotate := false
	// 如果上一条记录的时间和这条记录的时间不是同一分钟，需要尝试增加log文件
	nowTimeUnix60 := r.timeUnix / 60
	if rec.lastRecordTimeUnix60 != nowTimeUnix60 {
		needTryRotate = true
		rec.lastRecordTime = time.Unix(r.timeUnix, 0)
		// 更新最后一条记录的时间
		rec.lastRecordTimeUnix60 = nowTimeUnix60
	}
	if needTryRotate {
		// 需要增加log文件
		if err := rec.wm.RotateByTime(&rec.lastRecordTime); err != nil {
			syslog.Println("r.RotateByTime error", err)
			return err
		}
	}
	// 写入log
	if err := rec.wm.Write(r); err != nil {
		syslog.Println("doWriteRecord for w.Write error", err)
		return err
	}
	recordPool.Put(r)

	// flush
	if !opt.AsyncWrite {
		if err := rec.wm.Flush(); err != nil {
			syslog.Println("doWriteRecord f.Flush error", err)
			return err
		}
	}
	return nil
}

// tryRotate 尝试将文件输出器转储
func (rec *recordWriter) tryRotate() {
	if err := rec.wm.Rotate(); err != nil {
		syslog.Println("r.Rotate error", err)
	}
}

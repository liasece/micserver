package log

import (
	"errors"
	syslog "log"
	"strings"
	"sync"
	"time"
)

const tunnelSizeDefault = 1024

const (
	writerTypeConsole = 1
	writerTypeFile    = 2
)

// writerCfg config of log writer
type writerCfg struct {
	w Writer
	f Flusher
}

// writerManager writer manager
type writerManager struct {
	ws []*writerCfg
	l  sync.Mutex
}

func (wm *writerManager) AddWriter(w Writer) {
	wm.l.Lock()
	defer wm.l.Unlock()

	var f Flusher
	if fi, ok := w.(Flusher); ok {
		f = fi
	}

	wm.ws = append(wm.ws, &writerCfg{
		w: w,
		f: f,
	})
}

func (wm *writerManager) Flush() error {
	var errs []string
	for _, w := range wm.ws {
		if w.f != nil {
			err := w.f.Flush()
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}
	if len(errs) <= 0 {
		return nil
	}
	return errors.New("error list: " + strings.Join(errs, "; "))
}

func (wm *writerManager) Write(r *Record) error {
	var errs []string
	for _, w := range wm.ws {
		if w.w != nil {
			err := w.w.Write(r)
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}
	if len(errs) <= 0 {
		return nil
	}
	return errors.New("error list: " + strings.Join(errs, "; "))
}

func (wm *writerManager) RenameFile(filebasename, pattern string) error {
	wm.l.Lock()
	defer wm.l.Unlock()

	var errs []string
	for _, w := range wm.ws {
		if r, ok := w.w.(Rotater); ok {
			err := r.SetPathPattern(filebasename, pattern)
			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
			err = r.Rotate()
			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
		}
	}
	if len(errs) <= 0 {
		return nil
	}
	return errors.New("error list: " + strings.Join(errs, "; "))
}

func (wm *writerManager) RotateByTime(t *time.Time) error {
	wm.l.Lock()
	defer wm.l.Unlock()

	var errs []string
	for _, w := range wm.ws {
		if r, ok := w.w.(Rotater); ok {
			if err := r.RotateByTime(t); err != nil {
				errs = append(errs, err.Error())
				continue
			}
		}
	}
	if len(errs) <= 0 {
		return nil
	}
	return errors.New("error list: " + strings.Join(errs, "; "))
}

func (wm *writerManager) Rotate() error {
	wm.l.Lock()
	defer wm.l.Unlock()

	var errs []string
	for _, w := range wm.ws {
		if r, ok := w.w.(Rotater); ok {
			if err := r.Rotate(); err != nil {
				errs = append(errs, err.Error())
				continue
			}
		}
	}
	if len(errs) <= 0 {
		return nil
	}
	return errors.New("error list: " + strings.Join(errs, "; "))
}

// recordWriter log写入器
type recordWriter struct {
	wm                   *writerManager
	tunnel               chan *Record
	c                    chan bool
	stopchan             chan struct{}
	opt                  *Options
	lastRecordTimeUnix60 int64
	lastRecordTime       time.Time
}

// Init 初始化log写入器
func (rec *recordWriter) Init(opt *Options) {
	rec.opt = opt
	rec.wm = &writerManager{}
	rec.tunnel = make(chan *Record, tunnelSizeDefault)
	rec.stopchan = make(chan struct{})
	rec.c = make(chan bool, 1)

	go rec.boostrapLogWriter()
}

// AddLogFile 增加一个文件输出器
func (rec *recordWriter) AddLogFile(filename string) error {
	//	fmt.Printf("log filename,%s \n", filename)
	filebasename := filename
	filename += ".%Y%M%D-%H"
	w := newFileWriter()
	w.RedirectError = rec.opt.RedirectError
	err := w.SetPathPattern(filebasename, filename)
	if err != nil {
		return err
	}
	err = w.Init()
	if err != nil {
		return err
	}
	rec.registerLogWriter(w)
	return nil
}

// ChangeLogFile 修改所有的文件输出器的目标文件
func (rec *recordWriter) ChangeLogFile(filename string) {
	filebasename := filename
	filename += ".%Y%M%D-%H"
	rec.wm.RenameFile(filebasename, filename)
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
		// close(rec.tunnel)
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

// write 写入一条日志记录，等待后续异步处理
func (rec *recordWriter) write(r *Record) {
	select {
	case <-rec.stopchan:
		return
	default:
	}

	// no async
	if !rec.opt.AsyncWrite {
		rec.doWriteRecord(r)
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
func (rec *recordWriter) boostrapLogWriter() {
	var (
		r  *Record
		ok bool
	)

	if r, ok = <-rec.tunnel; !ok {
		rec.c <- true
		return
	}

	if err := rec.wm.Write(r); err != nil {
		syslog.Println("boostrapLogWriter w.Write error", err)
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Millisecond * 100)
	//	rotateTimer := time.NewTimer(time.Second * 10)
	for {
		select {
		case r, ok = <-rec.tunnel:
			if !ok {
				rec.c <- true
				return
			}
			rec.doWriteRecord(r)
		case <-rec.stopchan:
			rec.c <- true
			return
		case <-flushTimer.C:
			// for _, w := range rec.writers {
			// 	if f, ok := w.(Flusher); ok {
			// 		if err := f.Flush(); err != nil {
			// 			syslog.Println("boostrapLogWriter f.Flush error", err)
			// 		}
			// 	}
			// }
			// flushTimer.Reset(time.Millisecond * 500)
		case <-rotateTimer.C:
			//	fmt.Printf("start rotate file,actions, 1111\n")
			rec.tryRotate()
			rotateTimer.Reset(time.Second * 60)
		}
	}
}

func (rec *recordWriter) doWriteRecord(r *Record) error {
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
		syslog.Println("boostrapLogWriter for w.Write error", err)
		return err
	}
	recordPool.Put(r)
	if err := rec.wm.Flush(); err != nil {
		syslog.Println("boostrapLogWriter f.Flush error", err)
		return err
	}
	return nil
}

// tryRotate 尝试将文件输出器转储
func (rec *recordWriter) tryRotate() {
	if err := rec.wm.Rotate(); err != nil {
		syslog.Println("r.Rotate error", err)
	}
}

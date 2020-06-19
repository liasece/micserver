package log

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// writerCfg config of log writer
type writerCfg struct {
	w Writer
	f Flusher
	r Rotater
}

func newWriterCfg(w Writer) *writerCfg {
	cfg := &writerCfg{
		w: w,
	}

	if v, ok := w.(Flusher); ok {
		cfg.f = v
	}

	if v, ok := w.(Rotater); ok {
		cfg.r = v
	}
	return cfg
}

// writerManager writer manager
type writerManager struct {
	ws []*writerCfg
	l  sync.Mutex
}

func (wm *writerManager) AddWriter(w Writer) {
	cfg := newWriterCfg(w)
	wm.l.Lock()
	defer wm.l.Unlock()
	wm.ws = append(wm.ws, cfg)
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

func (wm *writerManager) RotateByTime(t *time.Time) error {
	wm.l.Lock()
	defer wm.l.Unlock()

	var errs []string
	for _, w := range wm.ws {
		if w.r != nil {
			if err := w.r.RotateByTime(t); err != nil {
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
		if w.r != nil {
			if err := w.r.Rotate(); err != nil {
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

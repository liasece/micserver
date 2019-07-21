package module

import (
	"sync"
	"time"
)

type TimerCallback func(time.Duration)

type Timer struct {
	cb              TimerCallback
	timeDuration    time.Duration
	limitTimes      int64
	hasTriggerTimes int64
	hasKilled       bool
	killChan        chan struct{}
}

func (this *Timer) IsStop() bool {
	if this.hasKilled {
		return true
	}
	if this.limitTimes != 0 && this.hasTriggerTimes >= this.limitTimes {
		return true
	}
	return false
}

func (this *Timer) Select(exChan chan *Timer) {
	this.killChan = make(chan struct{})
	t := time.NewTimer(this.timeDuration)
	for !this.IsStop() {
		select {
		case <-t.C:
			this.hasTriggerTimes++
			if !this.IsStop() {
				t.Reset(this.timeDuration)
			}
			if exChan != nil {
				exChan <- this
			} else {
				this.cb(this.timeDuration)
			}
		case <-this.killChan:
			break
		}
	}
}

func (this *Timer) KillTimer() {
	this.hasKilled = true
	this.killChan <- struct{}{}
}

type Register struct {
	timerList            sync.Map
	timeTriggerChan      chan *Timer
	timeTriggerChanMutex sync.Mutex
	hasKilled            bool
}

// limitTimes 限制了这个定时任务重复执行的次数，如果为 0 ，那么将不限制
// 其执行的次数。
// 如果 engross 为 true，那么这个 timer 的执行将独占一个协程
// 如果一个定时操作很耗时，你应该将它作为一个单独的协程去处理，但是这样你可能
// 要考虑并行执行带来的问题
func (this *Register) RegTimer(duration time.Duration, limitTimes int64,
	engross bool, cb func(time.Duration)) {
	timer := &Timer{
		cb:           cb,
		timeDuration: duration,
		limitTimes:   limitTimes,
	}
	// 定时器初始化的检查
	if this.timeTriggerChan == nil {
		this.timeTriggerChanMutex.Lock()
		if this.timeTriggerChan == nil {
			this.timeTriggerChan = make(chan *Timer, 100)
			go this.goSelectTimer()
		}
		this.timeTriggerChanMutex.Unlock()
	}
	// 注册定时器
	resi, _ := this.timerList.LoadOrStore(duration, make([]*Timer, 0))
	if res, ok := resi.([]*Timer); ok {
		res = append(res, timer)
		if engross {
			go timer.Select(nil)
		} else {
			// 触发协程
			go timer.Select(this.timeTriggerChan)
		}
	}
}

func (this *Register) KillRegister() {
	this.timerList.Range(func(ki, listi interface{}) bool {
		if list, ok := listi.([]*Timer); ok {
			for _, t := range list {
				if t != nil {
					t.KillTimer()
				}
			}
		}
		return true
	})
	this.hasKilled = true
	this.timeTriggerChan <- nil
}

func (this *Register) goSelectTimer() {
	for !this.hasKilled {
		select {
		case t := <-this.timeTriggerChan:
			if t != nil {
				t.cb(t.timeDuration)
			}
		}
	}
}

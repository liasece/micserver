/*
定时器实现
*/
package timer

import (
	"sync"
	"time"
)

// 返回是否还将定时
type TimerCallback func(time.Duration) bool

// 定时器
type Timer struct {
	cb              TimerCallback
	timeDuration    time.Duration
	limitTimes      int64
	hasTriggerTimes int64
	hasKilled       bool
	killChan        chan struct{}
	mutex           sync.Mutex
}

// 判断定时器是否已经停止
func (this *Timer) IsStop() bool {
	if this.hasKilled {
		return true
	}
	if this.limitTimes != 0 && this.hasTriggerTimes >= this.limitTimes {
		return true
	}
	return false
}

// 开始监听定时器
func (this *Timer) Select(exChan chan *Timer) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.killChan = make(chan struct{}, 1)
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
				if !this.cb(this.timeDuration) {
					this.killTimer()
				}
			}
		case <-this.killChan:
			break
		}
	}
}

// 关闭定时器
func (this *Timer) KillTimer() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.killTimer()
}

// 关闭定时器
func (this *Timer) killTimer() {
	if this.hasKilled || this.killChan == nil {
		return
	}
	this.hasKilled = true
	this.killChan <- struct{}{}
}

// 定时器管理器
type TimerManager struct {
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
func (this *TimerManager) RegTimer(duration time.Duration, limitTimes int64,
	engross bool, cb func(time.Duration) bool) *Timer {
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
	return timer
}

// 关闭所有当前管理器中的定时器
func (this *TimerManager) KillRegister() {
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

func (this *TimerManager) goSelectTimer() {
	for !this.hasKilled {
		select {
		case t := <-this.timeTriggerChan:
			if t != nil {
				t.cb(t.timeDuration)
			}
		}
	}
}

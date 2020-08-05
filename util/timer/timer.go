/*
Package timer 定时器实现
*/
package timer

import (
	"sync"
	"time"
)

// FTimerCallback 返回是否还将定时
type FTimerCallback func(time.Duration) bool

// Timer 定时器
type Timer struct {
	cb              FTimerCallback
	timeDuration    time.Duration
	limitTimes      int64
	hasTriggerTimes int64
	hasKilled       bool
	killChan        chan struct{}
	mutex           sync.Mutex
}

// IsStop 判断定时器是否已经停止
func (tr *Timer) IsStop() bool {
	if tr.hasKilled {
		return true
	}
	if tr.limitTimes != 0 && tr.hasTriggerTimes >= tr.limitTimes {
		return true
	}
	return false
}

// Select 开始监听定时器
func (tr *Timer) Select(exChan chan *Timer) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	tr.killChan = make(chan struct{}, 1)
	t := time.NewTimer(tr.timeDuration)
	for !tr.IsStop() {
		select {
		case <-t.C:
			tr.hasTriggerTimes++
			if !tr.IsStop() {
				t.Reset(tr.timeDuration)
			}
			if exChan != nil {
				exChan <- tr
			} else {
				if !tr.cb(tr.timeDuration) {
					tr.killTimer()
				}
			}
		case <-tr.killChan:
			break
		}
	}
}

// KillTimer 关闭定时器
func (tr *Timer) KillTimer() {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	tr.killTimer()
}

// killTimer 关闭定时器
func (tr *Timer) killTimer() {
	if tr.hasKilled || tr.killChan == nil {
		return
	}
	tr.hasKilled = true
	tr.killChan <- struct{}{}
}

// Manager 定时器管理器
type Manager struct {
	timerList            sync.Map
	timeTriggerChan      chan *Timer
	timeTriggerChanMutex sync.Mutex
	hasKilled            bool
}

// RegTimer limitTimes 限制了这个定时任务重复执行的次数，如果为 0 ，那么将不限制
// 其执行的次数。
// 如果 engross 为 true，那么这个 timer 的执行将独占一个协程
// 如果一个定时操作很耗时，你应该将它作为一个单独的协程去处理，但是这样你可能
// 要考虑并行执行带来的问题
func (tr *Manager) RegTimer(duration time.Duration, limitTimes int64, engross bool, cb func(time.Duration) bool) *Timer {
	timer := &Timer{
		cb:           cb,
		timeDuration: duration,
		limitTimes:   limitTimes,
	}
	// 定时器初始化的检查
	if tr.timeTriggerChan == nil {
		tr.timeTriggerChanMutex.Lock()
		if tr.timeTriggerChan == nil {
			tr.timeTriggerChan = make(chan *Timer, 100)
			go tr.goSelectTimer()
		}
		tr.timeTriggerChanMutex.Unlock()
	}
	// 注册定时器
	resi, _ := tr.timerList.LoadOrStore(duration, make([]*Timer, 0))
	if res, ok := resi.([]*Timer); ok {
		res = append(res, timer)
		if engross {
			go timer.Select(nil)
		} else {
			// 触发协程
			go timer.Select(tr.timeTriggerChan)
		}
	}
	return timer
}

// KillRegister 关闭所有当前管理器中的定时器
func (tr *Manager) KillRegister() {
	tr.timerList.Range(func(ki, listi interface{}) bool {
		if list, ok := listi.([]*Timer); ok {
			for _, t := range list {
				if t != nil {
					t.KillTimer()
				}
			}
		}
		return true
	})
	tr.hasKilled = true
	tr.timeTriggerChan <- nil
}

func (tr *Manager) goSelectTimer() {
	for !tr.hasKilled {
		select {
		case t := <-tr.timeTriggerChan:
			if t != nil {
				t.cb(t.timeDuration)
			}
		}
	}
}

package xtime

import (
	"time"
)

// MockTimer 模拟time.Time
type MockTimer struct {
	C       <-chan time.Time // 参数同time.Time，非mock模式下对接到time.Time.C
	c       chan time.Time   // mock模式下内部通道
	next    time.Time        // mock模式下下一次tick时间
	mock    *mock            // mock对象
	fn      func()           // mock模式下AfterFunc 指定回调参数
	stopped bool             // mock模式下是否已关闭
	timer   *time.Timer      // 如果设定则走非Mock逻辑
}

// Stop turns off the ticker.
func (t *MockTimer) Stop() bool {
	if t.timer != nil {
		return t.timer.Stop()
	}
	t.mock.rw.Lock()
	registered := !t.stopped
	t.mock.removeClockTimer((*internalTimer)(t))
	t.stopped = true
	t.mock.rw.Unlock()
	return registered
}

// Reset 重置结束时间
func (t *MockTimer) Reset(d time.Duration) bool {
	if t.timer != nil {
		return t.timer.Reset(d)
	}
	var registered = false
	t.mock.doWithLock(func() {
		t.next = t.mock.nowWithoutLock().Add(d)
		registered = !t.stopped
		if t.stopped {
			t.mock.timers = append(t.mock.timers, (*internalTimer)(t))
		}
		t.stopped = false
	})
	return registered
}

type internalTimer MockTimer

func (t *internalTimer) Next() time.Time { return t.next }
func (t *internalTimer) Tick(now time.Time) {
	defer t.mock.cc.Gosched()
	t.mock.doWithLock(func() {
		t.stopped = true
		t.mock.removeClockTimer(t)
	})
	if t.fn != nil {
		defer func() {
			go t.fn()
		}()
	} else {
		t.c <- now
	}
}

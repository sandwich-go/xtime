package xtime

import (
	"time"
)

// MockTicker 模拟time.Ticker
type MockTicker struct {
	C        <-chan time.Time // 参数同time.Ticker，非mock模式下对接到time.Ticker.C
	c        chan time.Time   // mock下的内部通道
	next     time.Time        // mock下下一次tick的时间
	mock     *mock            // mock对象
	duration time.Duration    // mock下设定的tick duration
	ticker   *time.Ticker     // 如果设定则走非Mock逻辑
}

// Stop 关闭ticker.
func (t *MockTicker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		return
	}
	t.mock.doWithLock(func() {
		t.mock.removeClockTimer((*internalTicker)(t))
	})
}

// Reset 重置Duration
func (t *MockTicker) Reset(duration time.Duration) {
	if t.ticker != nil {
		t.ticker.Reset(duration)
		return
	}
	t.mock.doWithLock(func() {
		t.duration = duration
		t.next = t.mock.nowWithoutLock().Add(duration)
	})
}

type internalTicker MockTicker

func (t *internalTicker) Next() time.Time { return t.next }
func (t *internalTicker) Tick(now time.Time) {
	select {
	case t.c <- now:
	default:
	}
	t.next = now.Add(t.duration)
	t.mock.cc.Gosched()
}

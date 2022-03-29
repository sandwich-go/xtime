package xtime

import (
	"context"
	"time"
)

// TickerInterface 便于逻辑层传递Tikcer时使用，兼容系统接口与MockTicker
type TickerInterface interface {
	Stop()
	Reset(d time.Duration)
}

// TickerChan 获取底层ticker的chan
func TickerChan(t TickerInterface) <-chan time.Time {
	if st, ok := t.(*time.Ticker); ok {
		return st.C
	}
	return t.(*MockTicker).C
}

// TimerInterface 便于逻辑层传递Tikcer时使用，兼容系统接口与MockTimer
type TimerInterface interface {
	Stop() bool
	Reset(d time.Duration) bool
}

// TimerChan 获取底层timer的chan
func TimerChan(t TimerInterface) <-chan time.Time {
	if st, ok := t.(*time.Timer); ok {
		return st.C
	}
	return t.(*MockTimer).C
}

// Clock 兼容系统timer方法
type Clock interface {
	After(d time.Duration) <-chan time.Time
	Now() time.Time
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	NewTicker(d time.Duration) *MockTicker
	Timer(d time.Duration) *MockTimer
	AfterFunc(d time.Duration, f func()) *MockTimer
	WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc)
	WithTimeout(parent context.Context, t time.Duration) (context.Context, context.CancelFunc)
}

type Cop interface {
	Freeze(t time.Time)
	Travel(t time.Time)
	Scale(scale float64)
	Return()
	Mocked() bool
	ApplyOption(...Option) []Option
}

type Mock interface {
	Clock
	Cop
}

func NewSystemClock() Clock       { return &clock{} }
func NewMock(opts ...Option) Mock { return newMock(NewOptions(opts...)) }

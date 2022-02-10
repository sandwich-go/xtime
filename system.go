package xtime

import (
	"context"
	"time"
)

type clock struct{}

func (c *clock) After(d time.Duration) <-chan time.Time { return time.After(d) }

func (c *clock) AfterFunc(d time.Duration, f func()) *MockTimer {
	return &MockTimer{timer: time.AfterFunc(d, f)}
}

func (c *clock) Now() time.Time { return time.Now() }

func (c *clock) Since(t time.Time) time.Duration { return time.Since(t) }

func (c *clock) Until(t time.Time) time.Duration { return time.Until(t) }

func (c *clock) Sleep(d time.Duration) { time.Sleep(d) }

func (c *clock) Tick(d time.Duration) <-chan time.Time { return time.Tick(d) }

func (c *clock) NewTicker(d time.Duration) *MockTicker {
	t := time.NewTicker(d)
	return &MockTicker{C: t.C, ticker: t}
}

func (c *clock) Timer(d time.Duration) *MockTimer {
	t := time.NewTimer(d)
	return &MockTimer{C: t.C, timer: t}
}

func (c *clock) WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, d)
}

func (c *clock) WithTimeout(parent context.Context, t time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, t)
}

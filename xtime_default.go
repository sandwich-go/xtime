package xtime

import (
	"context"
	"time"
)

var defaultMock = newMockNotStart(newDefaultOptions())

func getDefaultMock() Mock {
	defaultMock.start()
	return defaultMock
}

func ApplyOption(opt ...Option)                      { getDefaultMock().ApplyOption(opt...) }
func Now() time.Time                                 { return getDefaultMock().Now() }
func Since(t time.Time) time.Duration                { return getDefaultMock().Since(t) }
func Until(t time.Time) time.Duration                { return getDefaultMock().Until(t) }
func Sleep(d time.Duration)                          { getDefaultMock().Sleep(d) }
func Tick(d time.Duration) <-chan time.Time          { return getDefaultMock().Tick(d) }
func After(d time.Duration) <-chan time.Time         { return getDefaultMock().After(d) }
func AfterFunc(d time.Duration, f func()) *MockTimer { return getDefaultMock().AfterFunc(d, f) }
func Timer(d time.Duration) *MockTimer               { return getDefaultMock().Timer(d) }
func NewTicker(d time.Duration) *MockTicker          { return getDefaultMock().NewTicker(d) }
func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return getDefaultMock().WithDeadline(parent, d)
}
func WithTimeout(parent context.Context, t time.Duration) (context.Context, context.CancelFunc) {
	return getDefaultMock().WithTimeout(parent, t)
}

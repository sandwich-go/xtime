package xtime

import (
	"context"
	"time"
)

var defaultMock = &mock{}

func ApplyOption(opt ...Option)                      { defaultMock.ApplyOption(opt...) }
func Now() time.Time                                 { return defaultMock.Now() }
func Since(t time.Time) time.Duration                { return defaultMock.Since(t) }
func Until(t time.Time) time.Duration                { return defaultMock.Until(t) }
func Sleep(d time.Duration)                          { defaultMock.Sleep(d) }
func Tick(d time.Duration) <-chan time.Time          { return defaultMock.Tick(d) }
func After(d time.Duration) <-chan time.Time         { return defaultMock.After(d) }
func AfterFunc(d time.Duration, f func()) *MockTimer { return defaultMock.AfterFunc(d, f) }
func Timer(d time.Duration) *MockTimer               { return defaultMock.Timer(d) }
func NewTicker(d time.Duration) *MockTicker          { return defaultMock.NewTicker(d) }
func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return defaultMock.WithDeadline(parent, d)
}
func WithTimeout(parent context.Context, t time.Duration) (context.Context, context.CancelFunc) {
	return defaultMock.WithTimeout(parent, t)
}

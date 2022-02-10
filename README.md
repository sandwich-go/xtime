# xtime
```golang
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
	ApplyOption(...Option) []Option
}

type Mock interface {
	Clock
	Cop
}
```

package xtime

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

var _ Mock = &mock{}

func newMock(opts *Options) Mock {
	m := &mock{cc: opts}
	m.continueTick = make(chan struct{}, 1)
	m.tickStopChan = make(chan struct{})
	m.freshTicker()
	return m
}

type mock struct {
	cc           *Options
	rw           sync.RWMutex
	scale        float64
	frozen       bool
	traveled     bool
	freezeTime   time.Time
	travelTime   time.Time
	timers       clockTimers // tickers & timers
	continueTick chan struct{}
	tickStopChan chan struct{}
}

func (m *mock) ApplyOption(opt ...Option) []Option {
	old := m.cc.ApplyOption(opt...)
	m.freshTicker()
	return old
}

func (m *mock) freshTicker() {
	close(m.tickStopChan)
	m.tickStopChan = make(chan struct{})
	go m.tick(m.tickStopChan, m.cc.TickIntervalUnderMock)
}

func (m *mock) debugLogN(format string, a ...interface{}) {
	if !m.cc.Debug {
		return
	}
	format += "\n"
	fmt.Fprintf(m.cc.DebugWriter, format, a...)
}

func (m *mock) tick(stop chan struct{}, duration time.Duration) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			m.tickCallback(stop)
			select {
			case <-stop:
				return
			default:
			}
		case <-stop:
			ticker.Stop()
			return
		}
	}
}
func (m *mock) tickCallback(stop chan struct{}) {
	frozen := false
	m.rw.RLock()
	frozen = m.frozen
	m.rw.RUnlock()
	if frozen {
		// 检测到frozen主动运行一次timer
		for {
			if !m.runNextTimer(m.Now()) {
				break
			}
			m.cc.Gosched()
		}
		m.debugLogN("frozen waiting for continueTick chan")
		select {
		case <-m.continueTick:
			m.tickCallback(stop)
		case <-stop:
			return
		}
	}
	for {
		if !m.runNextTimer(m.Now()) {
			break
		}
		m.cc.Gosched()
	}
}

// Scale 缩放时间流逝的比例,会自动Travel到当前节点，从当前节点开始scale
func (m *mock) Scale(scale float64) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.scale = scale
	if !m.traveled {
		m.travelUnlock(m.cc.NowProvider())
	}
}

func (m *mock) Return() {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.frozen = false
	m.traveled = false
	m.scale = 1
	m.notifyContueTick()
}

// Freeze 静止在指定时间
func (m *mock) Freeze(t time.Time) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.freezeTime = t
	m.frozen = true
	// 如果首次进入frozen,tickCallback中检测到frozen会主动执行一次timer，堵塞在continueTick
	// 如果在frozen状态下再次Freeze,则通知一次contueTick驱动执行timer再次堵塞
	m.notifyContueTick()
}

func (m *mock) notifyContueTick() {
	// travel 模式下依然需要依赖tick驱动ticker与timer
	select {
	case <-m.continueTick:
	default:
	}
	m.continueTick <- struct{}{}
}
func (m *mock) travelUnlock(t time.Time) {
	m.freezeTime = t
	// 获取Travel时所处的时间点
	m.travelTime = m.cc.NowProvider()
	m.traveled = true
	// 主动解冻，时间开始流逝
	m.frozen = false
	// 时间跳转之后，主动运转队列中的timer，同时通知协程开启tick
	m.notifyContueTick()
}

// Travel 跳转到指定时间后开始时间流逝,并自动解冻
func (m *mock) Travel(t time.Time) {
	m.rw.Lock()
	m.travelUnlock(t)
	m.rw.Unlock()
	m.cc.Gosched()
}

// runNextTimer timer调用
func (m *mock) runNextTimer(max time.Time) bool {
	m.rw.Lock()
	// Sort timers by time.
	sort.Sort(m.timers)
	// If we have no more timers then exit.
	if len(m.timers) == 0 {
		m.rw.Unlock()
		return false
	}
	// Retrieve next timer. Exit if next tick is after new time.
	t := m.timers[0]
	next := t.Next()
	if next.After(max) {
		m.rw.Unlock()
		return false
	}
	m.debugLogN("runNextTimer next:%s max:%s", next, max)
	m.rw.Unlock()
	now := next
	if m.cc.TickAtMockNow {
		now = m.nowWithoutLock()
	}
	t.Tick(now) // 如果在Travel时有一个执行频繁的ticker，可能会导致ticker的执行一直占用tick协程导致其他的timer无法被及时执行
	return true
}

// After waits for the duration to elapse and then sends the current time on the returned channel.
func (m *mock) After(d time.Duration) <-chan time.Time {
	return m.Timer(d).C
}

// AfterFunc waits for the duration to elapse and then executes a function.
// A Timer is returned that can be stopped.
func (m *mock) AfterFunc(d time.Duration, f func()) *MockTimer {
	t := m.Timer(d)
	t.C = nil
	t.fn = f
	return t
}

// Now returns the current wall time on the mock clock.
func (m *mock) Now() time.Time {
	if m.frozen || m.traveled {
		m.rw.RLock()
		defer m.rw.RUnlock()
	}
	return m.nowWithoutLock()
}

func (m *mock) nowWithoutLock() time.Time {
	if m.frozen {
		return m.freezeTime
	}
	if m.traveled {
		return m.freezeTime.Add(time.Duration(float64(time.Since(m.travelTime)) * m.scale))
	}
	return m.cc.NowProvider()
}

// Since returns time since `t` using the mock clock's wall time.
func (m *mock) Since(t time.Time) time.Duration {
	return m.Now().Sub(t)
}

// Until returns time until `t` using the mock clock's wall time.
func (m *mock) Until(t time.Time) time.Duration {
	return t.Sub(m.Now())
}

// Sleep pauses the goroutine for the given duration on the mock clock.
// The clock must be moved forward in a separate goroutine.
func (m *mock) Sleep(d time.Duration) {
	<-m.After(d)
}

// Tick is a convenience function for Ticker().
// It will return a ticker channel that cannot be stopped.
func (m *mock) Tick(d time.Duration) <-chan time.Time {
	return m.NewTicker(d).C
}

// Ticker creates a new instance of Ticker.
func (m *mock) NewTicker(duration time.Duration) *MockTicker {
	m.rw.Lock()
	defer m.rw.Unlock()
	ch := make(chan time.Time, 1)
	t := &MockTicker{
		C:        ch,
		c:        ch,
		mock:     m,
		duration: duration,
		next:     m.nowWithoutLock().Add(duration),
	}
	m.timers = append(m.timers, (*internalTicker)(t))
	return t
}

// Timer creates a new instance of Timer.
func (m *mock) Timer(d time.Duration) *MockTimer {
	m.rw.Lock()
	defer m.rw.Unlock()
	ch := make(chan time.Time, 1)
	t := &MockTimer{
		C:       ch,
		c:       ch,
		mock:    m,
		next:    m.nowWithoutLock().Add(d),
		stopped: false,
	}
	m.timers = append(m.timers, (*internalTimer)(t))
	return t
}

func (m *mock) doWithLock(f func()) {
	m.rw.Lock()
	defer m.rw.Unlock()
	f()
}

func (m *mock) removeClockTimer(t clockTimer) {
	for i, timer := range m.timers {
		if timer == t {
			copy(m.timers[i:], m.timers[i+1:])
			m.timers[len(m.timers)-1] = nil
			m.timers = m.timers[:len(m.timers)-1]
			break
		}
	}
	sort.Sort(m.timers)
}

// clockTimer represents an object with an associated start time.
type clockTimer interface {
	Next() time.Time
	Tick(time.Time)
}

// clockTimers represents a list of sortable timers.
type clockTimers []clockTimer

func (a clockTimers) Len() int           { return len(a) }
func (a clockTimers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a clockTimers) Less(i, j int) bool { return a[i].Next().Before(a[j].Next()) }

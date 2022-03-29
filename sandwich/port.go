package sandwich

import (
	"time"

	"bitbucket.org/funplus/sandwich/base/stime"
	"github.com/sandwich-go/xtime"
)

func PatchSandwichTime(mocked xtime.Mock) {
	stime.Now = func() time.Time {
		return mocked.Now()
	}
	stime.Mocked = func() bool { return true }

	stime.Freeze = func(t time.Time) {
		mocked.Freeze(t)
	}
	stime.Scale = func(scale float64) {
		mocked.Scale(scale)
	}
	stime.Travel = func(t time.Time) {
		mocked.Travel(t)
	}

	stime.Since = func(t time.Time) time.Duration {
		return mocked.Since(t)
	}

	stime.Sleep = func(d time.Duration) {
		mocked.Sleep(d)
	}

	stime.After = func(d time.Duration) <-chan time.Time {
		return mocked.After(d)
	}

	stime.Tick = func(d time.Duration) <-chan time.Time {
		return mocked.Tick(d)
	}

	stime.Return = func() {
		mocked.Return()
	}
}

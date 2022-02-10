package xtime

// Patch 替换系统time package中的实现，但由于NewTicker，Timer，AfterFunc返回的类型与time package中一致无法实现Patch
// func Patch(m Clock) {
// 	monkey.Patch(time.Now, func() time.Time {
// 		return m.Now()
// 	})

// 	monkey.Patch(time.After, func(d time.Duration) <-chan time.Time {
// 		return m.After(d)
// 	})

// 	monkey.Patch(time.Since, func(t time.Time) time.Duration {
// 		return m.Since(t)
// 	})

// 	monkey.Patch(time.Until, func(t time.Time) time.Duration {
// 		return m.Until(t)
// 	})

// 	monkey.Patch(time.Sleep, func(d time.Duration) {
// 		m.Sleep(d)
// 	})

// 	monkey.Patch(time.Tick, func(d time.Duration) <-chan time.Time {
// 		return m.Tick(d)
// 	})

// 	monkey.Patch(context.WithDeadline, func(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
// 		return m.WithDeadline(parent, d)
// 	})
// 	monkey.Patch(context.WithTimeout, func(parent context.Context, t time.Duration) (context.Context, context.CancelFunc) {
// 		return m.WithTimeout(parent, t)
// 	})
// }

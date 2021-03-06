package xtime

import (
	"io"
	"os"
	"time"
)

type NowProvider = func() time.Time
type SleepProviderUnderFrozen = func(m Mock, d time.Duration)

var timeNow = time.Now
var timeSleep = time.Sleep

//go:generate optionGen  --option_return_previous=true
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@Gosched(comment="Gosched让出CPU防止忙占")
		"Gosched": func() { timeSleep(1 * time.Millisecond) },
		// annotation@TickIntervalUnderMock(comment="真实的tick时间间隔，用于驱动mock模式下的tiker、timer")
		"TickIntervalUnderMock": time.Duration(time.Millisecond),
		// annotation@TickAtMockNow(comment="timer,ticker在tick的时间为mock的当前时间，而不是Next时间，如果为Next时间,在时间跳转后会导致循环执行同一个ticker,timer")
		"TickAtMockNow": false,
		// annotation@NowProvider(comment="系统时间")
		"NowProvider": NowProvider(func() time.Time {
			return timeNow()
		}),
		// annotation@SleepProviderUnderFrozen(comment="Frozen模式下sleep将再次Freeze一次,将时间向前推进")
		"SleepProviderUnderFrozen": SleepProviderUnderFrozen(func(m Mock, d time.Duration) {
			m.Freeze(m.Now().Add(d))
		}),
		// annotation@Debug(comment="debug模式下以会向DebugWriter写日志")
		"Debug": false,
		// annotation@DebugWriter(comment="调试日志输出")
		"DebugWriter": io.Writer(os.Stdout),
	}
}

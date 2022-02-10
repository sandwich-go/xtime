package xtime_test

import (
	"net/http"
	"os"
	"sync/atomic"
	"testing"
	"time"

	_ "net/http/pprof"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/sandwich-go/xtime"
)

func TestMain(m *testing.M) {
	go func() { http.ListenAndServe(":8888", nil) }()
	os.Exit(m.Run())
}

var m = xtime.NewMock(xtime.WithDebug(true))
var stubTime = time.Unix(1522549800, 0) //Human time (GMT): Sunday, April 1, 2018 2:30:00 AM

func subTestReturn(m xtime.Mock) {
	Convey("Return之后应该继续之前的时间", func() {
		m.Return()
		now := time.Now().Unix()
		So(m.Now().Unix(), ShouldBeBetween, now-1, now+1)
	})
}
func TestMock(t *testing.T) {
	m.ApplyOption(xtime.WithTickIntervalUnderMock(time.Millisecond * 5))
	Convey("Freeze", t, func() {
		m.Freeze(stubTime)
		time.Sleep(time.Second)
		Convey("时间不再流逝", func() {
			So(m.Now(), ShouldEqual, stubTime)
		})
		Convey("AfterFunc应该停止运行不会被触发", func() {
			timerGotTs := int64(0)
			timer := m.AfterFunc(time.Millisecond*10, func() {
				timerGotTs = time.Now().Unix()
			})
			defer timer.Stop()
			time.Sleep(time.Second)
			So(timerGotTs, ShouldEqual, 0)
			Convey("恢复之后, timer被触发", func() {
				subTestReturn(m)
				time.Sleep(time.Second)
				So(m.Now().Unix()-timerGotTs, ShouldBeBetween, 0, 2)
			})
		})
		Convey("After应该停止运行不会被触发", func() {
			afterChan := m.After(time.Second)
			timerGotTs := int64(0)
			go func() {
				<-afterChan
				timerGotTs = time.Now().Unix()
			}()
			time.Sleep(time.Second)
			So(timerGotTs, ShouldEqual, 0)
			Convey("Freeze推进，timer应该被触发", func() {
				m.Freeze(stubTime.Add(time.Second * 5))
				// 给ticker一个执行的机会
				time.Sleep(time.Second)
				So(time.Now().Unix()-timerGotTs, ShouldBeBetween, 0, 2)
			})
			Convey("恢复之后, timer被触发", func() {
				subTestReturn(m)
				time.Sleep(time.Second)
				So(m.Now().Unix()-timerGotTs, ShouldBeBetween, 0, 2)
			})
		})
		Convey("Tick应该停止运行不会被触发", func() {
			ticker := m.NewTicker(time.Millisecond * 10)
			var tickerCount AtomicInt32
			go func() {
				for {
					<-ticker.C
					tickerCount.Add(1)
					if tickerCount.Get() == 10 {
						break
					}
				}
			}()
			time.Sleep(time.Second)
			So(tickerCount.Get(), ShouldEqual, 0)
			Convey("Freeze推进5s，当Tick选取当前时间作为Tick的时间点时，只会运行一次，因为时间停止了", func() {
				previous := m.ApplyOption(xtime.WithTickAtMockNow(true))
				m.Freeze(stubTime.Add(time.Second * 5))
				time.Sleep(time.Second * 1)
				So(tickerCount.Get(), ShouldEqual, 1)
				m.ApplyOption(previous...)
			})
			Convey("恢复之后, ticker运行，当Tick选取的时间为正常的Next时间，则会Tick多次", func() {
				previous := m.ApplyOption(xtime.WithTickAtMockNow(false))
				m.Freeze(stubTime.Add(time.Second * 5))
				subTestReturn(m)
				time.Sleep(time.Second * 1)
				So(tickerCount.Get(), ShouldEqual, 10)
				m.ApplyOption(previous...)
			})
			ticker.Stop()
		})
	})
	Convey("Travel: 时间跳转后继续流逝", t, func() {
		m.Travel(stubTime)
		time.Sleep(time.Second)
		Convey("时间跳转后继续流逝", func() {
			sub := m.Now().Sub(stubTime)
			So(sub.Seconds(), ShouldBeBetween, 0, 2)
		})
		Convey("AfterFunc会被立刻触发", func() {
			timerGotTs := int64(0)
			timer := m.AfterFunc(time.Hour, func() {
				timerGotTs = time.Now().Unix()
			})
			defer timer.Stop()
			m.Travel(stubTime.Add(time.Second * 5).Add(time.Hour))
			time.Sleep(time.Second * 2)
			// 间隙，让ticker有机会执行
			So(timerGotTs, ShouldNotEqual, 0)
		})
		subTestReturn(m)
	})
}

type AtomicInt32 int32

func (i *AtomicInt32) Add(n int32) int32 {
	return atomic.AddInt32((*int32)(i), n)
}

func (i *AtomicInt32) Set(n int32) {
	atomic.StoreInt32((*int32)(i), n)
}

func (i *AtomicInt32) Get() int32 {
	return atomic.LoadInt32((*int32)(i))
}

func (i *AtomicInt32) CompareAndSwap(oldval, newval int32) (swapped bool) {
	return atomic.CompareAndSwapInt32((*int32)(i), oldval, newval)
}

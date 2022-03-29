package sandwich

import (
	"net/http"
	"os"
	"testing"
	"time"

	_ "net/http/pprof"

	. "github.com/smartystreets/goconvey/convey"

	"bitbucket.org/funplus/sandwich/base/stime"
	"github.com/sandwich-go/xtime"
)

func TestMain(m *testing.M) {
	go func() { http.ListenAndServe(":8888", nil) }()
	os.Exit(m.Run())
}

var stubTime = time.Unix(1522549800, 0) //Human time (GMT): Sunday, April 1, 2018 2:30:00 AM
func TestMock(t *testing.T) {
	Convey("Freeze", t, func() {
		mock := xtime.NewMock(
			xtime.WithDebug(true),
			xtime.WithTickIntervalUnderMock(time.Millisecond*5),
		)
		PatchSandwichTime(mock)
		stime.Freeze(stubTime)
		Convey("时间不再流逝", func() {
			So(stime.Now(), ShouldEqual, stubTime)
		})
	})
}

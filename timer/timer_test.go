package timer_test

import (
	. "deepin-file-manager/time"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	Convey("just stop", t, func() {
		timer := NewTimer()
		timer.Stop()
		So(timer.Elapsed(), ShouldEqual, 0)
	})

	Convey("get elapse without stop", t, func() {
		timer := NewTimer()
		timer.Start()

		time.Sleep(time.Second)
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second-time.Millisecond, time.Second+time.Millisecond)

		time.Sleep(time.Second)
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second*2-time.Millisecond, time.Second*2+time.Millisecond*2)
	})

	Convey("stop and elapse", t, func() {
		timer := NewTimer()
		timer.Start()

		time.Sleep(time.Second)
		timer.Stop()
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second-time.Millisecond, time.Second+time.Millisecond)

		time.Sleep(time.Second)
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second-time.Millisecond, time.Second+time.Millisecond)
	})

	Convey("stop and continue", t, func() {
		timer := NewTimer()
		timer.Start()

		time.Sleep(time.Second)
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second-time.Millisecond, time.Second+time.Millisecond)

		timer.Stop()
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second-time.Millisecond, time.Second+time.Millisecond)

		time.Sleep(time.Second)
		timer.Continue()
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second-time.Millisecond, time.Second+time.Millisecond)

		time.Sleep(time.Second)
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second*2-time.Millisecond, time.Second*2+time.Millisecond)
	})

	Convey("reset", t, func() {
		timer := NewTimer()
		timer.Start()

		time.Sleep(time.Second)
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second-time.Millisecond, time.Second+time.Millisecond)

		timer.Reset()
		time.Sleep(time.Second)
		So(timer.Elapsed(), ShouldBeBetweenOrEqual, time.Second-time.Millisecond, time.Second+time.Millisecond)
	})
}

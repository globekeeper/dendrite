package ratelimit

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestRtFailedLogin(t *testing.T) {
	is := is.New(t)
	dut := NewRtFailedLogin(&RtFailedLoginConfig{
		Enabled:  true,
		Limit:    3,
		Interval: 10 * time.Millisecond,
	})
	var (
		can        bool
		remaining  time.Duration
		remainingB time.Duration
	)
	for i := 0; i < 3; i++ {
		can, remaining = dut.CanAct("foo")
		is.True(can)
		is.Equal(remaining, time.Duration(0))
		dut.Act("foo")
	}
	can, remaining = dut.CanAct("foo")
	is.True(!can)
	is.True(remaining > time.Millisecond*9)
	can, remainingB = dut.CanAct("bar")
	is.True(can)
	is.Equal(remainingB, time.Duration(0))
	dut.Act("bar")
	dut.Act("bar")
	time.Sleep(remaining + time.Millisecond)
	can, remaining = dut.CanAct("foo")
	is.True(can)
	is.Equal(remaining, time.Duration(0))
}

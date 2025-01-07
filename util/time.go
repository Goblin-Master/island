package util

import (
	"tgwp/log/zlog"
	"time"
)

// RecordTime a tool to record time
// e.g [defer util.RecordTime(time.Now())()]
func RecordTime(start time.Time) func() {
	return func() {
		end := time.Now()
		zlog.Debugf("use time:%d", end.Unix()-start.Unix())
	}
}

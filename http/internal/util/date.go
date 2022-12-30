package util

import (
	"sync"
	"time"

	"go.uber.org/atomic"
	"go.x2ox.com/sorbifolia/coarsetime"
)

var (
	serverDate     atomic.Value
	serverDateOnce sync.Once
)

func init() {
	serverDateOnce.Do(func() {
		refreshServerDate()
		go func() {
			for {
				time.Sleep(time.Second)
				refreshServerDate()
			}
		}()
	})
}

func refreshServerDate() {
	b := AppendHTTPDate(nil, coarsetime.Now())
	serverDate.Store(b)
}

// AppendHTTPDate appends HTTP-compliant (RFC1123) representation of date
// to dst and returns the extended dst.
func AppendHTTPDate(dst []byte, date time.Time) []byte {
	dst = date.In(time.UTC).AppendFormat(dst, time.RFC1123)
	copy(dst[len(dst)-3:], "GMT")
	return dst
}

func GetDate() []byte {
	return serverDate.Load().([]byte)
}

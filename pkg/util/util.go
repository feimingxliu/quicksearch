package util

import "time"

func ExecTime(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

package main

import (
	"time"
)

// preciseSleep sleeps for the specified number of milliseconds with improved accuracy
func preciseSleep(milsec int64) error {
	start := time.Now()
	for time.Since(start).Milliseconds() < milsec {
		time.Sleep(time.Nanosecond)
	}
	return nil
}

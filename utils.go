package main

import (
	"time"
)

func preciseSleep(milsec int64) error {
	start := time.Now()
	for time.Since(start).Milliseconds() < milsec {
		time.Sleep(time.Nanosecond)
	}
	return nil
}

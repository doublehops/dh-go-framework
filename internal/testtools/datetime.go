package testtools

import "time"

func GetTolerance(seconds int) (time.Time, time.Duration) {
	duration := time.Duration(seconds)
	timeNow := time.Now()
	expectedTime := timeNow.Add(-duration)
	tolerance := 10 * time.Second

	return expectedTime, tolerance
}

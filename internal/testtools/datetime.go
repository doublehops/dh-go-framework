package testtools

import "time"

func GetTolerance(seconds int) (time.Time, time.Duration) {
	duration := time.Duration(seconds)
	timeNow := time.Now()
	expectedTime := timeNow.Add(-duration * time.Second) // Adjust this as needed for your test
	tolerance := 10 * time.Second

	return expectedTime, tolerance
}

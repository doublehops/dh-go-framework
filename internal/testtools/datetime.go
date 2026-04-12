// Package testtools provides shared helper utilities for integration and unit tests.
package testtools

import "time"

// GetTolerance returns an expected time and tolerance duration for time-based assertions in tests.
func GetTolerance(seconds int) (time.Time, time.Duration) {
	duration := time.Duration(seconds)
	timeNow := time.Now()
	expectedTime := timeNow.Add(-duration)
	tolerance := 10 * time.Second

	return expectedTime, tolerance
}

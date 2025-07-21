package utils

import "time"

// GetTimeRoundedTo rounds the current time to the nearest given interval (in nanoseconds)
func GetTimeRoundedTo(interval time.Duration) time.Time {
	now := time.Now()

	// Calculate the nearest multiple of the interval
	unixTime := now.UnixNano()                                        // Convert time to nanoseconds
	roundedTimeNano := (unixTime / int64(interval)) * int64(interval) // Round down

	// If past halfway, round up
	if unixTime%int64(interval) >= int64(interval)/2 {
		roundedTimeNano += int64(interval)
	}

	// Convert nanoseconds back to time.Time
	roundedTime := time.Unix(0, roundedTimeNano)

	return roundedTime
}

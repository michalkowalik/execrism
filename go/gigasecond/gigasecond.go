// Package gigasecond provides solution to the gigasecond problem
package gigasecond

import "time"

const testVersion = 4

// AddGigasecond returns time updated by 1e9 seconds
func AddGigasecond(t time.Time) time.Time {
	return t.Add(time.Duration(1e9) * time.Second)
}

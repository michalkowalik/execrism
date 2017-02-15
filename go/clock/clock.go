// Package clock contains solution to ex. 3 - Clock
package clock

import (
	"fmt"
)

// The value of testVersion here must match `targetTestVersion` in the file
// clock_test.go.
const testVersion = 4

// Clock type represents clock
type Clock struct {
	Hour   int
	Minute int
}

// New returns new Clock instance
func New(hour, minute int) Clock {
	h := (hour + (minute / 60)) % 24
	m := minute % 60
	if m < 0 {
		h--
		m = 60 + m
	}
	if h < 0 {
		h = 24 + h
	}

	return Clock{Hour: h, Minute: m}
}

// String returns string representation of clock state
func (c Clock) String() string {
	return fmt.Sprintf("%02d:%02d", c.Hour, c.Minute)
}

// Add set time by adding minutes to the existing clock
func (c Clock) Add(minutes int) Clock {
	return New(c.Hour, c.Minute+minutes)
}

package util

import "time"

// Clock - wrapper around built in time
type Clock struct {
}

// CreateNewClock - Create a new instance of Clock
func CreateNewClock() *Clock {
	return &Clock{}
}

// Now - uses built in time.Now()
func (c *Clock) Now() time.Time {
	return time.Now()
}

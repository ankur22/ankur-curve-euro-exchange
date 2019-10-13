package util

import "time"

// Clock - wrapper around built in time
type Clock struct {
}

// Now - uses built in time.Now()
func (c *Clock) Now() time.Time {
	return time.Now()
}

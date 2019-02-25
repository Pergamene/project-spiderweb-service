package clock

import "time"

// Clock is used to mock out time.Now() for testing statically defined times
// When providing a Clock, either provide RealClock or MockClock
// (or another struct that interfaces Clock)
type Clock interface {
	Now() time.Time
}

// RealClock acts just like time.Now(). Should be used as Clock for normal runs.
type RealClock struct{}

// Now is the same as time.Now()
func (t RealClock) Now() time.Time {
	return time.Now()
}

// MockClock replaces time.Now() with the MockedTime.
type MockClock struct {
	MockedTime *time.Time
}

// Now replaces time.Now()
func (t MockClock) Now() time.Time {
	if t.MockedTime != nil {
		return *t.MockedTime
	}
	return time.Now()
}

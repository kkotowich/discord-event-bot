package types

import "time"

// Date struct that will be used for Date part only
type Date struct {
	Date time.Time
}

// MilitaryTime struct that willbe used for Time part only in Military format
type MilitaryTime struct {
	Time time.Time
}

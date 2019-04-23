package data

import "time"

func (s Act) Date() time.Time {
	return time.Date( s.Year, s.Month, s.Day,0,0,0,0, time.Local)
}


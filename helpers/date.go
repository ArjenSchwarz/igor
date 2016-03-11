package helpers

import (
	"time"
)

// RoughDay turns a unix date into a rough string approximation
// * Anything within an hour before or after is Now
// * Anything on the same day outside of that is Today
// * Anything the next day is Tomorrow
// * Anything else returns its weekday
func RoughDay(dateInt int64) string {
	date := time.Unix(dateInt, 0)
	beforeMargin, _ := time.ParseDuration("-1h")
	afterMargin, _ := time.ParseDuration("1h")
	start := time.Now().UTC().Add(beforeMargin)
	end := time.Now().UTC().Add(afterMargin)
	if date.After(start) && date.Before(end) {
		return "Now"
	}
	pYear, pMonth, pDay := date.Date()
	cYear, cMonth, cDay := time.Now().UTC().Date()
	if pYear == cYear && pMonth == cMonth && pDay == cDay {
		return "Today"
	}
	dayMargin, _ := time.ParseDuration("24h")
	cYear, cMonth, cDay = time.Now().UTC().Add(dayMargin).Date()
	if pYear == cYear && pMonth == cMonth && pDay == cDay {
		return "Tomorrow"
	}
	return date.Weekday().String()
}

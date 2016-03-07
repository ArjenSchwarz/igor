package helpers

import (
	"time"
)

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

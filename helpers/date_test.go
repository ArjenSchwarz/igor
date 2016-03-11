package helpers_test

import (
	"fmt"
	"github.com/ArjenSchwarz/igor/helpers"
	"testing"
	"time"
)

func TestRoughDay(t *testing.T) {
	zone, _ := time.Now().Zone()
	now := time.Now().UTC().Unix()
	year, month, day := time.Now().UTC().Date()
	layout := "2006-January-2 15:04 (MST)"
	startOfDay, _ := time.Parse(layout, fmt.Sprintf("%v-%v-%v 00:01 (%s)", year, month, day, zone))
	endOfDay, _ := time.Parse(layout, fmt.Sprintf("%v-%v-%v 23:59 (%s)", year, month, day, zone))
	plus59Min := now + (59 * 60)
	minus59Min := now - (59 * 60)
	tomorrow := time.Now().UTC().AddDate(0, 0, 1).Unix()
	dayAfter := time.Now().UTC().AddDate(0, 0, 2)
	var roughDayTests = []struct {
		input    int64
		expected string
	}{
		{now, "Now"},
		{plus59Min, "Now"},
		{minus59Min, "Now"},
		{tomorrow, "Tomorrow"},
		{dayAfter.Unix(), dayAfter.Weekday().String()},
		{startOfDay.Unix(), "Today"},
		{endOfDay.Unix(), "Today"},
	}

	for _, tt := range roughDayTests {
		actual := helpers.RoughDay(tt.input)
		if actual != tt.expected {
			t.Errorf("RoughDay(%v): expected %v, actual %v", tt.input, tt.expected, actual)
		}
	}
}

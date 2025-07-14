package ytime

import (
	"time"
)

// BeginMin return begin with minute
func BeginMin(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.Add(time.Duration(n[0]) * time.Minute)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

// EndMin return end with minute
func EndMin(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.Add(time.Duration(n[0]) * time.Minute)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 59, 999999999, t.Location())
}

// BeginHour return begin with hour
func BeginHour(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.Add(time.Duration(n[0]) * time.Hour)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

// EndHour return end with hour
func EndHour(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.Add(time.Duration(n[0]) * time.Hour)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 59, 999999999, t.Location())
}

// BeginDay return begin with day
func BeginDay(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.AddDate(0, 0, n[0])
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndDay return end with day
func EndDay(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.AddDate(0, 0, n[0])
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// BeginMon return begin with month
func BeginMon(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.AddDate(0, n[0], 0)
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndMon return end with month
func EndMon(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.AddDate(0, n[0], 0)
	}
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).Add(-1)
}

// BeginWeek return begin with week
func BeginWeek(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.AddDate(0, 0, n[0]*7)
	}
	return time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, t.Location())
}

// EndWeek return end with week
func EndWeek(t time.Time, n ...int) time.Time {
	if len(n) != 0 {
		t = t.AddDate(0, 0, n[0]*7)
	}
	return time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday())+6, 23, 59, 59, 999999999, t.Location())
}

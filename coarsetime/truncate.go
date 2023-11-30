package coarsetime

import (
	"time"
)

func truncate(d time.Duration) time.Time { return coarseTime.Load().(*time.Time).Truncate(d) }

func TruncateHour() time.Time {
	t := Now()
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location())
}
func TruncateMinute() time.Time { return truncate(time.Minute) }
func TruncateSecond() time.Time { return truncate(time.Second) }
func TruncateDay() time.Time {
	t := Now()
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func TruncateWeekday() time.Time {
	t := TruncateDay()
	return t.AddDate(0, 0, -int(t.Weekday()))
}

func TruncateMonth() time.Time {
	t := coarseTime.Load().(*time.Time)
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

func TruncateYear() time.Time {
	t := coarseTime.Load().(*time.Time)
	y, _, _ := t.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())
}

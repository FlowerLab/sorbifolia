package coarsetime

import (
	"testing"
)

func TestTruncateHour(t *testing.T) {
	t.Parallel()

	TruncateHour()
}

func TestTruncateMinute(t *testing.T) {
	t.Parallel()

	TruncateMinute()
}

func TestTruncateSecond(t *testing.T) {
	t.Parallel()

	TruncateSecond()
}

func TestTruncateDay(t *testing.T) {
	t.Parallel()

	TruncateDay()
}

func TestTruncateWeekday(t *testing.T) {
	t.Parallel()

	TruncateWeekday()
}

func TestTruncateMonth(t *testing.T) {
	t.Parallel()

	TruncateMonth()
}

func TestTruncateYear(t *testing.T) {
	t.Parallel()

	TruncateYear()
}

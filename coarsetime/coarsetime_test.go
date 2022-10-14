package coarsetime

import (
	"testing"
	"time"
)

func TestFloorTime(t *testing.T) {
	t.Parallel()

	FloorTime()
}

func TestCeilingTime(t *testing.T) {
	t.Parallel()

	CeilingTime()
}

func BenchmarkFloorTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FloorTime()
	}
}

func BenchmarkCeilingTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CeilingTime()
	}
}

func BenchmarkTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

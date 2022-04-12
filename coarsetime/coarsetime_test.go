package coarsetime

import (
	"testing"
)

func TestFloorTime(t *testing.T) {
	FloorTime()
}

func TestCeilingTime(t *testing.T) {
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

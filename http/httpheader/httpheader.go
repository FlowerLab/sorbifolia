package httpheader

type QualityValue struct {
	Value    []byte
	Priority float64 // 1.00 - 0.00
}

type (
	EachKeyQualityValue func(key []byte, val QualityValue) bool
	EachQualityValue    func(val QualityValue) bool
	EachValue           func(val []byte) bool
	EachKeyValue        func(key, val []byte) bool
	EachRanger          func(r Ranger) bool
)

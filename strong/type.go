package strong

type Type interface {
	Integer | UInteger | bool
}

type Integer interface {
	int | int8 | int16 | int32 | int64
}

type UInteger interface {
	uint | uint8 | uint16 | uint32 | uint64
}

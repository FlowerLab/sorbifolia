package strong

import (
	"testing"
)

func TestParse(t *testing.T) {
	t.Log(Parse[bool]("t"))
	t.Log(Parse[int64]("4487418"))
	t.Log(Parse[int32]("-12312"))
	t.Log(Parse[int]("-213123"))
	t.Log(Parse[uint]("4848"))
}

func TestFormat(t *testing.T) {
	t.Log(Format(true))
	t.Log(Format(1894189))
	t.Log(Format(-123123))
	t.Log(Format(uint8(123)))
}

func TestAppend(t *testing.T) {
	t.Log(Append(nil, true))
	t.Log(Append(nil, 4487418))
	t.Log(Append(nil, -12312))
	t.Log(Append(nil, uint16(4158)))
}

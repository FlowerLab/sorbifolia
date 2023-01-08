package httpmessage

type state uint8

const (
	_Init  state = 0
	_Close       = 1
	_Read        = 1 << (iota - 1)
	_Write
	_Method
	_URI
	_Version
	_Header
	_Body

	_Status = _URI
)

func (rs *state) Readable() bool {
	if *rs == _Init {
		*rs += _Read
	}
	return !rs.IsClose() && *rs&_Read == _Read
}
func (rs *state) Writable() bool {
	if *rs == _Init {
		*rs += _Write
	}
	return !rs.IsClose() && *rs&_Write == _Write
}
func (rs *state) Operate() state     { return *rs >> 3 << 3 }
func (rs *state) SetOperate(o state) { *rs = *rs<<5>>5 + o }
func (rs *state) IsClose() bool      { return *rs&_Close == _Close }
func (rs *state) Close()             { *rs = _Close }

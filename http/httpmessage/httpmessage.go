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

func (rs *state) Readable() bool     { return !rs.IsClose() && *rs&_Read == _Read }
func (rs *state) Writable() bool     { return !rs.IsClose() && *rs&_Write == _Write }
func (rs *state) Operate() state     { return *rs >> 3 << 3 }
func (rs *state) SetOperate(o state) { *rs = *rs<<5>>5 + o }
func (rs *state) IsClose() bool      { return *rs&_Close == _Close }
func (rs *state) Close()             { *rs = _Close }
func (rs *state) SetRead()           { rs.setRW(_Read) }
func (rs *state) SetWrite()          { rs.setRW(_Write) }
func (rs *state) setRW(s state) {
	if *rs != _Init && *rs&s == s {
		panic("BUG: not set state rw")
	}
	*rs += s
}

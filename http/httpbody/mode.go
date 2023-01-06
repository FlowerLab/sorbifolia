package httpbody

type rwcMode uint8

const (
	ModeReadWrite rwcMode = iota
	ModeWrite
	ModeRead
	ModeClose
)

func (m *rwcMode) SetMode(mode rwcMode) {
	switch *m { // 0:rw, 1:w, 2:r, 3:c
	case ModeReadWrite:
		*m = mode
	case ModeRead:
		panic("get io.Writer multiple times")
	case ModeWrite:
		panic("in read state")
	case ModeClose:
		panic("has been closed")
	default:
		panic("BUG: unknown state")
	}
}

func (m *rwcMode) IsMode(mode rwcMode) bool { return *m == mode }
func (m *rwcMode) IsReadWrite() bool        { return m.IsMode(ModeReadWrite) }
func (m *rwcMode) IsWrite() bool            { return m.IsMode(ModeWrite) }
func (m *rwcMode) IsRead() bool             { return m.IsMode(ModeRead) }
func (m *rwcMode) IsClose() bool            { return m.IsMode(ModeClose) }

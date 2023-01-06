package httpbody

type rwcMode uint8

const (
	ModeReadWrite rwcMode = iota
	ModeWrite
	ModeRead
	ModeClose
)

func (m *rwcMode) IsMode(mode rwcMode) bool { return *m == mode }
func (m *rwcMode) IsReadWrite() bool        { return m.IsMode(ModeReadWrite) }
func (m *rwcMode) IsWrite() bool            { return m.IsMode(ModeWrite) }
func (m *rwcMode) IsRead() bool             { return m.IsMode(ModeRead) }
func (m *rwcMode) IsClose() bool            { return m.IsMode(ModeClose) }

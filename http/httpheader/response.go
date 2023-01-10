package httpheader

type ResponseHeader struct {
	Header

	Close bool
}

func (rh *ResponseHeader) Reset() {
	rh.Header.Reset()
}

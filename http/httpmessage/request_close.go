package httpmessage

func (r *Request) Reset() {
	r.state = _Init
	r.buf.Reset()
	r.rp = 0
	r.bodyLength = 0
	r.Header.Reset()
	if r.Body != nil {
		_ = r.Body.Close()
		r.Body = nil
	}
}

func (r *Request) Close() error {
	r.buf.Reset()
	r.state.Close()
	return nil
}

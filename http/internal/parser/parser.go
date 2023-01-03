package parser

//
// type RequestParserFSM struct {
// 	Match func([]byte) int
// 	Next  int
// 	Limit int
// 	Call  func(r *Request, b *bufpool.Buffer) error
// 	More  func(r *Request, b *bufpool.Buffer) bool
// }
//
// var sss = [...]struct {
// 	Match func([]byte) int
// 	Next  int
// 	Limit int
// 	Call  func(r *Request, b *bufpool.Buffer) error
// 	More  func(r *Request, b *bufpool.Buffer) bool
// }{
// 	{
// 		Match: func(b []byte) int { return bytes.Index(b, char.Spaces) },
// 		Next:  1, Limit: 5,
// 		Call: func(r *Request, b *bufpool.Buffer) error {
// 			r.Method = method.Parse(b.Bytes())
// 			return nil
// 		},
// 	},
// 	{
// 		Match: func(b []byte) int {
// 			if i := bytes.Index(b, char.Spaces); i >= 0 {
// 				return i
// 			}
// 			return bytes.Index(b, char.CRLF)
// 		},
// 		Next: 1, Limit: 2048,
// 		Call: func(r *Request, b *bufpool.Buffer) error {
// 			r.Header.RequestURI = append(r.Header.RequestURI, b.Bytes()...)
// 			return nil
// 		},
// 	},
// 	{
// 		Match: func(b []byte) int { return bytes.Index(b, char.CRLF) },
// 		Limit: 10,
// 		Call: func(r *Request, b *bufpool.Buffer) error {
// 			if b.Len() == 0 {
// 				r.ver.Major, r.ver.Minor = 0, 9
// 				return nil
// 			}
//
// 			var ok bool
// 			if r.ver, ok = version.Parse(b.Bytes()); !ok {
// 				return httperr.ParseHTTPVersionErr
// 			}
// 			return nil
// 		},
// 	},
// 	{
// 		Match: func(b []byte) int { return bytes.Index(b, char.CRLF2) },
// 		Limit: 2048,
// 		Call: func(r *Request, buf *bufpool.Buffer) error {
// 			b := buf.Bytes()
// 			r.Header.KVs.preAlloc(bytes.Count(b, char.CRLF))
//
// 			for idx := bytes.Index(b, char.CRLF); len(b) > 0; idx = bytes.Index(b, char.CRLF) {
// 				r.Header.KVs.addHeader(b[:idx])
// 				b = b[idx+2:]
// 			}
//
// 			return r.Header.RawParse()
// 		},
// 	},
// }

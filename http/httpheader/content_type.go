package httpheader

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
)

type ContentType []byte

func (v ContentType) MIME() []byte {
	if i := bytes.IndexByte(v, char.Semi); i >= 0 {
		return v[:i]
	}
	return v
}

func (v ContentType) Charset() []byte {
	var charset []byte = nil

	if i := bytes.Index(v, char.Charset); i >= 0 {
		charset = v[i+len(char.Charset):]
	}
	if len(charset) == 0 || charset[0] != char.Equal {
		return nil
	}
	charset = charset[1:]

	if i := bytes.IndexByte(charset, char.Semi); i >= 0 {
		charset = charset[:i]
	}

	return cleanQuotationMark(cleanSuffixSpace(charset))
}

func (v ContentType) Boundary() []byte {
	var boundary []byte = nil

	if i := bytes.Index(v, char.Boundary); i >= 0 {
		boundary = v[i+len(char.Boundary):]
	}
	if len(boundary) == 0 || boundary[0] != char.Equal {
		return nil
	}
	boundary = boundary[1:]

	if i := bytes.IndexByte(boundary, char.Semi); i >= 0 {
		boundary = boundary[:i]
	}

	return cleanSuffixSpace(boundary)
}

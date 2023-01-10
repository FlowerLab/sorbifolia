package render

import (
	"encoding/json"

	"go.x2ox.com/sorbifolia/http/internal/bufpool"
)

var jsonContentType = []byte("application/json; charset=utf-8")

func JSON(data any) Render {
	r := &render{}
	buf := &bufpool.ReadBuffer{}
	err := json.NewEncoder(buf).Encode(data)
	r.length = int64(buf.Len())

	r.read = func(p []byte) (int, error) {
		if err != nil {
			return 0, err
		}
		return buf.Read(p)
	}
	r.contentType = jsonContentType
	return r
}

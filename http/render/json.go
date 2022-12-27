package render

import (
	"arena"
	"encoding/json"

	"go.x2ox.com/sorbifolia/http/internal/util"
)

var jsonContentType = []byte("application/json; charset=utf-8")

func JSON(a *arena.Arena, data any) Render {
	r := arena.New[render](a)
	buf := arena.New[util.Buffer](a)
	buf.A = a
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

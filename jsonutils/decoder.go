package jsonutils

import (
	"encoding/json"
	"io"
)

func Decode[T any](r io.Reader) (arr []T, err error) {
	dec := json.NewDecoder(r)

	for {
		arg := new(T)
		if err = dec.Decode(arg); err == nil {
			arr = append(arr, *arg)
		}

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
	}
}

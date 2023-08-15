package sa

import (
	"fmt"
	"log/slog"
)

func Error(key string, err error) (attr slog.Attr) {
	if err == nil {
		return String(key, "<nil>")
	}

	switch e := err.(type) {
	case fmt.Formatter:
		attr.Value = slog.StringValue(fmt.Sprintf("%+v", e))
	case slog.LogValuer:
		attr.Value = e.LogValue()
	default:
		attr.Value = slog.StringValue(err.Error())
	}

	return
}

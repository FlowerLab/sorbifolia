package httputils

import (
	"errors"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

func TestHTTP_ParseData(t *testing.T) {
	t.Parallel()

	var resp X2oxIPResp
	err := Get("https://api.ip.x2ox.com/api/ip?type=json").
		SetContentType(AppJSON).
		Request(1, nil, 5*time.Second).
		ParserData(JSONParser(&resp)).
		DoRelease()
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.IP)

	if err = Get("https://api.ip.x2ox.com/api/ip?type=json").
		SetContentType(AppJSON).
		Request(1, nil, 5*time.Second).
		ParserData(func(resp *fasthttp.Response) error { return errors.New("err") }).
		DoRelease(); err == nil {
		t.Error("err should")
	}
}

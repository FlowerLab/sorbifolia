package httputils

import (
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

type X2oxIPResp struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}

func TestHTTP(t *testing.T) {
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
}

func TestHTTPRetry(t *testing.T) {
	t.Parallel()

	var resp X2oxIPResp
	err := Get("https://api.ip.x2ox.com/api/ip?type=json").
		SetContentType(AppJSON).
		Request(-1, nil, 5*time.Second).
		ParserData(JSONParser(&resp)).
		DoRelease()

	if err != nil {
		t.Error(err)
	}
	t.Log(resp.IP)
}

func TestHTTPError(t *testing.T) {
	t.Parallel()

	var resp X2oxIPResp
	err := Post("https://api.ip.x2ox.com/api/ip?type=json").
		SetContentType(AppJSON).
		SetBodyWithEncoder(JSON(), struct {
			A chan struct{}
		}{A: nil}).
		Request(3, func(err error, response *fasthttp.Response) bool {
			return true
		}, 5*time.Second).
		ParserData(JSONParser(&resp)).
		DoRelease()

	if err == nil {
		t.Error("err")
	}
}

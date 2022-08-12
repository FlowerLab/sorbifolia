package httputils

import (
	"testing"
	"time"
)

type X2oxIPResp struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}

func TestHTTP(t *testing.T) {
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

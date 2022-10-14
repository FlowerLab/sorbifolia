package httputils

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestHTTP_SetClient(t *testing.T) {
	t.Parallel()

	h := Get().SetClient(&fasthttp.Client{Name: "foo"})

	if h.client.Name != "foo" {
		t.Error("SetClient err")
	}
}

func TestHTTP_SetProxy(t *testing.T) {
	t.Parallel()

	if h := Get().SetProxy("socks5://proxy.x2ox.com:8888"); h.client == nil ||
		h.client.Dial == nil {
		t.Error("SetProxy err")
	}

	if h := Get().SetProxy("socks5h://proxy.x2ox.com:8888"); h.client == nil ||
		h.client.Dial == nil {
		t.Error("SetProxy err")
	}

	if h := Get().SetProxy("http://proxy.x2ox.com:8888"); h.client == nil ||
		h.client.Dial == nil {
		t.Error("SetProxy err")
	}

	if h := Get().SetProxy("https://proxy.x2ox.com:8888"); h.client == nil ||
		h.client.Dial == nil {
		t.Error("SetProxy err")
	}

	if h := Get().SetProxy("un://proxy.x2ox.com:8888"); h.client == nil ||
		h.client.Dial != nil {
		t.Error("SetProxy err")
	}

	if h := Get().SetProxy("un:// .x2ox.com:8888"); h.client == nil ||
		h.client.Dial != nil {
		t.Error("SetProxy err")
	}
}

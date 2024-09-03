package crpc

import (
	"net/http"
	"path"
	"sync/atomic"

	"github.com/VictoriaMetrics/metrics"
)

type Health interface {
	IsLive() bool
	IsReady() bool

	SetLive()
	SetNoLive(e string)

	SetReady()
	SetNoReady(e string)
}

type healthAndMetrics struct {
	live, ready atomic.Value
	addr, path  string
}

func (x *healthAndMetrics) IsLive() bool  { return x.live.Load() == nil }
func (x *healthAndMetrics) IsReady() bool { return x.ready.Load() == nil }

func (x *healthAndMetrics) SetLive()           { x.live.Store(nil) }
func (x *healthAndMetrics) SetNoLive(e string) { x.live.Store(e) }

func (x *healthAndMetrics) SetReady()           { x.ready.Store(nil) }
func (x *healthAndMetrics) SetNoReady(e string) { x.ready.Store(e) }

func (x *healthAndMetrics) Register(h HttpHandle) {
	h.HandleFunc(path.Join(x.path, "/metrics"), func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, true)
	})

	h.HandleFunc(path.Join(x.path, "/livez"), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		if x.IsLive() {
			_, _ = w.Write([]byte("ok"))
		} else {
			val, _ := x.ready.Load().(string)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte(val))
		}
	})

	h.HandleFunc(path.Join(x.path, "/readyz"), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		if x.IsReady() {
			_, _ = w.Write([]byte("ok"))
		} else {
			val, _ := x.ready.Load().(string)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte(val))
		}
	})

	h.HandleFunc(path.Join(x.path, "/"), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<h2>Server: core-pb/tag</h2><br/>
	<a href="livez">livez</a> - liveness checking<br/>
	<a href="readyz">readyz</a> - readiness checking<br/>
	<a href="metrics">metrics</a> - available service metrics<br/>
	`))
	})
}

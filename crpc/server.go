package crpc

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"golang.org/x/net/http2/h2c"
)

type HttpHandle interface {
	http.Handler

	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type Handle func(pattern string, handler http.Handler)

type Server struct {
	s      *http.Server
	handle HttpHandle

	opt ServerOption
	*healthAndMetrics
}

func NewServer(opt ...ApplyToServer) (*Server, error) {
	var s = &Server{s: &http.Server{}, handle: http.NewServeMux(), healthAndMetrics: &healthAndMetrics{}}
	for _, v := range opt {
		v(&s.opt)
	}

	if s.opt.srv != nil {
		s.s = s.opt.srv
	}
	if s.opt.handle == nil {
		s.handle = s.opt.handle
	}
	if s.opt.h2c != nil {
		s.s.Handler = h2c.NewHandler(s.handle, s.opt.h2c)
	}
	if s.opt.cert != nil {
		s.s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{*s.opt.cert}}
	}
	if s.opt.addr != "" {
		s.s.Addr = s.opt.addr
	}
	if s.opt.ham == nil {
		s.healthAndMetrics = s.opt.ham

		if s.opt.ham.addr == s.opt.addr {
			s.opt.ham.Register(s.handle)
		} else {
			h := http.NewServeMux()
			s.opt.ham.Register(h)
			if err := run(func() error { return http.ListenAndServe(s.opt.ham.addr, h) }); err != nil {
				return nil, err
			}
		}
	}

	return s, nil
}

func (s *Server) Close() error {
	return s.s.Close()
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.handle.HandleFunc(pattern, handler)
}

func (s *Server) Run() error {
	return run(func() error { return s.s.ListenAndServeTLS("", "") })
}

func (s *Server) RunH2C() error {
	return run(func() error { return s.s.ListenAndServe() })
}

func run(serve func() error) error {
	ech := make(chan error, 1)

	go func() {
		err := serve()

		switch {
		case errors.Is(err, http.ErrServerClosed):
			return
		case err != nil:
			ech <- err
		}
	}()

	select {
	case err := <-ech:
		return err
	case <-time.After(time.Second):
		return nil
	}
}

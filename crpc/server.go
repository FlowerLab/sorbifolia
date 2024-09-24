package crpc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/http2/h2c"
)

type HttpHandle interface {
	http.Handler

	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type Server struct {
	s      *http.Server
	handle HttpHandle

	*healthAndMetrics
}

func NewServer(opt ...ApplyToServer) (*Server, error) {
	var (
		handle = http.NewServeMux()
		s      = &Server{s: &http.Server{Handler: handle}, handle: handle, healthAndMetrics: &healthAndMetrics{}}
		so     = &ServerOption{}
	)

	for _, v := range opt {
		v(so)
	}

	if so.srv != nil {
		s.s = so.srv
	}
	if so.handle != nil {
		s.handle = so.handle
		s.s.Handler = s.handle
	}
	if so.h2c != nil {
		s.s.Handler = h2c.NewHandler(s.handle, so.h2c)
	}
	if so.cors != nil {
		s.s.Handler = so.cors(s.s.Handler)
	}
	if so.cert != nil {
		s.s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{*so.cert}, MinVersion: tls.VersionTLS12}
	}
	if so.addr != "" {
		s.s.Addr = so.addr
	}
	if so.ham != nil {
		s.healthAndMetrics = so.ham

		if so.addr != "" && so.ham.addr == so.addr {
			so.ham.Register(s.handle)
		} else {
			h := http.NewServeMux()
			so.ham.Register(h)
			if err := run(func() error { return http.ListenAndServe(so.ham.addr, h) }); err != nil {
				return nil, err
			}
		}
	}

	return s, nil
}

func (s *Server) Close() error {
	s.SetNoLive("server is closing")
	err := s.s.Close()
	s.SetNoLive("server is closed")
	return err
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.handle.Handle(pattern, handler)
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.handle.HandleFunc(pattern, handler)
}

func (s *Server) Run() error {
	s.SetNoLive("server is starting")

	if err := run(func() error { return s.s.ListenAndServeTLS("", "") }); err != nil {
		s.SetNoLive(fmt.Sprintf("server start fail: %s", err))
		return err
	}
	s.SetLive()
	return nil
}

func (s *Server) RunH2C() error {
	s.SetNoLive("server is starting")

	if err := run(func() error { return s.s.ListenAndServe() }); err != nil {
		s.SetNoLive(fmt.Sprintf("server start fail: %s", err))
		return err
	}

	s.SetLive()
	return nil
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

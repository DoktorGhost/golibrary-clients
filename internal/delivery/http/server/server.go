package server

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server *http.Server
	errors chan error
}

func NewHttpServer(hanler http.Handler, port string) *Server {
	if port == "" {
		port = _defaultAddr
	}

	httpServer := &http.Server{
		Handler: hanler,
		Addr:    port,
	}

	return &Server{
		server: httpServer,
		errors: make(chan error),
	}
}

func (s Server) Serve() {
	go func() {
		s.errors <- s.server.ListenAndServe()
		close(s.errors)
	}()
}

func (s Server) Notify() <-chan error {
	return s.errors
}

func (s Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

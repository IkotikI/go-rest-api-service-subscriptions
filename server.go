package microservice

import (
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(httpServer *http.Server) *Server {
	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}

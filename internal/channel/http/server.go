package http

import (
	"context"
	"github.com/DeOne4eg/http-multiplexer/config"
	"net/http"
	"strconv"
)

// Server is simple HTTP server
type Server struct {
	httpServer *http.Server
}

// NewServer create instance of Server
func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + strconv.Itoa(cfg.HTTP.Port),
			Handler:      handler,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
		},
	}
}

// Run is starting HTTP server
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Stop is stopping HTTP server
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

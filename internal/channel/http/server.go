package http

import (
	"fmt"
	"github.com/DeOne4eg/http-multiplexer/config"
	"net/http"
)

type Server struct {
	httpServer *http.ServeMux
}

func NewServer(handler *http.Handler) *Server {
	return &Server{
		httpServer: http.NewServeMux(),
	}
}

func (s *Server) Run(cfg *config.Config) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port), s.httpServer)
}
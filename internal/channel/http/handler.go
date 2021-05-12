package http

import (
	"github.com/DeOne4eg/http-multiplexer/config"
	"net/http"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) userLogin(c *gin.Context) {

}

func (h *Handler) handlePing(w http.ResponseWriter, r *http.Request) {
	write, err := w.Write([]byte("pong"))
	if err != nil {
		return 
	}
}

func (h *Handler) handleUrls(w http.ResponseWriter, r *http.Request) {

}
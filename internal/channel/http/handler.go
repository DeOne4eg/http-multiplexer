package http

import (
	"encoding/json"
	"net/http"
)

const (
	maxConn = 100
)

var (
	limiter = make(chan struct{}, maxConn)
)

type Handler struct {}

type input struct {
	URL []string `json:"urls"`
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Init() http.Handler {
	return http.HandlerFunc(h.handle)
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	limiter <- struct{}{}
	defer func() { <-limiter }()

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		h.handlePing(w, r)
	case "POST":
		h.handleUrls(w, r)
	}
}

func (h *Handler) handlePing(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func (h *Handler) handleUrls(w http.ResponseWriter, r *http.Request) {
	i := input{}
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, "Error parsing JSON")
		return
	}

	switch {
	case len(i.URL) == 0:
		NewErrorResponse(w, http.StatusBadRequest, "Empty URL list")
		return
	case len(i.URL) > 20:
		NewErrorResponse(w, http.StatusBadRequest, "The number of URLs should not be more than 20")
		return
	}

	_, _ = w.Write([]byte("ok"))
}
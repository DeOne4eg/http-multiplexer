package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
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

type urlResponse struct {
	URL string `json:"url"`
	StatusCode int `json:"status_code"`
	Response string `json:"response"`
	Headers map[string]string `json:"headers"`
}

type successResponse struct {
	OK bool `json:"ok"`
	Result []urlResponse `json:"result"`
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
		NewErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
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

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan urlResponse, 4)

	defer cancel()

	for _, u := range i.URL {
		_, err := url.ParseRequestURI(u)
		if err != nil {
			NewErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid URL: %s", u))
			return
		}

		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				fmt.Println("Goroutine stopped")
				return
			default:
				ch <- urlResponse{
					URL:        "https://vk.com",
					StatusCode: 200,
					Response:   "test",
					Headers:    nil,
				}
			}
		}(ctx)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	var res []urlResponse
	for ur := range ch {
		if ur.StatusCode != 200 {
			NewErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error handling URL: %v", ur.URL))
			return
		}

		res = append(res, ur)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&successResponse{
		OK:   true,
		Result: res,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
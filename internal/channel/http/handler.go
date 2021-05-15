package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	// limit connections to server
	maxConn = 100
)

var (
	// buff channel for counting connections
	limiter = make(chan struct{}, maxConn)
)

// Handler is a struct for handle request
type Handler struct {}

// input describe JSON struct of HTTP request body
type input struct {
	URL []string `json:"urls"`
}

// urlResponse contains information about result of request to received URLs
type urlResponse struct {
	URL string `json:"url"`
	StatusCode int `json:"status_code"`
	Response string `json:"response"`
	Headers map[string]string `json:"headers"`
}

// successResponse returns on successful result
type successResponse struct {
	OK bool              `json:"ok"`
	Result []urlResponse `json:"result"`
}

// NewHandler create instance of Handler
func NewHandler() *Handler {
	return &Handler{}
}

// Init calls the handler function and return http.Handler
func (h *Handler) Init() http.Handler {
	return http.HandlerFunc(h.handle)
}

// handle is main handler for all HTTP requests.
// If request method is POST then calls handleUrls function.
// If request method is not POST then calls handlePing function.
func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	limiter <- struct{}{}
	defer func() { <-limiter }()

	log.Printf("%s request: %s", r.Method, r.URL)
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		h.handleUrls(w, r)
	default:
		h.handlePing(w, r)
	}
}

// handlePing simple ping pong handler
func (h *Handler) handlePing(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

//handleUrls is main function for handle received URLs.
func (h *Handler) handleUrls(w http.ResponseWriter, r *http.Request) {
	i := input{}
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// input data validation in switch case
	switch {
	case len(i.URL) == 0:
		NewErrorResponse(w, http.StatusBadRequest, "Empty URL list")
		return
	case len(i.URL) > 20:
		NewErrorResponse(w, http.StatusBadRequest, "The number of URLs should not be more than 20")
		return
	}

	var wg sync.WaitGroup

	// create context with cancel for stop all goroutines
	ctx, cancel := context.WithCancel(r.Context())

	// main channel for getting handle result of URL
	// 4 - limit goroutines
	ch := make(chan urlResponse, 4)

	defer cancel()

	for _, u := range i.URL {
		_, err := url.ParseRequestURI(u)
		if err != nil {
			NewErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid URL: %s", u))
			return
		}

		wg.Add(1)
		go func(ctx context.Context, url string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				// exit the goroutine
				return
			default:
				// create client with timeout 1 second
				client := http.Client{Timeout: 1 * time.Second}
				// do GET request
				resp, err := client.Get(url)
				if err != nil {
					ch <- urlResponse{
						URL:        url,
						StatusCode: 0,
						Response:   "",
						Headers:    nil,
					}
					return
				}
				defer func() {
					_ = resp.Body.Close()
				}()

				// get body
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error read body; %v", err)
				}

				// get headers from response
				headers := make(map[string]string)
				getHeaders(resp.Header, headers)

				// write result
				ch <- urlResponse{
					URL:        url,
					StatusCode: resp.StatusCode,
					Response:   string(body),
					Headers:    headers,
				}
			}
		}(ctx, u)
	}

	go func() {
		// close channel
		defer close(ch)
		// wait the end of all goroutines
		wg.Wait()
	}()

	var res []urlResponse
	for ur := range ch {
		// if server return error then throw errorResponse with status code 500
		if ur.StatusCode < 200 || ur.StatusCode >= 300 {
			NewErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error handling URL: %v", ur.URL))
			return
		}

		// write result in urlResponse slice
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

// getHeaders receive HTTP headers list and convert it to map
func getHeaders(in http.Header, out map[string]string)  {
	for k, v := range in {
		out[k] = v[0]
	}
}
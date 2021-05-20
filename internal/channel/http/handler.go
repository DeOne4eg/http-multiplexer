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

// Handler is a struct for handle request.
type Handler struct{}

// input describe JSON struct of HTTP request body.
type input struct {
	URL []string `json:"urls"`
}

// urlResponse contains information about result of request to received URLs.
type urlResponse struct {
	URL        string            `json:"url"`
	StatusCode int               `json:"status_code"`
	Response   string            `json:"response"`
	Headers    map[string]string `json:"headers"`
}

// successResponse returns on successful result.
type successResponse struct {
	OK     bool          `json:"ok"`
	Result []urlResponse `json:"result"`
}

// NewHandler create instance of Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Init calls the handler function and return http.Handler.
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
	case http.MethodPost:
		h.handleUrls(w, r)
	default:
		h.handlePing(w, r)
	}
}

// handlePing simple ping pong handler.
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
				return
			default:
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					ch <- urlResponse{URL: url}
					log.Fatalf("Error creating request: %v", err)
				}
				req = req.WithContext(ctx)

				// create HTTP client with timeout 1 second
				client := http.Client{Timeout: 1 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					ch <- urlResponse{URL: url}
					return
				}
				defer func() {
					err = resp.Body.Close()
					if err != nil {
						log.Printf("Error close body response: %v", err.Error())
					}
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
		defer close(ch)
		wg.Wait()
	}()

	var res []urlResponse
	for ur := range ch {
		// if server return error then throw errorResponse with status code 500
		if ur.StatusCode < http.StatusOK || ur.StatusCode >= http.StatusMultipleChoices {
			NewErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error handling URL: %v", ur.URL))
			return
		}

		// write result in urlResponse slice
		res = append(res, ur)
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&successResponse{
		OK:     true,
		Result: res,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// getHeaders receive HTTP headers list and convert it to map
func getHeaders(in http.Header, out map[string]string) {
	for k, v := range in {
		out[k] = v[0]
	}
}

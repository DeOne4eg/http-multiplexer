package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleUrls(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	handler := NewHandler()

	ts := httptest.NewServer(http.HandlerFunc(handler.handleUrls))
	defer ts.Close()

	tests := []struct {
		name              string
		ok                bool
		payload           string
		expectStatusCode  int
		expectDescription string
	}{
		{
			name:             "correct with one url",
			ok:               true,
			payload:          `{"urls":["https://jsonplaceholder.typicode.com/posts/1"]}`,
			expectStatusCode: 200,
		},
		{
			name:              "invalid json",
			ok:                false,
			expectStatusCode:  400,
			expectDescription: "Invalid JSON",
		},
		{
			name:              "invalid json",
			ok:                false,
			payload:           `{`,
			expectStatusCode:  400,
			expectDescription: "Invalid JSON",
		},
		{
			name:              "more than 20 urls",
			ok:                false,
			payload:           `{"urls":["https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1","https://jsonplaceholder.typicode.com/posts/1"]}`,
			expectStatusCode:  400,
			expectDescription: "The number of URLs should not be more than 20",
		}, {
			name:              "more than 20 urls",
			ok:                false,
			payload:           `{"urls":[]}`,
			expectStatusCode:  400,
			expectDescription: "Empty URL list",
		},
	}

	for _, tc := range tests {
		res, err := http.Post(ts.URL, "application/json", strings.NewReader(tc.payload))
		if err != nil {
			t.Error(err)
		}
		defer func() {
			_ = res.Body.Close()
		}()

		if res.StatusCode != tc.expectStatusCode {
			t.Fatalf("Test failed: %s. Expected %d, got %d", tc.name, tc.expectStatusCode, res.StatusCode)
		}

		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		if tc.ok {
			response := successResponse{}
			err := json.Unmarshal(content, &response)
			if err != nil {
				t.Fatalf("Error unmarshal json")
				return
			}

			if response.OK != tc.ok {
				t.Fatalf("Test failed: %s. Expected %v, got %v", tc.name, tc.ok, response.OK)
			}
		} else {
			response := errorResponse{}
			err := json.Unmarshal(content, &response)
			if err != nil {
				t.Fatalf("Error unmarshal json")
				return
			}

			if response.OK != tc.ok {
				t.Fatalf("Test failed: %s. Expected %v, got %v", tc.name, tc.ok, response.OK)
			} else if response.Description != tc.expectDescription {
				t.Fatalf("Test failed: %s. Expected %s, got %s", tc.name, tc.expectDescription, response.Description)
			}
		}
	}
}

package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleUrls(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	handler := NewHandler()

	ts := httptest.NewServer(http.HandlerFunc(handler.handleUrls))
	defer ts.Close()
}
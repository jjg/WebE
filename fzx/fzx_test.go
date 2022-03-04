package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jjg/WebE/fzx/quiet"
)

func TestRequestMethods(t *testing.T) {

	// Mute output while running tests.
	defer quiet.BeQuiet()()

	cases := []struct {
		method string
		url    string
	}{
		{http.MethodHead, "http://example.com"},
		{http.MethodGet, "http://example.com"},
		{http.MethodPost, "http://example.com"},
		{http.MethodPut, "http://example.com"},
		{http.MethodDelete, "http://example.com"},
		// NOTE: EXECUTE is a made-up method specific to fzx,
		// maybe we'll create an RFC for it once this takes off...
		{"EXECUTE", "http://example.com"},
	}

	for _, c := range cases {
		req := httptest.NewRequest(c.method, c.url, nil)
		w := httptest.NewRecorder()

		Handler(w, req)

		if want, got := http.StatusOK, w.Result().StatusCode; want != got {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

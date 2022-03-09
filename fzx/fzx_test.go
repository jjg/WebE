package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jjg/WebE/fzx/quiet"
	"github.com/stretchr/testify/assert"
)

func TestReadMethods(t *testing.T) {

	// Mute output while running tests.
	defer quiet.BeQuiet()()

	cases := []struct {
		method string
		url    string
	}{
		{http.MethodHead, "http://example.com"},
		{http.MethodGet, "http://example.com"},
		//{http.MethodPost, "http://example.com"},
		//{http.MethodPut, "http://example.com"},
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

func TestPostPut(t *testing.T) {

	testFileContents := "A plain text file to test the POST and PUT methods."

	// Write testFileContents to a file.
	f, err := os.CreateTemp("", "posttest.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	if _, err := f.Write([]byte(testFileContents)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	// POST file to /testing/posttest.txt
	req := httptest.NewRequest(http.MethodPost, "http://localhost:7302/testing/posttest.txt", f)
	w := httptest.NewRecorder()

	Handler(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Make a HEAD request for /testing/posttest.txt and compare header values.
	req = httptest.NewRequest(http.MethodHead, "http://localhost:7302/testing/posttest.txt", nil)

	Handler(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Make a GET request for /testing/posttest.txt and compare the contents.
	req = httptest.NewRequest(http.MethodGet, "http://localhost:7302/testing/posttest.txt", nil)

	Handler(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, testFileContents, w.Body.String())
}

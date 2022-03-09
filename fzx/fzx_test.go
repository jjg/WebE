package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jjg/WebE/fzx/inode"
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

	testFileUrl := "http://localhost:7302/testing/posttest.txt"
	testFileFzxPath := ".localhost:7302/testing/posttest.txt"
	testFileContents := "A plain text file to test the POST and PUT methods."

	// Write testFileContents to a file.
	f, err := os.CreateTemp("", "*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	if _, err := f.Write([]byte(testFileContents)); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	// POST the file.
	req := httptest.NewRequest(http.MethodPost, testFileUrl, f)
	w := httptest.NewRecorder()

	Handler(w, req)

	// POST should return a string of JSON data with details about what was stored.
	// To use this for testing, we need to extract the body and parse the JSON.
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	postRequestResultBody := buf.String()
	i := &inode.Inode{}
	if err = json.Unmarshal([]byte(postRequestResultBody), i); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	assert.Equal(t, testFileFzxPath, i.FzxPath)
	// TODO: Consider additional checks to validate POST response JSON.

	// Make a HEAD request and test the metadata.
	req = httptest.NewRequest(http.MethodHead, testFileUrl, nil)

	Handler(w, req)

	log.Print(w.Result().Header["FzxPath"])

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, i.FzxPath, w.Result().Header["FzxPath"])
	// TODO: Compare additional values returned by POST request to headers.

	/*
		// Make a GET request and test the contents.
		req = httptest.NewRequest(http.MethodGet, testFileUrl, nil)

		Handler(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, testFileContents, w.Body.String())
	*/
}

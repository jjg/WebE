package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jjg/WebE/fzx/inode"
	"github.com/jjg/WebE/fzx/utils"
	"github.com/stretchr/testify/assert"
)

func TestHttpMethods(t *testing.T) {

	// Mute output while running tests.
	defer utils.BeQuiet()()

	// Setup test file parameters.
	testFilename := fmt.Sprintf("%v.txt", time.Now().Unix())
	testFileUrl := fmt.Sprintf("http://localhost:7302/testing/%v", testFilename)
	testFileFzxPath := fmt.Sprintf(".localhost:7302/testing/%v", testFilename)
	testFileContents := "A plain text file to test the POST and PUT methods."
	testFileContentType := "Text/Plain"
	testFileLength := len(testFileContents)

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

	// Re-open the file so we can POST it.
	// TODO: There *must* be a better way to do this...
	f2, _ := os.Open(f.Name())
	defer f2.Close()

	t.Run("POST", func(t *testing.T) {

		// POST the file.
		req := httptest.NewRequest(http.MethodPost, testFileUrl, f2)
		req.Header.Set("Content-Type", testFileContentType)
		postRecorder := httptest.NewRecorder()

		Handler(postRecorder, req)

		// POST should return a string of JSON data with details about what was stored.
		// To use this for testing, we need to extract the body and parse the JSON.
		buf := new(bytes.Buffer)
		buf.ReadFrom(postRecorder.Result().Body)
		postRequestResultBody := buf.String()
		i := &inode.Inode{}
		if err = json.Unmarshal([]byte(postRequestResultBody), i); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusCreated, postRecorder.Result().StatusCode)
		assert.Equal(t, testFileFzxPath, i.FzxPath)
		// TODO: Consider additional checks to validate POST response JSON.
	})

	t.Run("HEAD", func(t *testing.T) {

		// Make a HEAD request and test the metadata.
		req := httptest.NewRequest(http.MethodHead, testFileUrl, nil)
		headRecorder := httptest.NewRecorder()

		Handler(headRecorder, req)

		assert.Equal(t, http.StatusOK, headRecorder.Result().StatusCode)
		assert.Equal(t, testFileContentType, headRecorder.Result().Header["Content-Type"][0])
		// NOTE: Using Sprintf to convert between int and string seems dumb.
		assert.Equal(t, fmt.Sprintf("%v", testFileLength), headRecorder.Result().Header["Content-Length"][0])
	})

	t.Run("GET", func(t *testing.T) {
		// Make a GET request and test the contents.
		req := httptest.NewRequest(http.MethodGet, testFileUrl, nil)
		getRecorder := httptest.NewRecorder()

		Handler(getRecorder, req)

		assert.Equal(t, http.StatusOK, getRecorder.Result().StatusCode)
		assert.Equal(t, testFileContents, getRecorder.Body.String())
	})

	// TODO: Test PUT.
	// TODO: Test EXECUTE.
	// TODO: Test DELETE.
}

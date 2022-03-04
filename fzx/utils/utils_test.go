package utils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFzxPathFromRequest(t *testing.T) {

	cases := []struct {
		url     string
		fzxPath string
	}{
		{"http://example.com", ".com.example"},
		{"http://example.com/", ".com.example/"},
		{"http://example.com/foo/bar.html", ".com.example/foo/bar.html"},
		{"http://example.com/.com.debug.www/baz/fizz.json", ".com.debug.www/baz/fizz.json"},
	}

	for _, c := range cases {

		req, _ := http.NewRequest(http.MethodGet, c.url, nil)
		target, _ := FzxPathFromRequest(req)

		assert.Equal(t, c.fzxPath, target)
	}
}

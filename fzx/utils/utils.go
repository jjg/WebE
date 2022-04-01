package utils

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// TODO: Write a test for this function.
func StringToSha1(in string) string {
	h := sha1.New()
	h.Write([]byte(in))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// TODO: This could probably be combined with
// StringToSha1 into something more generic.
func BytesToSha1(in []byte) string {
	h := sha1.New()
	h.Write(in)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func Reverse(input []string) {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
}

func FzxPathFromRequest(r *http.Request) (string, error) {

	var err error
	path := r.URL.Path

	// If the caller provided a native fzx path, use it.
	if len(path) > 1 {
		if string(path[1]) == "." {
			return string(path[1:]), err
		}
	}

	// Split the hostname appart.
	splitHost := strings.Split(r.Host, ".")

	// Add one more element so we get a leading "."
	splitHost = append(splitHost, "")

	// Reverse the array (slice?)
	Reverse(splitHost)

	// Finally glue it together backwards.
	reversedHost := strings.Join(splitHost, ".")

	// Add the path and return.
	return fmt.Sprintf("%v%v", reversedHost, path), err
}

// Hides output from fmt, log, etc.
func BeQuiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}

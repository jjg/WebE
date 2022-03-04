package utils

import (
	"fmt"
	"net/http"
	"strings"
)

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

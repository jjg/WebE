package main

import (
	"io"
	"log"
	"net/http"
)

// TODO: Move these to flags, env, etc.
const (
	LISTEN_PORT      = ":7302"
	STORAGE_LOCATION = "./blocks"
)

func Handler(w http.ResponseWriter, req *http.Request) {

	// DEBUG
	log.Print(req)

	// TODO: Translate DNS name to fzx namespace.
	// TODO: Try to load inode.
	// TODO: Check authorization.
	// TODO: Set default response status code.
	// TODO: Set default response headers.

	switch req.Method {
	case "HEAD":
		log.Print("Got HEAD")
		// TODO: Add/update response-specific headers, status.
	case "GET":
		log.Print("Got GET")
		// TODO: Locate blocks.
		// TODO: Add/update response-specific headers, status.
		// TODO: Return blocks.
	case "POST":
		log.Print("Got POST")
		// TODO: Write blocks.
		// TODO: Write inode.
		// TODO: Add/update response-specific headers, status.
	case "PUT":
		log.Print("Got PUT")
		// TODO: Handle PUT request
		// NOTE: This is probably identical to POST.
		// TODO: Write blocks.
		// TODO: Write inode.
		// TODO: Add/update response-specific headers, status.
	case "DELETE":
		log.Print("Got DELETE")
		// TODO: Delete inode
		// TODO: Add/update response-specific headers, status.
	case "EXECUTE":
		log.Print("Got EXECUTE")
		// TODO: Handle EXECUTE request.
		// TODO: Execute specified binary.
		// TODO: Return output.
		// TODO: Add/update response-specific headers, status.
	default:
		log.Printf("I don't know what to do with method %v", req.Method)
	}

	// TODO: Finalize response (flush buffers, etc.).
	io.WriteString(w, "Hello, world!\n")
}

func main() {

	log.Printf("fzx listening on port %v", LISTEN_PORT)

	// Wire-up handler.
	http.HandleFunc("/", Handler)

	// Listen for incoming HTTP requests.
	// NOTE: This blocks anything below it from running.
	log.Fatal(http.ListenAndServe(LISTEN_PORT, nil))
}

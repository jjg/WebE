package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jjg/WebE/fzx/inode"
	"github.com/jjg/WebE/fzx/methods"
	"github.com/jjg/WebE/fzx/utils"
)

const (
	DEFAULT_LISTEN_PORT      = 7302
	DEFAULT_STORAGE_LOCATION = "blocks"
	DEFAULT_BLOCK_SIZE       = 1024 * 1024
)

// Get port, data directory from command line.
var listenPort = flag.Int("p", DEFAULT_LISTEN_PORT, "Override the default port")
var blockSize = flag.Int64("bs", DEFAULT_BLOCK_SIZE, "Block size in bytes")
var storageLocation = flag.String("storage", DEFAULT_STORAGE_LOCATION, "Block storage location")

func Handler(w http.ResponseWriter, req *http.Request) {

	var err error
	var fzxPath string

	// Translate DNS name to fzx namespace.
	if fzxPath, err = utils.FzxPathFromRequest(req); err != nil {
		log.Printf("Error extracting fzx path from request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Try to load inode (it's OK if this fails for POST/PUT/etc.).
	// TODO: I'm not sure anInode is getting reinitialized as expected during
	// each request cycle.  Tests seem to show that values "stick" between
	// requests and I'm not sure if this is just a side-effect of the test
	// harness or a problem with the way I'm handling pointers.  Something
	// to look into.
	anInode := &inode.Inode{FzxPath: fzxPath, StorageLocation: storageLocation}
	if err := anInode.Load(storageLocation, fzxPath); err != nil {
		log.Printf("Error loading inode for %v, %v", fzxPath, err)
		anInode.BlockSize = *blockSize
		anInode.StorageLocation = storageLocation

		// Set additional inode metadata based on request.
		v, ok := req.Header["Content-Type"]
		if ok {
			anInode.ContentType = v[0]
		}
	} else {
		log.Printf("Loaded inode for %v", fzxPath)

		// Set response headers using inode data.
		w.Header().Add("Content-Type", anInode.ContentType)
		w.Header().Add("Content-Length", fmt.Sprintf("%v", anInode.FileSize))
		w.Header().Add("FzxPath", anInode.FzxPath)
	}

	// TODO: Process request authorization.

	// TODO: I don't love how the response is sometimes updated in the handler
	// and sometimes updated here.  Consider other approaches that keep the HTTP
	// stuff limited to this layer.
	switch req.Method {
	case "HEAD":
		log.Print("Got HEAD")
		if anInode.Status != http.StatusNotFound {
			methods.Head(w, req, anInode)
		}
	case "OPTIONS":
		log.Print("Got OPTIONS")
		// TODO: Does OPTIONS care if we have an inode or not?
		methods.Options(w, req, anInode)
	case "GET":
		log.Print("Got GET")
		if anInode.Status != http.StatusNotFound {
			methods.Get(w, req, anInode)
		}
	case "POST":
		log.Print("Got POST")
		// POST is not be allowed if an inode already exists!
		if anInode.Status != http.StatusNotFound {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Can't POST over an existing file, try PUT instead."))
		} else {
			methods.Post(w, req, anInode)
		}
	case "PUT":
		log.Print("Got PUT")
		if anInode.Status == http.StatusNotFound {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Can't PUT unless file exists, try POST instead."))
		} else {
			// TODO: Abstract storage into something shared by both
			// POST and PUT; for now, just use the POST function for both.
			methods.Post(w, req, anInode)
			//methods.Put(w, req, anInode)
		}
	case "DELETE":
		log.Print("Got DELETE")
		if anInode.Status != http.StatusNotFound {
			methods.Delete(w, req, anInode)
		}
	case "EXECUTE":
		log.Print("Got EXECUTE")
		if anInode.Status != http.StatusNotFound {
			methods.Execute(w, req, anInode)
		}
	default:
		log.Printf("I don't know what to do with method %v", req.Method)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func main() {

	flag.Parse()

	log.Printf("fzx listening on port %v", *listenPort)

	// Wire-up handler.
	http.HandleFunc("/", Handler)

	// Listen for incoming HTTP requests.
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *listenPort), nil))
}

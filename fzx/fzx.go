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

	// DEBUG
	log.Print(req)

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
	} else {
		log.Printf("Loaded inode for %v", fzxPath)

		// Set headers using inode data.
		w.Header().Add("Content-Type", anInode.ContentType)
		w.Header().Add("Content-Length", fmt.Sprintf("%v", anInode.FileSize))
		w.Header().Add("FzxPath", anInode.FzxPath)
	}

	// DEBUG
	log.Print(anInode)

	switch req.Method {

	case "HEAD":
		log.Print("Got HEAD")
		methods.Head(w, req, anInode)
	case "OPTIONS":
		log.Print("Got OPTIONS")
		methods.Options(w, req, anInode)
	case "GET":
		log.Print("Got GET")
		methods.Get(w, req, anInode)
	case "POST":
		log.Print("Got POST")
		methods.Post(w, req, anInode)
	case "PUT":
		log.Print("Got PUT")
		methods.Put(w, req, anInode)
	case "DELETE":
		log.Print("Got DELETE")
		methods.Delete(w, req, anInode)
	case "EXECUTE":
		log.Print("Got EXECUTE")
		methods.Execute(w, req, anInode)
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
	// NOTE: This blocks anything below it from running.
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *listenPort), nil))
}

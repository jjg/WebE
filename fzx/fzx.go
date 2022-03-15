package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jjg/WebE/fzx/inode"
	"github.com/jjg/WebE/fzx/utils"
)

// TODO: Move these to flags, env, etc.
const (
	LISTEN_PORT      = ":7302"
	STORAGE_LOCATION = "./blocks"
)

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
	anInode := &inode.Inode{FzxPath: fzxPath, StorageLocation: STORAGE_LOCATION}
	if err := anInode.Load(STORAGE_LOCATION, fzxPath); err != nil {
		log.Printf("Error loading inode for %v, %v", fzxPath, err)
	}

	switch req.Method {
	case "HEAD":
		log.Print("Got HEAD")
		// TODO: Check to see if inode was actually loaded.
		// TODO: Check authorization.

		// Set headers using inode data.
		w.Header().Add("Content-Type", anInode.ContentType)
		// TODO: Determine if FileSize is really equivalent here, and if so consider renaming it.
		w.Header().Add("Content-Length", fmt.Sprintf("%v", anInode.FileSize))
		// TODO: Set additional headers?

		// Return result.
		w.WriteHeader(http.StatusOK)
	case "GET":
		log.Print("Got GET")
		// TODO: Check to see if inode was actually loaded.
		// TODO: Check authorization.

		// Set headers using inode data.
		w.Header().Add("FzxPath", anInode.FzxPath)
		// TODO: Set additional headers.

		// TODO: Locate blocks.
		// TODO: Return blocks.

		// Return result.
		w.WriteHeader(http.StatusOK)

	case "POST":
		log.Print("Got POST")

		// TODO: Check authorization.

		// Write blocks.
		log.Print("Begin processing uploaded data.")

		// TODO: Move this to flags, config, etc.
		var blockSize int64
		blockSize = 8

		// Generate a hash of the block data to use as a block filename.
		// TODO: There's got to be a better way to init this than using an empty string...
		blockDataToHash := bytes.NewBufferString("")
		if _, err := io.CopyN(blockDataToHash, req.Body, blockSize); err != nil {
			log.Fatal(err)
		}

		// DEBUG
		log.Print(blockDataToHash)

		// Open block file.
		blockF, err := os.Create("blockhash.blk")
		if err != nil {
			log.Fatal(err)
		}

		blockW := bufio.NewWriter(blockF)
		defer blockF.Close()

		log.Print("Block file created.")

		// Read one block worth of bytes out of req.Body into block file
		//if blockBytesWritten, err := io.CopyN(blockW, req.Body, blockSize); err != nil {
		if blockBytesWritten, err := io.CopyN(blockW, blockDataToHash, blockSize); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("%v bytes written to blockfile", blockBytesWritten)
		}

		// TODO: Keep reading & writing blocks until EOF.

		// Write inode.
		if err := anInode.Save(); err != nil {
			log.Print(err)
		}

		// Return result.
		if s, err := anInode.JsonString(); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, fmt.Sprintf("Error parsing inode: %v", err))
		} else {
			w.WriteHeader(http.StatusCreated)
			io.WriteString(w, s)
		}
	case "PUT":
		log.Print("Got PUT")
		// NOTE: This is probably identical to POST.
	case "DELETE":
		log.Print("Got DELETE")
		// TODO: Check authorization.
		// TODO: Delete inode.
		// TODO: Return result.
	case "EXECUTE":
		log.Print("Got EXECUTE")
		// TODO: Check authorization.
		// TODO: Handle EXECUTE request.
		// TODO: Execute specified binary.
		// TODO: Return output.
	default:
		log.Printf("I don't know what to do with method %v", req.Method)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func main() {

	log.Printf("fzx listening on port %v", LISTEN_PORT)

	// Wire-up handler.
	http.HandleFunc("/", Handler)

	// Listen for incoming HTTP requests.
	// NOTE: This blocks anything below it from running.
	log.Fatal(http.ListenAndServe(LISTEN_PORT, nil))
}

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jjg/WebE/fzx/inode"
	"github.com/jjg/WebE/fzx/utils"
)

// TODO: Move these to flags, env, etc.
const (
	LISTEN_PORT      = ":7302"
	STORAGE_LOCATION = "/perm"
)

var blockSize int64

func Handler(w http.ResponseWriter, req *http.Request) {

	// TODO: Move this to flags, env, etc.
	blockSize = 1024 * 1024 // 1MB  //8

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
	anInode := &inode.Inode{FzxPath: fzxPath, StorageLocation: STORAGE_LOCATION}
	if err := anInode.Load(STORAGE_LOCATION, fzxPath); err != nil {
		log.Printf("Error loading inode for %v, %v", fzxPath, err)
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

		// Check to see if inode was actually loaded.
		if anInode.Status == 404 {
			w.WriteHeader(http.StatusNotFound)

			// NOTE: I don't know why this is needed, but if we don't
			// write something, the connection hangs open.
			w.Write([]byte("Not found"))
			break
		}

		// TODO: Check authorization.

		// Return result.
		w.WriteHeader(http.StatusOK)

	case "GET":
		log.Print("Got GET")

		// Check to see if inode was actually loaded.
		if anInode.Status == 404 {
			w.WriteHeader(http.StatusNotFound)
			break
		}

		// TODO: Check authorization.

		// TODO: Locate blocks.

		// Return blocks.
		for _, blockName := range anInode.Blocks {
			var err error
			var blockData []byte

			// Read block from disk.
			if blockData, err = ioutil.ReadFile(fmt.Sprintf("%v/%v", STORAGE_LOCATION, blockName)); err != nil {
				panic(err)
			}

			// Write block to response.
			w.Write(blockData)
		}

	case "POST":
		log.Print("Got POST")

		// Check to see if inode was actually loaded.
		// POST should not be allowed if an inode exists,
		// so if we loaded an inode above, reject this request
		// and maybe recommend PUT instead?
		if anInode.Status != 404 {
			w.Write([]byte("Can't POST over an existing file, try PUT instead."))
			w.WriteHeader(http.StatusMethodNotAllowed)
			break
		}

		// TODO: Check authorization.

		// Set inode metadata
		v, ok := req.Header["Content-Type"]
		if ok {
			anInode.ContentType = v[0]
		}

		// Write blocks.
		log.Print("Begin processing uploaded data.")
		blockData := make([]byte, blockSize)

		// Read data from request and store it as blocks.
		totalBlockBytesWritten := 0
		for {

			// Step 1 - Get a block of data from req.Body as a byte array.
			log.Printf("Get a block of data from the request body.")
			bodyBytesRead, err := req.Body.Read(blockData)

			// DEBUG
			log.Printf("bodyBytesRead: %v", bodyBytesRead)

			// If there's no more data to read, eject.
			// TODO: See if there is a better way to detect EOF.
			if bodyBytesRead == 0 {
				break
			}

			// Step 2 - Hash the block to get the block name as a string.
			log.Printf("Hash block to get block name.")
			blockHash := utils.BytesToSha1(blockData)

			// DEBUG
			log.Printf("blockHash: %v", blockHash)

			// Step 3 - Write the block data to a file named for the block's hash.
			log.Printf("Write block data to file.")
			blockFile, err := os.Create(fmt.Sprintf("%v/%v", STORAGE_LOCATION, blockHash))
			if err != nil {
				panic(err)
			}
			defer blockFile.Close()

			blockBytesWritten, err := blockFile.Write(blockData[0:bodyBytesRead])
			totalBlockBytesWritten += blockBytesWritten

			// DEBUG
			log.Printf("blockBytesWritten: %v", blockBytesWritten)
			log.Printf("totalBlockBytesWritten: %v", totalBlockBytesWritten)

			// Step 4 - Add the block name (hash) to the inode as a string.
			log.Printf("Add block to inode.")
			anInode.Blocks = append(anInode.Blocks, blockHash)
		}

		// Update inode metadata
		anInode.FileSize = totalBlockBytesWritten

		// Write the inode.
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
		// Check to see if inode was actually loaded.
		if anInode.Status == 404 {
			w.WriteHeader(http.StatusNotFound)
			break
		}

		// TODO: Check authorization.
		// TODO: If authorized, proceed akin to POST.

	case "DELETE":
		log.Print("Got DELETE")

		// Check to see if inode was actually loaded.
		if anInode.Status == 404 {
			w.WriteHeader(http.StatusNotFound)
			break
		}

		// TODO: Check authorization.
		// TODO: Delete inode.
		// TODO: Return result.

	case "EXECUTE":
		log.Print("Got EXECUTE")

		// Check to see if inode was actually loaded.
		if anInode.Status == 404 {
			w.WriteHeader(http.StatusNotFound)
			break
		}

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

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
	STORAGE_LOCATION = "./blocks"
)

var blockSize int64

func Handler(w http.ResponseWriter, req *http.Request) {

	// TODO: Move this to flags, env, etc.
	blockSize = 8

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
		for _, blockName := range anInode.Blocks {
			var err error
			var blockData []byte

			// Read block from disk.
			if blockData, err = ioutil.ReadFile(fmt.Sprintf("%v/%v", STORAGE_LOCATION, blockName)); err != nil {
				panic(err)
			}

			// DEBUG
			log.Printf("blockData: >%v<", string(blockData[:]))

			// TODO: Write block to response.
			w.Write(blockData)
		}

		// Return result.
		w.WriteHeader(http.StatusOK)

	case "POST":
		log.Print("Got POST")

		// TODO: Check authorization.

		// Write blocks.
		log.Print("Begin processing uploaded data.")
		blockData := make([]byte, blockSize)

		// Read data from request and store it as blocks.
		for {

			// Step 1 - Get a block of data from req.Body as a byte array.
			log.Printf("Get a block of data from the request body.")
			bodyBytesRead, err := req.Body.Read(blockData)

			// DEBUG
			log.Printf("bodyBytesRead: %v", bodyBytesRead)
			log.Printf("err: %v", err)
			log.Printf("blockData: >%v<", string(blockData[:]))

			// If there's no more data to read, eject.
			// TODO: See if there is a better way to detect EOF.
			if err != nil {
				if bodyBytesRead == 0 {
					break
				} else {
					panic(err)
				}
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

			// DEBUG
			log.Printf("err: %v", err)

			blockBytesWritten, err := blockFile.Write(blockData[0:bodyBytesRead])

			// DEBUG
			log.Printf("blockBytesWritten: %v", blockBytesWritten)
			log.Printf("err: %v", err)

			// Step 4 - Add the block name (hash) to the inode as a string.
			log.Printf("Add block to inode.")
			anInode.Blocks = append(anInode.Blocks, blockHash)
		}

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

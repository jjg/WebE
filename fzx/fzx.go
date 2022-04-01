package main

import (
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

		// Step 1 - Get a block of data from req.body as a byte array
		log.Printf("Get a block of data from the request body")
		blockData := make([]byte, blockSize)
		bodyBytesRead, err := req.Body.Read(blockData)

		// DEBUG
		log.Printf("bodyBytesRead: %v", bodyBytesRead)
		log.Printf("err: %v", err)
		log.Printf("blockData: >%v<", string(blockData[:]))

		// Step 2 - Hash the block to get the block name as a string
		log.Printf("Hash block to get block name")
		blockHash := utils.BytesToSha1(blockData)

		// DEBUG
		log.Printf("blockHash: %v", blockHash)

		// Step 3 - Write the block data to a file named for the block's hash
		log.Printf("Write block data to file")
		blockFile, err := os.Create(fmt.Sprintf("%v/%v", STORAGE_LOCATION, blockHash))
		defer blockFile.Close()

		// DEBUG
		log.Printf("err: %v", err)

		blockBytesWritten, err := blockFile.Write(blockData)

		// DEBUG
		log.Printf("blockBytesWritten: %v", blockBytesWritten)
		log.Printf("err: %v", err)

		// Step 4 - Add the block name (hash) to the inode as a string
		log.Printf("Add block to inode")
		anInode.Blocks = append(anInode.Blocks, blockHash)

		/*
				// Generate a hash of the block data to use as a block filename.
				// TODO: There's got to be a better way to init this than using an empty string...
				// NOTE: blockData is only used to initialize blockDataBuffer afaik.
				var blockData []byte //:= bytes.NewBufferString("")
				blockDataBuffer := bytes.NewBuffer(blockData)

				// TODO: Make sure we properly handle the last/partial block.
				//if _, err := io.CopyN(blockData, req.Body, blockSize); err != nil {
				if _, err := io.CopyN(blockDataBuffer, req.Body, blockSize); err != nil {
					log.Fatal(err)
				}

				// DEBUG
				log.Print("blockDataBuffer: (contents of block) ", blockDataBuffer)

				// Generate a hash of the incoming block data.
				blockHash := utils.BytesToSha1(blockDataBuffer)

				// DEBUG
				log.Print("blockHash (block filename): ", blockHash)

					// Open block file.
					blockF, err := os.Create(blockHash)
					if err != nil {
						log.Fatal(err)
					}

						blockW := bufio.NewWriter(blockF)
						defer blockF.Close()

						// Write block data into block file.
						// TODO: Since we carve-up the data into blocks above,
						// this can probably use io.Copy() instead.
						if blockBytesWritten, err := io.CopyN(blockW, blockData, blockSize); err != nil {
							log.Fatal(err)
						} else {
							log.Printf("%v bytes written to blockfile %v", blockBytesWritten, blockHash)
						}

						log.Printf("Block file %v created.", blockHash)

			// Add block hash to inode block array.
			anInode.Blocks = append(anInode.Blocks, blockHash)

			// TODO: Keep reading & writing blocks until EOF.
		*/

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

package methods

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jjg/WebE/fzx/inode"
	"github.com/jjg/WebE/fzx/utils"
)

// TODO: Since writing is shared between POST and PUT, consider abstracting
// the plubming here into a more generic function shared by the POST and PUT handlers.
func Post(w http.ResponseWriter, req *http.Request, anInode *inode.Inode) {

	// Write blocks.
	log.Print("Begin processing uploaded data.")
	blockData := make([]byte, anInode.BlockSize)

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
		blockFile, err := os.Create(fmt.Sprintf("%v/%v", *anInode.StorageLocation, blockHash))
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

	// DEBUG
	log.Print(anInode)

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
}

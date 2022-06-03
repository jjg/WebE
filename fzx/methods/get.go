package methods

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jjg/WebE/fzx/inode"
)

func Get(w http.ResponseWriter, req *http.Request, anInode *inode.Inode) {

	// Return blocks.
	for _, blockName := range anInode.Blocks {
		var err error
		var blockData []byte

		// Read block from disk.
		if blockData, err = ioutil.ReadFile(fmt.Sprintf("%v/%v", *anInode.StorageLocation, blockName)); err != nil {
			panic(err)
		}

		// Write block to response.
		w.Write(blockData)
	}
}

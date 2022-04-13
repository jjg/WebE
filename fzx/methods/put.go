package methods

import (
	"net/http"

	"github.com/jjg/WebE/fzx/inode"
)

func Put(w http.ResponseWriter, req *http.Request, anInode *inode.Inode) {

	// TODO: Check authorization.
	// TODO: Update the file.
}

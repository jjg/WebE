package methods

import (
	"net/http"

	"github.com/jjg/WebE/fzx/inode"
)

func Execute(w http.ResponseWriter, req *http.Request, anInode *inode.Inode) {
	// Check to see if inode was actually loaded.
	if anInode.Status == 404 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: Check authorization.
	// TODO: Execute the program.
}

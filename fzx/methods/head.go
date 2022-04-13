package methods

import (
	"net/http"

	"github.com/jjg/WebE/fzx/inode"
)

func Head(w http.ResponseWriter, req *http.Request, anInode *inode.Inode) {

	// Check to see if inode was actually loaded.
	if anInode.Status == 404 {
		w.WriteHeader(http.StatusNotFound)

		// NOTE: I don't know why this is needed, but if we don't
		// write something, the connection hangs open.
		w.Write([]byte("Not found"))
		return
	}

	// TODO: Check authorization.

	// Return result.
	w.WriteHeader(http.StatusOK)
}

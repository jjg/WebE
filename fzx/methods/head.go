package methods

import (
	"net/http"

	"github.com/jjg/WebE/fzx/inode"
)

// TODO: Consider eliminating this module as it does very little right now.
func Head(w http.ResponseWriter, req *http.Request, anInode *inode.Inode) {

	// Return result.
	w.WriteHeader(http.StatusOK)
}

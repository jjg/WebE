package inode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInode(t *testing.T) {

	fzxPath := ".com.example/foo/bar.txt"
	storageLocation := "./blocks"

	// Initialize new inode
	inode := &Inode{FzxPath: fzxPath, StorageLocation: storageLocation}

	// TODO: Set some values.

	// Save inode
	if err := inode.Save(); err != nil {
		t.Fatal(err)
	}

	// Load inode
	loadedInode := &Inode{}
	if err := loadedInode.Load(storageLocation, fzxPath); err != nil {
		t.Fatal(err)
	}

	// Compare loaded to original values.
	assert.Equal(t, inode.FzxPath, loadedInode.FzxPath)
}

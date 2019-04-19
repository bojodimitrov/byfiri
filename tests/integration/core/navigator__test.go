package coreint__test

import (
	"testing"

	"github.com/bojodimitrov/byfiri/core"
)

func TestEnterDirectory(t *testing.T) {
	storage, currentDir := setupFileSystem()

	dirInode := core.AllocateDirectory(storage, currentDir, "dir")

	currentDir, _ = core.EnterDirectory(storage, currentDir, "dir")

	if currentDir.DirectoryInode != dirInode {
		t.Errorf("got %d, want %d", currentDir.DirectoryInode, dirInode)
	}
}

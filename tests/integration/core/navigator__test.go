package coreint__test

import (
	"testing"

	"github.com/bojodimitrov/byfiri/core"
)

func TestEnterDirectory(t *testing.T) {
	storage, currentDir := setupFileSystem(t)

	dirInode, err := core.AllocateDirectory(storage, currentDir, "dir")
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}

	currentDir, err = core.EnterDirectory(storage, currentDir, "dir")
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}

	if currentDir.DirectoryInode != dirInode {
		t.Errorf("got %d, want %d", currentDir.DirectoryInode, dirInode)
	}
}

func TestEnterDirectoryPath(t *testing.T) {
	storage, currentDir := setupFileSystem(t)

	core.AllocateDirectory(storage, currentDir, "dir")

	currentDir, _ = core.EnterDirectory(storage, currentDir, "dir")

	dirInode, err := core.AllocateDirectory(storage, currentDir, "dir")
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}

	currentDir, err = core.EnterDirectory(storage, currentDir, "dir")
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}

	if currentDir.DirectoryInode != dirInode {
		t.Errorf("got %d, want %d", currentDir.DirectoryInode, dirInode)
	}
}

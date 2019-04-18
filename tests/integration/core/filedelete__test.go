package coreint__test

import (
	"testing"

	"github.com/bojodimitrov/byfiri/core"
)

func TestFileDelete(t *testing.T) {
	storage, currentDir := setupFileSystem()

	content := "file content"

	fileInode1 := core.AllocateFile(storage, currentDir, "test file", content)
	fileInode2 := core.AllocateFile(storage, currentDir, "test file 2", content)

	if fileInode1 != 2 || fileInode2 != 3 {
		t.Errorf("got %d, want %d", []int{fileInode1, fileInode2}, []int{2, 3})
	}
	core.DeleteFile(storage, fileInode1)

	fileInode3 := core.AllocateFile(storage, currentDir, "test file 3", content)

	if fileInode3 != 2 || fileInode2 != 3 {
		t.Errorf("got %d, want %d", []int{fileInode3, fileInode2}, []int{2, 3})
	}
}

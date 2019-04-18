package coreint__test

import (
	"testing"

	"github.com/bojodimitrov/gofys/core"
)

func TestFileDelete(t *testing.T) {
	storage := setupFileSystem()

	content := "file content"

	fileInode1 := core.AllocateFile(storage, 1, content)
	fileInode2 := core.AllocateFile(storage, 1, content)

	if fileInode1 != 2 && fileInode2 != 3 {
		t.Errorf("got %q, want %q", []int{fileInode1, fileInode2}, []int{2, 3})
	}
	core.DeleteFile(storage, fileInode1)

	fileInode3 := core.AllocateFile(storage, 1, content)

	if fileInode3 != 2 && fileInode2 != 3 {
		t.Errorf("got %q, want %q", []int{fileInode3, fileInode2}, []int{2, 3})
	}
}

package coreint__test

import (
	"testing"

	"github.com/bojodimitrov/byfiri/core"
	"github.com/bojodimitrov/byfiri/structures"
)

func TestFileDelete(t *testing.T) {
	storage, currentDir := setupFileSystem()

	content := "file content"

	fileInode1 := core.AllocateFile(storage, currentDir, "test file", content)
	fileInode2 := core.AllocateFile(storage, currentDir, "test file 2", content)

	if fileInode1 != 2 || fileInode2 != 3 {
		t.Errorf("got %d, want %d", []int{fileInode1, fileInode2}, []int{2, 3})
	}
	core.DeleteFile(storage, currentDir, fileInode1)

	fileInode3 := core.AllocateFile(storage, currentDir, "test file 3", content)

	if fileInode3 != 2 || fileInode2 != 3 {
		t.Errorf("got %d, want %d", []int{fileInode3, fileInode2}, []int{2, 3})
	}
}

func TestFileDeleteInDirectory(t *testing.T) {
	storage, root := setupFileSystem()

	content := "file content"

	fileInode := core.AllocateFile(storage, root, "test file", content)
	resultWithFile := core.ReadDirectory(storage, root.DirectoryInode)
	addedFile := []structures.DirectoryEntry{
		structures.DirectoryEntry{FileName: ".", Inode: uint32(root.DirectoryInode)},
		structures.DirectoryEntry{FileName: "..", Inode: 0},
		structures.DirectoryEntry{FileName: "test file", Inode: uint32(fileInode)},
	}
	for i, entry := range resultWithFile {
		if entry.FileName != addedFile[i].FileName || entry.Inode != addedFile[i].Inode {
			t.Errorf("got %q, want %q", resultWithFile, addedFile)
		}
	}
	core.DeleteFile(storage, root, fileInode)
	resultNoFile := core.ReadDirectory(storage, root.DirectoryInode)
	noFile := []structures.DirectoryEntry{
		structures.DirectoryEntry{FileName: ".", Inode: uint32(root.DirectoryInode)},
		structures.DirectoryEntry{FileName: "..", Inode: 0},
	}
	for i, entry := range resultNoFile {
		if entry.FileName != noFile[i].FileName || entry.Inode != noFile[i].Inode {
			t.Errorf("got %q, want %q", resultNoFile, noFile)
		}
	}
}

func TestDeleteDirectory(t *testing.T) {
	storage, dir := setupFileSystem()

	dirInode := core.AllocateDirectory(storage, dir, "root lv1 dir1")

	dir, _ = core.EnterDirectory(storage, dir, "root lv1 dir1")
	core.AllocateDirectory(storage, dir, "root lv2 dir1")
	dir, _ = core.EnterDirectory(storage, dir, "root lv2 dir1")
	core.AllocateFile(storage, dir, "root lv2 f1", "hello there")
	dir, _ = core.EnterDirectory(storage, dir, "..")
	dir, _ = core.EnterDirectory(storage, dir, "..")

	core.DeleteDirectory(storage, dir, dirInode)
	resultNonExistant := core.ReadDirectory(storage, dirInode)

	if resultNonExistant != nil {
		t.Errorf("got %q, want %q", resultNonExistant, "nil")
	}

	fileInode1 := core.AllocateFile(storage, dir, "test file", "test content")
	fileInode2 := core.AllocateFile(storage, dir, "test file 2", "test content")
	if fileInode1 != 2 || fileInode2 != 3 {
		t.Errorf("got %d, want %d", []int{fileInode1, fileInode2}, []int{2, 3})
	}
}

package coreint__test

import (
	"testing"

	"github.com/bojodimitrov/byfiri/core"
	"github.com/bojodimitrov/byfiri/structures"
)

func setupFileSystem(t *testing.T) ([]byte, *structures.DirectoryIterator) {
	size := 512 * 1048576
	storage := core.InitFsSpace(size)
	root, err := core.AllocateAllStructures(storage, size, structures.DefaultBlockSize)
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}
	return storage, root
}

func TestRoot(t *testing.T) {
	storage, _ := setupFileSystem(t)
	rootContent, err := core.ReadDirectory(storage, 1)
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}
	rootEntry := rootContent[0]
	expected := structures.DirectoryEntry{FileName: ".", Inode: 1}

	if rootEntry.FileName != "." && rootEntry.Inode != 1 {
		t.Errorf("got %q, want %q", rootContent, expected)
	}
}

func TestDirectoryAllocation(t *testing.T) {
	storage, root := setupFileSystem(t)
	dirInode, err := core.AllocateDirectory(storage, root, "test dir")
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}
	result, err := core.ReadDirectory(storage, dirInode)
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}
	expected := []structures.DirectoryEntry{
		structures.DirectoryEntry{FileName: ".", Inode: uint32(dirInode)},
		structures.DirectoryEntry{FileName: "..", Inode: uint32(root.DirectoryInode)},
	}
	for i, entry := range result {
		if entry.FileName != expected[i].FileName || entry.Inode != expected[i].Inode {
			t.Errorf("got %q, want %q", result, expected)
		}
	}
}

func TestFilesInDirectory(t *testing.T) {
	storage, root := setupFileSystem(t)
	dirInode, err := core.AllocateDirectory(storage, root, "test dir")
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}
	fileInode, err := core.AllocateFile(storage, root, "test file", "file content")
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}
	result, err := core.ReadDirectory(storage, root.DirectoryInode)
	if err != nil {
		t.Errorf("got %q, want %q", err, "nil")
	}
	expected := []structures.DirectoryEntry{
		structures.DirectoryEntry{FileName: ".", Inode: uint32(root.DirectoryInode)},
		structures.DirectoryEntry{FileName: "..", Inode: 0},
		structures.DirectoryEntry{FileName: "test dir", Inode: uint32(dirInode)},
		structures.DirectoryEntry{FileName: "test file", Inode: uint32(fileInode)},
	}
	for i, entry := range result {
		if entry.FileName != expected[i].FileName || entry.Inode != expected[i].Inode {
			t.Errorf("got %q, want %q", result, expected)
		}
	}
}

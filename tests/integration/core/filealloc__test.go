package coreint__test

import (
	"strings"
	"testing"

	"github.com/bojodimitrov/gofys/core"
	"github.com/bojodimitrov/gofys/structures"
)

func setupFileSystem() ([]byte, *structures.DirectoryIterator) {
	size := 512 * 1048576
	storage := core.InitFsSpace(size)
	currentDir := core.AllocateAllStructures(storage, size, structures.DefaultBlockSize)
	return storage, currentDir
}

func TestRoot(t *testing.T) {
	storage, _ := setupFileSystem()
	rootContent := core.ReadDirectory(storage, 1)[0]
	expected := structures.DirectoryEntry{FileName: ".", Inode: 1}

	if rootContent.FileName != "." && rootContent.Inode != 1 {
		t.Errorf("got %q, want %q", rootContent, expected)
	}
}

func TestFileAllocation(t *testing.T) {
	storage, currentDir := setupFileSystem()
	content := "file content"
	fileInode := core.AllocateFile(storage, currentDir, "test file", content)
	result := core.ReadFile(storage, fileInode)

	if content != result {
		t.Errorf("got %q, want %q", result, content)
	}
}

func TestFileAllocationLarge(t *testing.T) {
	storage, currentDir := setupFileSystem()
	content := strings.Repeat("a", structures.DefaultBlockSize+10)
	fileInode := core.AllocateFile(storage, currentDir, "test file", content)
	result := core.ReadFile(storage, fileInode)

	if content != result {
		t.Errorf("got %q, want %q", result, content)
	}
}

func TestFileAllocationConsequtiveInodes(t *testing.T) {
	storage, currentDir := setupFileSystem()
	content := "file content"
	rootContent := core.ReadDirectory(storage, 1)[0]
	fileInode1 := core.AllocateFile(storage, currentDir, "test file", content)
	fileInode2 := core.AllocateFile(storage, currentDir, "test file 2", content)
	fileInode3 := core.AllocateFile(storage, currentDir, "test file 3", content)

	if rootContent.Inode != 1 || fileInode1 != 2 || fileInode2 != 3 || fileInode3 != 4 {
		t.Errorf("got %d, want %d", []int{int(rootContent.Inode), fileInode1, fileInode2, fileInode3}, []int{1, 2, 3, 4})
	}
}

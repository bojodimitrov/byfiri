package coreint__test

import (
	"strings"
	"testing"

	"github.com/bojodimitrov/gofys/core"
	"github.com/bojodimitrov/gofys/structures"
)

func TestRoot(t *testing.T) {
	size := 512 * 1048576
	storage := core.InitFsSpace(size)
	core.AllocateAllStructures(storage, size, structures.DefaultBlockSize)
	rootContent := core.ReadDirectory(storage, 1)[0]
	expected := structures.DirectoryContent{FileName: ".", Inode: 1}

	if rootContent.FileName != "." && rootContent.Inode != 1 {
		t.Errorf("got %q, want %q", rootContent, expected)
	}
}

func TestFileAllocation(t *testing.T) {
	size := 512 * 1048576
	storage := core.InitFsSpace(size)
	core.AllocateAllStructures(storage, size, structures.DefaultBlockSize)
	content := "file content"
	fileInode := core.AllocateFile(storage, 1, content)
	result := core.ReadFile(storage, fileInode)

	if content != result {
		t.Errorf("got %q, want %q", content, result)
	}
}

func TestFileAllocationLarge(t *testing.T) {
	size := 512 * 1048576
	storage := core.InitFsSpace(size)
	core.AllocateAllStructures(storage, size, structures.DefaultBlockSize)
	content := strings.Repeat("a", structures.DefaultBlockSize+10)
	fileInode := core.AllocateFile(storage, 1, content)
	result := core.ReadFile(storage, fileInode)

	if content != result {
		t.Errorf("got %q, want %q", content, result)
	}
}

func TestFileAllocationConsequtiveInodes(t *testing.T) {
	size := 512 * 1048576
	storage := core.InitFsSpace(size)
	core.AllocateAllStructures(storage, size, structures.DefaultBlockSize)

	content := "file content"
	rootContent := core.ReadDirectory(storage, 1)[0]
	fileInode1 := core.AllocateFile(storage, 1, content)
	fileInode2 := core.AllocateFile(storage, 1, content)
	fileInode3 := core.AllocateFile(storage, 1, content)

	if rootContent.Inode != 1 && fileInode1 != 2 && fileInode2 != 3 && fileInode3 != 4 {
		t.Errorf("got %q, want %q", []int{int(rootContent.Inode), fileInode1, fileInode2, fileInode3}, []int{1, 2, 3, 4})
	}
}

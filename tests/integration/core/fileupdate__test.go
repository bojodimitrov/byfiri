package coreint__test

import (
	"strings"
	"testing"

	"github.com/bojodimitrov/gofys/core"
)

func TestFileUpdate(t *testing.T) {
	storage := setupFileSystem()

	content := "file content"
	updatedContent := "update file content"
	fileInode := core.AllocateFile(storage, 1, content)
	resultPure := core.ReadFile(storage, fileInode)
	core.UpdateFile(storage, fileInode, updatedContent)
	resultTouched := core.ReadFile(storage, fileInode)

	if resultPure == resultTouched || resultTouched != updatedContent {
		t.Errorf("got %q, want %q", resultTouched, updatedContent)
	}
}

func TestFileEnlarge(t *testing.T) {
	storage := setupFileSystem()

	content := "file content"
	enlargeContent := strings.Repeat("x", 20000)

	fileInode1 := core.AllocateFile(storage, 1, content)
	fileInode2 := core.AllocateFile(storage, 1, content)
	core.UpdateFile(storage, fileInode1, enlargeContent)

	resultFile1 := core.ReadFile(storage, fileInode1)
	resultFile2 := core.ReadFile(storage, fileInode2)

	if resultFile1 != enlargeContent || resultFile2 != content {
		t.Errorf("file1 content len: got %d, want %d", len(resultFile1), len(enlargeContent))
		t.Errorf("file2: got %q, want %q", resultFile2, content)
	}
}

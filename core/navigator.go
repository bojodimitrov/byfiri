package core

import (
	"fmt"

	"github.com/bojodimitrov/byfiri/structures"
)

// EnterDirectory returns DirectoryIterator of the desired directory within current directory
func EnterDirectory(storage []byte, current *structures.DirectoryIterator, directory string) (*structures.DirectoryIterator, error) {
	currentDirectoryContent := ReadDirectory(storage, current.DirectoryInode)
	inode := 0
	for _, dirEntry := range currentDirectoryContent {
		if dirEntry.FileName == directory {
			inode = int(dirEntry.Inode)
		}
	}
	if inode == 0 {
		return nil, fmt.Errorf("enter directory: name not found")
	}
	fsdata := ReadMetadata(storage)
	inodeInfo := ReadInode(storage, fsdata, inode)
	if inodeInfo.Mode != 0 {
		return nil, fmt.Errorf("enter directory: file is not directory")
	}
	wantedDirectoryContent := ReadDirectory(storage, inode)
	wantedDirectoryIt := structures.DirectoryIterator{DirectoryInode: inode, DirectoryContent: wantedDirectoryContent}
	return &wantedDirectoryIt, nil
}

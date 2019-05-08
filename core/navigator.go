package core

import (
	"fmt"

	"github.com/bojodimitrov/byfiri/logger"
	"github.com/bojodimitrov/byfiri/structures"
)

// EnterDirectory returns DirectoryIterator of the desired directory within current directory
func EnterDirectory(storage []byte, current *structures.DirectoryIterator, directory string) (*structures.DirectoryIterator, error) {
	currentDirectoryContent, err := ReadDirectory(storage, current.DirectoryInode)
	if err != nil {
		return nil, err
	}
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
	inodeInfo, err := ReadInode(storage, fsdata, inode)
	if err != nil {
		// Log err
		logger.Log("enter directory: " + err.Error())
		return nil, fmt.Errorf("enter directory: could not read inode")
	}
	if inodeInfo.Mode != 0 {
		return nil, fmt.Errorf("enter directory: file is not directory")
	}
	wantedDirectoryContent, err := ReadDirectory(storage, inode)
	if err != nil {
		return nil, err
	}
	wantedDirectoryIt := structures.DirectoryIterator{DirectoryInode: inode, DirectoryContent: wantedDirectoryContent}
	return &wantedDirectoryIt, nil
}

// IsDirectory checks if name is directory
func IsDirectory(storage []byte, current *structures.DirectoryIterator, name string) (bool, error) {

	inode := GetInode(storage, current, name)
	fsdata := ReadMetadata(storage)
	inodeInfo, err := ReadInode(storage, fsdata, inode)
	if err != nil {
		logger.Log("enter directory: " + err.Error())
		return false, fmt.Errorf("is directory: could not read inode")
	}
	return inodeInfo.Mode == 0, nil
}

// GetInode returns inode behind a name
func GetInode(storage []byte, current *structures.DirectoryIterator, name string) int {
	inode := 0
	for _, entry := range current.DirectoryContent {
		if entry.FileName == name {
			inode = int(entry.Inode)
		}
	}
	return inode
}

//GetCurrentName returns name of current directory
func GetCurrentName(storage []byte, current *structures.DirectoryIterator) string {
	currentDirName := ""
	if current.DirectoryInode == 1 {
		currentDirName = "root"
	} else {
		parent, _ := EnterDirectory(storage, current, "..")
		for _, value := range parent.DirectoryContent {
			if value.Inode == uint32(current.DirectoryInode) {
				currentDirName = value.FileName
			}
		}
	}
	return currentDirName
}

//GetPath returns path to current directory
func GetPath(storage []byte, current *structures.DirectoryIterator) string {
	paths := []string{GetCurrentName(storage, current)}
	for current.DirectoryInode != 1 {
		current, _ = EnterDirectory(storage, current, "..")
		paths = append(paths, GetCurrentName(storage, current))
	}
	path := ""
	for i := len(paths) - 1; i >= 0; i-- {
		path += paths[i] + "\\"
	}
	path = path[:len(path)-1]
	path += ">"
	return path
}

package core

import (
	"fmt"
	"strings"

	"github.com/bojodimitrov/byfiri/structures"
)

// EnterDirectory returns DirectoryIterator of the desired directory within current directory
func EnterDirectory(storage []byte, current *structures.DirectoryIterator, path string) (*structures.DirectoryIterator, error) {
	fields := strings.Split(path, "\\")
	var err error
	for _, directory := range fields {
		current, err = navigate(storage, current, directory)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

func navigate(storage []byte, current *structures.DirectoryIterator, directory string) (*structures.DirectoryIterator, error) {
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

// IsDirectory checks if name is directory
func IsDirectory(storage []byte, current *structures.DirectoryIterator, name string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("read directory: inode does not exist")
		}
	}()

	inode := GetInode(storage, current, name)
	fsdata := ReadMetadata(storage)
	inodeInfo := ReadInode(storage, fsdata, inode)
	return inodeInfo.Mode == 0
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

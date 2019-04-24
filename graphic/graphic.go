package graphic

import (
	"fmt"
	"strings"

	"github.com/bojodimitrov/byfiri/core"
	"github.com/bojodimitrov/byfiri/structures"
)

func iterateDirectoryRecursively(storage []byte, fsdata *structures.Metadata, currentDirectory *structures.DirectoryIterator, level int) error {
	for _, entry := range currentDirectory.DirectoryContent {
		if entry.FileName != "." && entry.FileName != ".." {
			inodeInfo, err := core.ReadInode(storage, fsdata, int(entry.Inode))
			if err != nil {
				return err
			}
			fmt.Print(strings.Repeat("|   ", level))
			fmt.Print("|___")
			if inodeInfo.Mode == 0 {
				child, err := core.EnterDirectory(storage, currentDirectory, entry.FileName)
				if err != nil {
					return err
				}
				fmt.Println(entry.FileName + ":")
				err = iterateDirectoryRecursively(storage, fsdata, child, level+1)
				if err != nil {
					return err
				}
			} else {
				fmt.Println(entry.FileName)
			}
		}
	}
	return nil
}

//DisplayDirectoryTree prints directory tree in readable way
func DisplayDirectoryTree(storage []byte, currentDirectory *structures.DirectoryIterator) error {
	fsdata := core.ReadMetadata(storage)

	if currentDirectory.DirectoryInode == 1 {
		fmt.Println("root:")
	} else {
		currentDirName := ""
		parent, _ := core.EnterDirectory(storage, currentDirectory, "..")
		for _, value := range parent.DirectoryContent {
			if value.Inode == uint32(currentDirectory.DirectoryInode) {
				currentDirName = value.FileName
			}
		}
		fmt.Println(currentDirName + ":")
	}
	err := iterateDirectoryRecursively(storage, fsdata, currentDirectory, 0)
	if err != nil {
		return err
	}
	return nil
}

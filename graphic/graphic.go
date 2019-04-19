package graphic

import (
	"fmt"
	"strings"

	"github.com/bojodimitrov/byfiri/core"
	"github.com/bojodimitrov/byfiri/structures"
)

func iterateDirectoryRecursively(storage []byte, fsdata *structures.Metadata, currentDirectory *structures.DirectoryIterator, level int) {
	for _, entry := range currentDirectory.DirectoryContent {
		if entry.FileName != "." && entry.FileName != ".." {
			inodeInfo := core.ReadInode(storage, fsdata, int(entry.Inode))
			fmt.Print(strings.Repeat("|   ", level))
			fmt.Print("|___")
			if inodeInfo.Mode == 0 {
				child, err := core.EnterDirectory(storage, currentDirectory, entry.FileName)
				if err != nil {
					panic("delete directory: could not recursively delete content")
				}
				fmt.Println(entry.FileName + ":")
				iterateDirectoryRecursively(storage, fsdata, child, level+1)
			} else {
				fmt.Println(entry.FileName)
			}
		}
	}
}

//DisplayDirectoryTree prints directory tree in readable way
func DisplayDirectoryTree(storage []byte, currentDirectory *structures.DirectoryIterator) {
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
		fmt.Println(currentDirName)
	}
	iterateDirectoryRecursively(storage, fsdata, currentDirectory, 0)
}

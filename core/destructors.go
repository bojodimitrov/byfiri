package core

import (
	"fmt"

	"github.com/bojodimitrov/byfiri/structures"
)

func deleteInode(storage []byte, fsdata *structures.Metadata, inode int) {
	markOnBitmap(storage, fsdata, false, inode, structures.Inodes)
}

func deleteBlocks(storage []byte, fsdata *structures.Metadata, inodeInfo *structures.Inode) {
	for _, block := range inodeInfo.BlocksLocations {
		markOnBitmap(storage, fsdata, false, int(block), structures.Blocks)
	}
}

func removeFileFromDirectory(storage []byte, current *structures.DirectoryIterator, fileInode int) {
	content := ReadDirectory(storage, current.DirectoryInode)

	index := 0
	for i, entry := range content {
		if entry.Inode == uint32(fileInode) {
			index = i
		}
	}

	content[index] = content[len(content)-1]
	content = content[:len(content)-1]

	UpdateDirectory(storage, current.DirectoryInode, content)

	current.DirectoryContent[index] = current.DirectoryContent[len(current.DirectoryContent)-1]
	current.DirectoryContent = current.DirectoryContent[:len(current.DirectoryContent)-1]
}

//DeleteFile deletes file
func DeleteFile(storage []byte, currentDirectory *structures.DirectoryIterator, inode int) {
	if inode == 0 {
		fmt.Println("delete file: inode cannot be 0")
		return
	}
	if inode == 1 {
		fmt.Println("delete file: cannot delete root")
		return
	}

	fsdata := ReadMetadata(storage)
	inodeInfo := ReadInode(storage, fsdata, inode)
	clearFile(storage, inodeInfo.BlocksLocations, fsdata)
	clearInode(storage, fsdata, inode)
	removeFileFromDirectory(storage, currentDirectory, inode)
	deleteBlocks(storage, fsdata, inodeInfo)
	deleteInode(storage, fsdata, inode)
}

func directoryContainsDirectory(storage []byte, fsdata *structures.Metadata, currentDirectory *structures.DirectoryIterator) bool {
	for _, entry := range currentDirectory.DirectoryContent {
		if entry.FileName != "." && entry.FileName != ".." {
			inodeInfo := ReadInode(storage, fsdata, int(entry.Inode))
			if inodeInfo.Mode == 0 {
				return true
			}
		}
	}
	return false
}

func iterateDirectoryRecursively(storage []byte, fsdata *structures.Metadata, currentDirectory *structures.DirectoryIterator) {
	//bottom of recursion
	if !directoryContainsDirectory(storage, fsdata, currentDirectory) {
		// we delete all other files
		for _, entry := range currentDirectory.DirectoryContent {
			if entry.FileName != "." && entry.FileName != ".." {
				DeleteFile(storage, currentDirectory, int(entry.Inode))
			}
		}
		return
	}
	for _, entry := range currentDirectory.DirectoryContent {
		if entry.FileName != "." && entry.FileName != ".." {
			inodeInfo := ReadInode(storage, fsdata, int(entry.Inode))
			if inodeInfo.Mode == 0 {
				child, err := EnterDirectory(storage, currentDirectory, entry.FileName)
				if err != nil {
					panic("delete directory: could not recursively delete content")
				}
				iterateDirectoryRecursively(storage, fsdata, child)
				DeleteFile(storage, currentDirectory, int(entry.Inode))
			}
			if inodeInfo.Mode == 1 {
				DeleteFile(storage, currentDirectory, int(entry.Inode))
			}
		}
	}
}

//DeleteDirectory deletes directory recursively
func DeleteDirectory(storage []byte, currentDirectory *structures.DirectoryIterator, inode int) {
	fsdata := ReadMetadata(storage)
	fileName := ""
	for _, entry := range currentDirectory.DirectoryContent {
		if entry.Inode == uint32(inode) {
			fileName = entry.FileName
		}
	}
	deletedDirectory, _ := EnterDirectory(storage, currentDirectory, fileName)
	iterateDirectoryRecursively(storage, fsdata, deletedDirectory)
	DeleteFile(storage, currentDirectory, inode)
}

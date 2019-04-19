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

package core

import (
	"github.com/bojodimitrov/gofys/structures"
)

func deleteInode(storage []byte, fsdata *structures.Metadata, inode int) {
	markOnBitmap(storage, fsdata, false, inode, structures.Inodes)
}

func deleteBlocks(storage []byte, fsdata *structures.Metadata, inodeInfo *structures.Inode) {
	for _, block := range inodeInfo.BlocksLocations {
		markOnBitmap(storage, fsdata, false, int(block), structures.Blocks)
	}
}

//DeleteFile deletes file
func DeleteFile(storage []byte, inode int) {
	fsdata := ReadMetadata(storage)
	inodeInfo := ReadInode(storage, fsdata, inode)
	clearFile(storage, inodeInfo.BlocksLocations, &fsdata)
	clearInode(storage, &fsdata, inode)
	deleteBlocks(storage, &fsdata, &inodeInfo)
	deleteInode(storage, &fsdata, inode)
}

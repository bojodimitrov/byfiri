package core

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/bojodimitrov/byfiri/util"

	"github.com/bojodimitrov/byfiri/diracts"
	"github.com/bojodimitrov/byfiri/errors"
	"github.com/bojodimitrov/byfiri/structures"
)

//ReadInode returns an Inode structure written behind inode location
func ReadInode(storage []byte, metadata *structures.Metadata, inode int) *structures.Inode {
	inodesBeginning := int(metadata.Root)
	inodeLocation := inodesBeginning + inode*int(structures.InodeSize)
	//First 3 bytes are mode
	modeStr := Read(storage, inodeLocation, 3)
	mode, err := strconv.Atoi(modeStr)
	errors.CorruptInode(err, inode)
	inodeLocation += 3
	//Next 10 bytes are size
	sizeStr := Read(storage, inodeLocation, 10)
	size, err := strconv.Atoi(sizeStr)
	errors.CorruptInode(err, inode)
	inodeLocation += 10

	var blocksGathered [12]uint32
	//Next 12 * 10 bytes are blocks of file
	for i := 0; i < 12; i++ {
		block := Read(storage, inodeLocation, 10)
		value, err := strconv.Atoi(block)
		errors.CorruptInode(err, inode)
		blocksGathered[i] = uint32(value)
		inodeLocation += 10
	}
	inodeInfo := structures.Inode{Mode: uint8(mode), Size: uint32(size), BlocksLocations: blocksGathered}
	return &inodeInfo
}

//ReadContent reads file content
func ReadContent(storage []byte, metadata *structures.Metadata, inodeInfo *structures.Inode) string {
	blocksBeginning := int(metadata.FirstBlock)
	blockSize := int(metadata.BlockSize)
	fileSize := int(inodeInfo.Size)
	var contentBuffer bytes.Buffer

	for i := 0; i < 12; i++ {
		if inodeInfo.BlocksLocations[i] != 0 {
			content := ReadRaw(storage, blocksBeginning+int(inodeInfo.BlocksLocations[i])*blockSize, util.Min(blockSize, fileSize))
			contentBuffer.WriteString(content)
			fileSize -= blockSize
		}
	}
	content := contentBuffer.String()

	return content
}

//ReadFile returns file content
func ReadFile(storage []byte, inode int) string {
	if inode == 0 {
		fmt.Println("read file: inode cannot be 0")
		return ""
	}

	fsdata := ReadMetadata(storage)
	inodeInfo := ReadInode(storage, fsdata, inode)
	if inodeInfo.Mode == 0 {
		fmt.Println("update file: file is directory")
		return ""
	}
	return ReadContent(storage, fsdata, inodeInfo)
}

//ReadDirectory returns directory content
func ReadDirectory(storage []byte, inode int) []structures.DirectoryEntry {
	if inode == 0 {
		fmt.Println("read directory: inode cannot be 0")
		return nil
	}

	fsdata := ReadMetadata(storage)
	inodeInfo := ReadInode(storage, fsdata, inode)
	if inodeInfo.Mode == 1 {
		fmt.Println("update directory: directory is file")
		return nil
	}
	content := ReadContent(storage, fsdata, inodeInfo)
	dirContent := diracts.DecodeDirectoryContent(content)
	return dirContent
}

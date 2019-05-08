package core

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/bojodimitrov/byfiri/logger"
	"github.com/bojodimitrov/byfiri/util"

	"github.com/bojodimitrov/byfiri/diracts"
	"github.com/bojodimitrov/byfiri/structures"
)

//ReadInode returns file inode
func ReadInode(storage []byte, metadata *structures.Metadata, inode int) (*structures.Inode, error) {
	inodesBeginning := int(metadata.Root)
	inodeLocation := inodesBeginning + inode*int(structures.InodeSize)
	//First 3 bytes are mode
	modeStr := Read(storage, inodeLocation, 3)
	mode, err := strconv.Atoi(modeStr)

	if err != nil {
		logger.Log("read inode: " + err.Error())
		return nil, fmt.Errorf("read inode: corrupt inode data")
	}

	inodeLocation += 3
	//Next 10 bytes are size
	sizeStr := Read(storage, inodeLocation, 10)
	size, err := strconv.Atoi(sizeStr)

	if err != nil {
		logger.Log("read inode: " + err.Error())
		return nil, fmt.Errorf("read inode: corrupt inode data")
	}

	inodeLocation += 10

	var blocksGathered [12]uint32
	//Next 12 * 10 bytes are blocks of file
	for i := 0; i < 12; i++ {
		block := Read(storage, inodeLocation, 10)
		value, err := strconv.Atoi(block)
		if err != nil {
			logger.Log("read inode: " + err.Error())
			return nil, fmt.Errorf("read inode: corrupt inode data")
		}
		blocksGathered[i] = uint32(value)
		inodeLocation += 10
	}
	inodeInfo := structures.Inode{Mode: uint8(mode), Size: uint32(size), BlocksLocations: blocksGathered}
	return &inodeInfo, nil
}

func readContent(storage []byte, metadata *structures.Metadata, inodeInfo *structures.Inode) string {
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
func ReadFile(storage []byte, inode int) (string, error) {
	if inode == 0 {
		return "", fmt.Errorf("read file: inode cannot be 0")
	}

	fsdata := ReadMetadata(storage)
	if !getInodeValue(storage, fsdata, structures.Inodes, inode) {
		return "", fmt.Errorf("read file: file does not exits")
	}
	inodeInfo, err := ReadInode(storage, fsdata, inode)
	if err != nil {
		logger.Log("read file: " + err.Error())
		return "", fmt.Errorf("read file: could not read inode")
	}
	if inodeInfo.Mode == 0 {
		return "", fmt.Errorf("read file: file is directory")
	}
	return readContent(storage, fsdata, inodeInfo), nil
}

//ReadDirectory returns directory content
func ReadDirectory(storage []byte, inode int) ([]structures.DirectoryEntry, error) {

	if inode == 0 {
		return nil, fmt.Errorf("read directory: inode cannot be 0")
	}

	fsdata := ReadMetadata(storage)
	if !getInodeValue(storage, fsdata, structures.Inodes, inode) {
		return nil, fmt.Errorf("read directory: directory does not exists")
	}
	inodeInfo, err := ReadInode(storage, fsdata, inode)
	if err != nil {
		logger.Log("read directory: " + err.Error())
		return nil, fmt.Errorf("read directory: could not read inode")
	}
	if inodeInfo.Mode == 1 {
		return nil, fmt.Errorf("read directory: directory is file")
	}

	content := readContent(storage, fsdata, inodeInfo)
	dirContent := diracts.DecodeDirectoryContent(content)
	return dirContent, nil
}

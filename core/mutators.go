package core

import (
	"fmt"
	"math"
	"strings"

	"github.com/bojodimitrov/gofys/structures"
)

func updateContent(storage []byte, fsdata *structures.Metadata, inode *structures.Inode, content string) ([]int, error) {
	numberOfRequiredBlocks := int(math.Ceil(float64(len(content)) / float64(fsdata.BlockSize)))
	if numberOfRequiredBlocks == 0 {
		numberOfRequiredBlocks = 1
	}
	if numberOfRequiredBlocks >= 12 {
		return nil, fmt.Errorf("file is too long")
	}
	numberOfTakenBlocks := int(math.Ceil(float64(inode.Size) / float64(fsdata.BlockSize)))
	if numberOfTakenBlocks == 0 {
		numberOfTakenBlocks = 1
	}

	var gatheredBlocks []int
	for _, value := range inode.BlocksLocations {
		if value != 0 {
			gatheredBlocks = append(gatheredBlocks, int(value))
		}
	}

	if numberOfTakenBlocks > numberOfRequiredBlocks {
		gatheredBlocks = gatheredBlocks[:len(gatheredBlocks)-(numberOfTakenBlocks-numberOfRequiredBlocks)]
	}

	if numberOfTakenBlocks < numberOfRequiredBlocks {
		for i := numberOfTakenBlocks; i < numberOfRequiredBlocks; i++ {
			freeBlock, err := findFreeBitmapPosition(storage, fsdata, structures.Blocks, gatheredBlocks)
			if err != nil {
				return nil, err
			}
			gatheredBlocks = append(gatheredBlocks, freeBlock)
		}
	}

	for _, value := range gatheredBlocks {
		blocksBeginning := int(fsdata.FirstBlock)
		offset := blocksBeginning + value*int(fsdata.BlockSize)
		Write(storage, content, offset)
		markOnBitmap(storage, fsdata, true, value, structures.Blocks)
	}
	return gatheredBlocks, nil
}

func clearFile(storage []byte, blocks [12]uint32, fsdata *structures.Metadata) {
	for _, block := range blocks {
		Write(storage, strings.Repeat("\x00", int(fsdata.BlockSize)), int(fsdata.FirstBlock+block*fsdata.BlockSize))
	}
}

func updateBlockIdsInInode(inodeInfo *structures.Inode, blocks []int) {
	for i := range inodeInfo.BlocksLocations {
		inodeInfo.BlocksLocations[i] = 0
	}
	for i, value := range blocks {
		inodeInfo.BlocksLocations[i] = uint32(value)
	}
}

func updateInode(storage []byte, fsdata *structures.Metadata, inodeInfo *structures.Inode, inode int) {
	inodesBeginning := int(fsdata.Root)
	offset := inodesBeginning + inode*structures.InodeSize

	Write(storage, fmt.Sprint(inodeInfo.Mode), offset)
	offset += 3
	Write(storage, fmt.Sprint(inodeInfo.Size), offset)
	offset += 10
	for i := 0; i < 12; i++ {
		Write(storage, fmt.Sprint(inodeInfo.BlocksLocations[i]), offset)
		offset += 10
	}
}

func clearInode(storage []byte, fsdata *structures.Metadata, inode int) {
	inodesBeginning := int(fsdata.Root)
	offset := inodesBeginning + inode*structures.InodeSize

	Write(storage, strings.Repeat("\x00", structures.InodeSize), offset)
}

//UpdateFile updates file content
func UpdateFile(storage []byte, inode int, content string) {
	fsdata := ReadMetadata(storage)
	inodeInfo := ReadInode(storage, fsdata, inode)
	clearFile(storage, inodeInfo.BlocksLocations, &fsdata)
	blocks, err := updateContent(storage, &fsdata, &inodeInfo, content)
	if err != nil {
		fmt.Println(err)
		return
	}
	updateBlockIdsInInode(&inodeInfo, blocks)
	inodeInfo.Size = uint32(len(content))
	clearInode(storage, &fsdata, inode)
	updateInode(storage, &fsdata, &inodeInfo, inode)
}

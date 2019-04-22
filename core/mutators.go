package core

import (
	"fmt"
	"math"
	"strings"

	"github.com/bojodimitrov/byfiri/diracts"
	"github.com/bojodimitrov/byfiri/structures"
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

	for i, value := range gatheredBlocks {
		blocksBeginning := int(fsdata.FirstBlock)
		offset := blocksBeginning + value*int(fsdata.BlockSize)
		Write(storage, cutContent(content, i, int(fsdata.BlockSize)), offset)
		markOnBitmap(storage, fsdata, true, value, structures.Blocks)
	}
	return gatheredBlocks, nil
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

//UpdateFile updates file content
func UpdateFile(storage []byte, inode int, content string) {
	if inode == 0 {
		fmt.Println("update file: inode cannot be 0")
		return
	}
	fsdata := ReadMetadata(storage)
	inodeInfo := ReadInode(storage, fsdata, inode)
	if inodeInfo.Mode == 0 {
		fmt.Println("update file: file is directory")
		return
	}
	clearFile(storage, inodeInfo.BlocksLocations, fsdata)
	blocks, err := updateContent(storage, fsdata, inodeInfo, content)
	if err != nil {
		fmt.Println(err)
		return
	}
	updateBlockIdsInInode(inodeInfo, blocks)
	inodeInfo.Size = uint32(len(content))
	clearInode(storage, fsdata, inode)
	updateInode(storage, fsdata, inodeInfo, inode)
}

//RenameFile renames file
func RenameFile(storage []byte, currentDirectory *structures.DirectoryIterator, inode int, newName string) {
	if inode == 0 {
		fmt.Println("rename file: inode cannot be 0")
		return
	}
	if strings.ContainsAny(newName, "\\: ") {
		fmt.Println("rename file: name cannot contain ", []string{"'\\'", "' '", "':'"})
		return
	}
	for i, dirEntry := range currentDirectory.DirectoryContent {
		if dirEntry.Inode == uint32(inode) {
			currentDirectory.DirectoryContent[i].FileName = newName
			UpdateDirectory(storage, currentDirectory.DirectoryInode, currentDirectory.DirectoryContent)
			return
		}
	}
	fmt.Println("rename file: file not found")
}

//UpdateDirectory updates file content
func UpdateDirectory(storage []byte, inode int, content []structures.DirectoryEntry) {
	if inode == 0 {
		fmt.Println("update directory: inode cannot be 0")
		return
	}
	fsdata := ReadMetadata(storage)
	inodeInfo := ReadInode(storage, fsdata, inode)
	if inodeInfo.Mode == 1 {
		fmt.Println("update directory: directory is file")
		return
	}
	clearFile(storage, inodeInfo.BlocksLocations, fsdata)

	encoded, err := diracts.EncodeDirectoryContent(content)
	blocks, err := updateContent(storage, fsdata, inodeInfo, encoded)
	if err != nil {
		fmt.Println(err)
		return
	}
	updateBlockIdsInInode(inodeInfo, blocks)
	inodeInfo.Size = uint32(len(encoded))
	clearInode(storage, fsdata, inode)
	updateInode(storage, fsdata, inodeInfo, inode)
}

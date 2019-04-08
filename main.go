package main

import (
	"fmt"
	"os"
	"strconv"
)

// Basic structures needed for the fs:
//	* File system metadata:
//				Contains information for the file system: storage size, block size, location of root inode, etc.
// 	* Inodes array:
//				all inodes will be stored there, an inode can be either free or describing a file
// 				the ID of the inode will be its index in the inodes array
// 	* Free space bitmap:
// 				an array describing free blocks. The length of the array is
// 				the number of all blocks. Since we use byte array, an index of the byte array
// 				could contain information for 8 blocks: for example if we have 64 blocks, we can
//				squeeze the bitmap in 8 locations: the first index will contain information of blocks
//				0-7, the second index- of blocks 8-15, etc.
//	* Blocks space:
//				The space where the data blocks will be stored
// Current planned structure of the file system array:
// [[FS Metadata] [Inodes array] [Free space bitmap] [Data blocks]]

const defaultBlockSize = 4096

type metadata struct { // Will be written in the beginning of the storage space
	storageSize  uint32
	blockSize    uint32
	blockCount   uint32
	root         uint32 // Location of the inode of root dir
	freeSpaceMap uint32 // Beginning of free space bitmap, its len will be blockSize / 8
}

// So far the metadata struct will need 5*10 = 50 spaces to be fitted in the storage space
// Let us define that as a constant, which will be used to calculate required space for the metadata
const metadataSize = 50

type inode struct {
	mode            uint8 // Determines wheter it is a file or directory
	size            uint32
	blocksLocations [12]uint32
}

// So far the inode struct will need 3 + 10 + 12*10 = 133 spaces to be fitted in the storage space
// Let us define that as a constant, which will be used to calculate required space for the inodes array
const inodeSize = 133

// Any change in the structs means changing the corresponding constant

// In order to determine the count of the inodes, thus defining the max files count,
// we will use this simple formula: {size of storage in bytes} / {block size} / 4
// any optimisation of the formila will be appreciated

func calculateInodesCount(storageSize int, blockSize int) int {
	return storageSize / blockSize / 4
}

// Now we have to calculate the size that the inodes array will use
func calculateInodesSpace(inodesCount int) int {
	return inodesCount * 133
}

// Now what remains is to fill in the remaining space with the data blocks
// but we have to save some space for the free space bitmap array
// We will do that bu calculating the remaining free space after assigning some
// for the metadata and the inodes array, calculate how many blocks could be fitted,
// based on that calculate how much space the bitmap will use, subtract from the free space
// and recalculate everything. Some insignificatn space will remain unusable
func calculateNumberOfBlocks(freeSpace int, blockSize int) int {
	fittableBlocks := freeSpace / blockSize
	bitmapBlockSpace := fittableBlocks / 8 / blockSize
	return fittableBlocks - bitmapBlockSpace
}

func initFsSpace(size int) []byte {
	storage := make([]byte, size)
	return storage
}

func readArgs() (int, int) {
	if !checkArguments() {
		fmt.Println("missing size configuration")
		os.Exit(3)
	}
	size := getSize()
	blockSize := getBlockSize()
	return size, blockSize
}

func getSize() int {
	requestedSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("incorrect size format:", os.Args[1])
		os.Exit(3)
	}
	sizeAbbreviation := os.Args[2]

	switch sizeAbbreviation {
	case "GB":
		requestedSize = requestedSize * 1073741824
	case "MB":
		requestedSize = requestedSize * 1048576
	}
	return requestedSize
}

func getBlockSize() int {
	if len(os.Args) == 3 {
		return defaultBlockSize
	}
	requestedSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("incorrect size format:", os.Args[1])
		os.Exit(3)
	}

	return requestedSize
}

func checkArguments() bool {
	numberOfArgs := len(os.Args)
	if numberOfArgs < 3 || numberOfArgs > 4 {
		return false
	}
	return true
}

func main() {
	size, blockSize := readArgs()
	storage := initFsSpace(size)
	sizeString := strconv.Itoa(size)
	write(storage, sizeString, 0)
	fmt.Println(read(storage, 0, 3))
	fmt.Println(blockSize)
	fmt.Println(blockSize)
}

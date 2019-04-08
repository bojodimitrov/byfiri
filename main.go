package main

import (
	"GoFyS/structures"
	"GoFyS/utility"
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
		return structures.DefaultBlockSize
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
	storage := utility.InitFsSpace(size)
	utility.AllocateAllStructures(storage, size, blockSize)
	fmt.Println(utility.ReadMetadata(storage))
}

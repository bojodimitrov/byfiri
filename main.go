package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bojodimitrov/byfiri/core"
	"github.com/bojodimitrov/byfiri/structures"
)

// Basic structures needed for the fs:
//	* File system metadata:
//				Contains information for the file system: storage size, block size, location of root inode, etc.
//	* Free inodes bitmap:
//				an array describing free inodes. The length of the array is
// 				the number of all inodes. Since we use byte array, an index of the byte array
// 				could contain information for 8 inodes: for example if we have 64 inodes, we can
//				squeeze the bitmap in 8 locations: the first index will contain information of inodes
//				0-7, the second index- of inodes 8-15, etc.
// 	* Inodes array:
//				all inodes will be stored there, an inode can be either free or describing a file
// 				the ID of the inode will be its index in the inodes array
// 	* Free blocks bitmap:
// 				an array describing free blocks. Works same way as inodes bitmap
//	* Blocks space:
//				The space where the data blocks will be stored
// Current planned structure of the file system array:
// [[FS Metadata] [Inodes bitmap] [Inodes array] [Free space bitmap] [Data blocks]]

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
	storage := core.InitFsSpace(size)
	dir := core.AllocateAllStructures(storage, size, blockSize)

	inodae := core.AllocateDirectory(storage, dir, "root lv1 dir1")
	core.AllocateDirectory(storage, dir, "root lv1 dir2")
	core.AllocateFile(storage, dir, "root lv1 f1", "man of culture")

	dir, _ = core.EnterDirectory(storage, dir, "root lv1 dir1")
	core.AllocateDirectory(storage, dir, "root lv2 dir1")
	core.AllocateFile(storage, dir, "root lv2 f1", "hello there")
	core.AllocateFile(storage, dir, "root lv2 f2", "thanos did nothing wrong")

	dir, _ = core.EnterDirectory(storage, dir, "root lv2 dir1")
	core.AllocateFile(storage, dir, "root lv3 f1", "i am the senate")

	dir, _ = core.EnterDirectory(storage, dir, "..")
	dir, _ = core.EnterDirectory(storage, dir, "..")
	core.DeleteDirectory(storage, dir, inodae)
	inode := core.AllocateFile(storage, dir, "bleh", "bleh")
	inode2 := core.AllocateFile(storage, dir, "bleh1", "thanos did nothing wrong 22323")
	inode3 := core.AllocateFile(storage, dir, "bleh2", "thanos did nothing wrong 324234 324 ")
	inode4 := core.AllocateFile(storage, dir, "bleh3", "thanos did nothing wrong 234f ssdf ")
	inode5 := core.AllocateFile(storage, dir, "bleh4", "thanos did nothing wrong sdffsd s e r3")
	fmt.Println(core.ReadDirectory(storage, dir.DirectoryInode))
	fmt.Println(inode, inode2, inode3, inode4, inode5)
	fmt.Println(core.ReadFile(storage, inode))
	fmt.Println(core.ReadFile(storage, inode2))
	fmt.Println(core.ReadFile(storage, inode3))
	fmt.Println(core.ReadFile(storage, inode4))
	fmt.Println(core.ReadFile(storage, inode5))
}

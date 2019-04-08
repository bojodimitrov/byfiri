package utility

import (
	"GoFyS/structures"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	return fittableBlocks - bitmapBlockSpace - (fittableBlocks-bitmapBlockSpace)%8
}

func fillWithZeros(str string, upToSize int) string {
	zeros := strings.Repeat("0", upToSize-len(str))
	return zeros + str
}

// AllocateMetadata writes metadata struct on storage
func AllocateMetadata(storage []byte, fsdata *structures.Metadata) {
	Write(storage, fillWithZeros(fmt.Sprint(fsdata.StorageSize), 20), 0)
	Write(storage, fillWithZeros(fmt.Sprint(fsdata.BlockSize), 10), 20)
	Write(storage, fillWithZeros(fmt.Sprint(fsdata.BlockCount), 10), 30)
	Write(storage, fillWithZeros(fmt.Sprint(fsdata.Root), 10), 40)
	Write(storage, fillWithZeros(fmt.Sprint(fsdata.FreeSpaceMap), 10), 50)
	Write(storage, fillWithZeros(fmt.Sprint(fsdata.FirstBlock), 10), 60)
}

// ReadMetadata read metadata information from storage and returns it
func ReadMetadata(storage []byte) structures.Metadata {
	var input [6]int
	sizeStr, err := Read(storage, 0, 20)
	throwCorruptError(err)
	size, err := strconv.Atoi(sizeStr)
	throwCorruptError(err)
	input[0] = size
	for i := 1; i < 6; i++ {
		valueStr, err := Read(storage, (i+1)*10, 10)
		throwCorruptError(err)
		value, err := strconv.Atoi(valueStr)
		throwCorruptError(err)
		input[i] = value
	}

	fsdata := structures.Metadata{StorageSize: uint64(input[0]),
		BlockSize:    uint32(input[1]),
		BlockCount:   uint32(input[2]),
		Root:         uint32(input[3]),
		FreeSpaceMap: uint32(input[4]),
		FirstBlock:   uint32(input[5])}
	return fsdata
}

func throwCorruptError(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Println("corrupt metadata")
		os.Exit(3)
	}
}

// AllocateInode writes an inode struct on storage
func AllocateInode(storage []byte, iNodeInfo *structures.Inode, location int) {
	offset := location
	Write(storage, fmt.Sprint(iNodeInfo.Mode), offset)
	offset += 3
	Write(storage, fmt.Sprint(iNodeInfo.Size), offset)
	offset += 10
	for i := 0; i < 12; i++ {
		Write(storage, fmt.Sprint(iNodeInfo.BlocksLocations[i]), offset)
		offset += 10
	}
}

func createRoot() {

}

// InitFsSpace creates storage array in memory
func InitFsSpace(size int) []byte {
	storage := make([]byte, size)
	return storage
}

// AllocateAllStructures writes all basic structures on storage
func AllocateAllStructures(storage []byte, size int, blockSize int) structures.Metadata {
	inodes := calculateInodesSpace(calculateInodesCount(size, blockSize))
	numberOfBlocks := uint32(calculateNumberOfBlocks(size, blockSize))
	fsdata := structures.Metadata{StorageSize: uint64(size),
		BlockSize:    uint32(blockSize),
		BlockCount:   numberOfBlocks,
		Root:         structures.MetadataSize,
		FreeSpaceMap: structures.MetadataSize + uint32(inodes),
		FirstBlock:   structures.MetadataSize + uint32(inodes) + numberOfBlocks/8}
	AllocateMetadata(storage, &fsdata)
	createRoot()
	return fsdata
}

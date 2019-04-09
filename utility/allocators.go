package utility

import (
	"GoFyS/error_handling"
	"GoFyS/structures"
	"fmt"
	"math"
	"strconv"
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

// AllocateMetadata writes metadata struct on storage
func AllocateMetadata(storage []byte, fsdata *structures.Metadata) {
	Write(storage, fmt.Sprint(fsdata.StorageSize), 0)
	Write(storage, fmt.Sprint(fsdata.InodesCount), 20)
	Write(storage, fmt.Sprint(fsdata.BlockSize), 30)
	Write(storage, fmt.Sprint(fsdata.BlockCount), 40)
	Write(storage, fmt.Sprint(fsdata.Root), 50)
	Write(storage, fmt.Sprint(fsdata.FreeSpaceMap), 60)
	Write(storage, fmt.Sprint(fsdata.FirstBlock), 70)
}

// ReadMetadata read metadata information from storage and returns it
func ReadMetadata(storage []byte) structures.Metadata {
	var input [7]int
	sizeStr, err := Read(storage, 0, 20)
	error_handling.CorruptMetadata(err)
	size, err := strconv.Atoi(sizeStr)
	error_handling.CorruptMetadata(err)
	input[0] = size
	for i := 1; i < 7; i++ {
		valueStr, err := Read(storage, (i+1)*10, 10)
		error_handling.CorruptMetadata(err)
		value, err := strconv.Atoi(valueStr)
		error_handling.CorruptMetadata(err)
		input[i] = value
	}

	fsdata := structures.Metadata{StorageSize: uint64(input[0]),
		InodesCount:  uint32(input[1]),
		BlockSize:    uint32(input[2]),
		BlockCount:   uint32(input[3]),
		Root:         uint32(input[4]),
		FreeSpaceMap: uint32(input[5]),
		FirstBlock:   uint32(input[6])}
	return fsdata
}

// AllocateInode writes an inode struct on storage
func AllocateInode(storage []byte, inodeInfo *structures.Inode, inodeNumber int) {
	fsdata := ReadMetadata(storage)
	inodesBeginning := int(fsdata.Root)
	offset := inodesBeginning + inodeNumber*structures.InodeSize

	Write(storage, fmt.Sprint(inodeInfo.Mode), offset)
	offset += 3
	Write(storage, fmt.Sprint(inodeInfo.Size), offset)
	offset += 10
	for i := 0; i < 12; i++ {
		Write(storage, fmt.Sprint(inodeInfo.BlocksLocations[i]), offset)
		offset += 10
	}
}

// AllocateContent writes a file content in a block on storage
func AllocateContent(storage []byte, content string) ([]int, error) {
	fsdata := ReadMetadata(storage)
	numberOfRequiredBlocks := int(math.Ceil(float64(len(content)) / float64(fsdata.BlockSize)))
	var gatheredBlocks []int
	for i := 0; i < numberOfRequiredBlocks; i++ {
		freeBlock, err := findFreeBitmapPosition(storage, structures.Blocks)
		if err != nil {
			return nil, err
		}
		gatheredBlocks = append(gatheredBlocks, freeBlock)
	}
	for _, value := range gatheredBlocks {
		blocksBeginning := int(fsdata.FirstBlock)
		offset := blocksBeginning + value*int(fsdata.BlockSize)
		Write(storage, content, offset)
	}
	return gatheredBlocks, nil
}

func readBitmap(storage []byte, bitmap structures.Bitmap) []bool {
	fsdata := ReadMetadata(storage)
	bitmapStart := 0
	bitmapLength := 0
	switch bitmap {
	case structures.Inodes:
		bitmapStart = int(fsdata.InodesMap)
		bitmapLength = int(fsdata.InodesCount / 8)
	case structures.Blocks:
		bitmapStart = int(fsdata.FreeSpaceMap)
		bitmapLength = int(fsdata.BlockCount / 8)
	default:
		error_handling.IncorrectFormat("Bitmap option")
	}

	bitmapHexStr, err := ReadRaw(storage, bitmapStart, bitmapLength)
	if err != nil {
		error_handling.CorruptBitmap(err, "inodes")
	}
	bitmapArray := StringToBin(bitmapHexStr)
	return bitmapArray
}

func findFreeBitmapPosition(storage []byte, bitmap structures.Bitmap) (int, error) {
	inodesBitmap := readBitmap(storage, bitmap)
	freePosition := -1
	for i, val := range inodesBitmap {
		if val == false {
			freePosition = i
			break
		}
	}
	if freePosition == -1 {
		return 0, fmt.Errorf("all inodes taken")
	}
	return freePosition, nil
}

// AllocateFile writes a file on storage
func AllocateFile(storage []byte, inodeInfo *structures.Inode, content string) {
	blocksGathered := AllocateContent(storage, content)

	freeInode, err := findFreeBitmapPosition(storage, structures.Inodes)
	if err != nil {
		fmt.Println(err)
		return
	}
	AllocateInode(storage, inodeInfo, freeInode)

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
	inodesCount := calculateInodesCount(size, blockSize)
	inodes := calculateInodesSpace(inodesCount)
	numberOfBlocks := uint32(calculateNumberOfBlocks(size, blockSize))
	fsdata := structures.Metadata{StorageSize: uint64(size),
		InodesCount:  uint32(inodesCount),
		BlockSize:    uint32(blockSize),
		BlockCount:   numberOfBlocks,
		InodesMap:    structures.MetadataSize,
		Root:         structures.MetadataSize + uint32(inodesCount)/8,
		FreeSpaceMap: structures.MetadataSize + uint32(inodesCount)/8 + uint32(inodes),
		FirstBlock:   structures.MetadataSize + uint32(inodesCount)/8 + uint32(inodes) + numberOfBlocks/8}
	AllocateMetadata(storage, &fsdata)
	createRoot()
	return fsdata
}

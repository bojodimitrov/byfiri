package core

import (
	"fmt"
	"math"

	"github.com/bojodimitrov/gofys/diracts"
	"github.com/bojodimitrov/gofys/structures"
)

// In order to determine the count of the inodes, thus defining the max files count,
// we will use this simple formula: {size of storage in bytes} / {block size} / 4
// any optimisation of the formila will be appreciated

func calculateInodesCount(storageSize int, blockSize int) int {
	return storageSize / blockSize / 4
}

// Now we have to calculate the size that the inodes array will use
func calculateInodesSpace(inodesCount int) int {
	return inodesCount * structures.InodeSize
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

// InitFsSpace creates storage array in memory
func InitFsSpace(size int) []byte {
	storage := make([]byte, size)
	return storage
}

// AllocateMetadata writes metadata struct on storage
func AllocateMetadata(storage []byte, fsdata *structures.Metadata) {
	Write(storage, fmt.Sprint(fsdata.StorageSize), 0)
	Write(storage, fmt.Sprint(fsdata.InodesCount), 20)
	Write(storage, fmt.Sprint(fsdata.BlockSize), 30)
	Write(storage, fmt.Sprint(fsdata.BlockCount), 40)
	Write(storage, fmt.Sprint(fsdata.InodesMap), 50)
	Write(storage, fmt.Sprint(fsdata.Root), 60)
	Write(storage, fmt.Sprint(fsdata.FreeSpaceMap), 70)
	Write(storage, fmt.Sprint(fsdata.FirstBlock), 80)
}

// AllocateInode writes an inode struct on storage
func AllocateInode(storage []byte, inodeInfo *structures.Inode) int {
	freeInode, err := findFreeBitmapPosition(storage, structures.Inodes, []int{})
	if err != nil {
		fmt.Println(err)
		return 0
	}

	fsdata := ReadMetadata(storage)
	inodesBeginning := int(fsdata.Root)
	offset := inodesBeginning + freeInode*structures.InodeSize

	Write(storage, fmt.Sprint(inodeInfo.Mode), offset)
	offset += 3
	Write(storage, fmt.Sprint(inodeInfo.Size), offset)
	offset += 10
	for i := 0; i < 12; i++ {
		Write(storage, fmt.Sprint(inodeInfo.BlocksLocations[i]), offset)
		offset += 10
	}
	markOnBitmap(storage, true, inodesBeginning, freeInode, structures.Inodes)
	return freeInode
}

// AllocateContent writes a file content in a block on storage
func AllocateContent(storage []byte, content string) ([]int, error) {
	fsdata := ReadMetadata(storage)
	numberOfRequiredBlocks := int(math.Ceil(float64(len(content)) / float64(fsdata.BlockSize)))
	if numberOfRequiredBlocks == 0 {
		numberOfRequiredBlocks = 1
	}
	var gatheredBlocks []int
	for i := 0; i < numberOfRequiredBlocks; i++ {
		freeBlock, err := findFreeBitmapPosition(storage, structures.Blocks, gatheredBlocks)
		if err != nil {
			return nil, err
		}
		gatheredBlocks = append(gatheredBlocks, freeBlock)
	}
	for _, value := range gatheredBlocks {
		blocksBeginning := int(fsdata.FirstBlock)
		offset := blocksBeginning + value*int(fsdata.BlockSize)
		Write(storage, content, offset)
		markOnBitmap(storage, true, int(fsdata.InodesMap), value, structures.Blocks)
	}
	return gatheredBlocks, nil
}

// AllocateFile writes a file on storage
func AllocateFile(storage []byte, mode uint8, content string) int {
	blocksGathered, err := AllocateContent(storage, content)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	inodeInfo := structures.Inode{Mode: mode, Size: uint32(len(content)), BlocksLocations: [12]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	addBlockIdsInInode(&inodeInfo, blocksGathered)

	return AllocateInode(storage, &inodeInfo)
}

func findFreeBitmapPosition(storage []byte, bitmap structures.Bitmap, inMemoryTakenPositions []int) (int, error) {
	inodesBitmap := GetBitmap(storage, bitmap)
	freePosition := -1
	for i, val := range inodesBitmap {
		if i != 0 && val == false && !Contains(inMemoryTakenPositions, i) {
			freePosition = i
			break
		}
	}
	if freePosition == -1 {
		return 0, fmt.Errorf("all inodes taken")
	}
	return freePosition, nil
}

func addBlockIdsInInode(inodeInfo *structures.Inode, blocks []int) {
	for i, value := range blocks {
		inodeInfo.BlocksLocations[i] = uint32(value)
	}
}

func markOnBitmap(storage []byte, value bool, offset int, location int, bitmap structures.Bitmap) {
	byteOctet := GetBitmapIndex(storage, bitmap, location/8)
	byteOctet[location%8] = value
	val := BinToByteValue(byteOctet)
	WriteBitmap(storage, bitmap, val, location/8)
}

func createRoot(storage []byte) {
	fmt.Println("creating root")
	content, err := diracts.EncodeDirectoryContent([]structures.DirectoryContent{structures.DirectoryContent{FileName: ".", Inode: 1}})
	if err != nil {
		return
	}
	AllocateFile(storage, 0, content)
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
	createRoot(storage)
	return fsdata
}

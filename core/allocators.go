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

// AllocateFile writes a file on storage
func AllocateFile(storage []byte, mode uint8, content string) int {
	fsdata := ReadMetadata(storage)
	blocksGathered, err := allocateContent(storage, &fsdata, content)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	inodeInfo := structures.Inode{Mode: mode, Size: uint32(len(content)), BlocksLocations: [12]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	addBlockIdsInInode(&inodeInfo, blocksGathered)

	return allocateInode(storage, &fsdata, &inodeInfo)
}

// AllocateDirectory writes a directory on storage
func AllocateDirectory(storage []byte, mode uint8, content []structures.DirectoryContent) int {
	fsdata := ReadMetadata(storage)
	encoded, err := diracts.EncodeDirectoryContent(content)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	blocksGathered, err := allocateContent(storage, &fsdata, encoded)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	inodeInfo := structures.Inode{Mode: mode, Size: uint32(len(encoded)), BlocksLocations: [12]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	addBlockIdsInInode(&inodeInfo, blocksGathered)

	return allocateInode(storage, &fsdata, &inodeInfo)
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
	allocateMetadata(storage, &fsdata)
	createRoot(storage)
	return fsdata
}

// InitFsSpace creates storage array in memory
func InitFsSpace(size int) []byte {
	storage := make([]byte, size)
	return storage
}

func allocateMetadata(storage []byte, fsdata *structures.Metadata) {
	Write(storage, fmt.Sprint(fsdata.StorageSize), 0)
	Write(storage, fmt.Sprint(fsdata.InodesCount), 20)
	Write(storage, fmt.Sprint(fsdata.BlockSize), 30)
	Write(storage, fmt.Sprint(fsdata.BlockCount), 40)
	Write(storage, fmt.Sprint(fsdata.InodesMap), 50)
	Write(storage, fmt.Sprint(fsdata.Root), 60)
	Write(storage, fmt.Sprint(fsdata.FreeSpaceMap), 70)
	Write(storage, fmt.Sprint(fsdata.FirstBlock), 80)
}

func allocateInode(storage []byte, fsdata *structures.Metadata, inodeInfo *structures.Inode) int {
	freeInode, err := findFreeBitmapPosition(storage, fsdata, structures.Inodes, []int{})
	if err != nil {
		fmt.Println(err)
		return 0
	}

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
	markOnBitmap(storage, fsdata, true, freeInode, structures.Inodes)
	return freeInode
}

func allocateContent(storage []byte, fsdata *structures.Metadata, content string) ([]int, error) {
	numberOfRequiredBlocks := int(math.Ceil(float64(len(content)) / float64(fsdata.BlockSize)))
	if numberOfRequiredBlocks == 0 {
		numberOfRequiredBlocks = 1
	}
	var gatheredBlocks []int
	for i := 0; i < numberOfRequiredBlocks; i++ {
		freeBlock, err := findFreeBitmapPosition(storage, fsdata, structures.Blocks, gatheredBlocks)
		if err != nil {
			return nil, err
		}
		gatheredBlocks = append(gatheredBlocks, freeBlock)
	}
	for _, value := range gatheredBlocks {
		blocksBeginning := int(fsdata.FirstBlock)
		offset := blocksBeginning + value*int(fsdata.BlockSize)
		Write(storage, content, offset)
		markOnBitmap(storage, fsdata, true, value, structures.Blocks)
	}
	return gatheredBlocks, nil
}

func findFreeBitmapPosition(storage []byte, fsdata *structures.Metadata, bitmap structures.Bitmap, inMemoryTakenPositions []int) (int, error) {
	inodesBitmap := GetBitmap(storage, fsdata, bitmap)
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

func markOnBitmap(storage []byte, fsdata *structures.Metadata, value bool, location int, bitmap structures.Bitmap) {
	byteOctet := GetBitmapIndex(storage, fsdata, bitmap, location/8)
	byteOctet[location%8] = value
	val := BinToByteValue(byteOctet)
	WriteBitmap(storage, fsdata, bitmap, val, location/8)
}

func createRoot(storage []byte) {
	fmt.Println("creating root")
	content := []structures.DirectoryContent{structures.DirectoryContent{FileName: ".", Inode: 1}}
	AllocateDirectory(storage, 0, content)
}

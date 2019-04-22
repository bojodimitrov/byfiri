package core

import (
	"fmt"
	"math"
	"strings"

	"github.com/bojodimitrov/byfiri/diracts"
	"github.com/bojodimitrov/byfiri/structures"
	"github.com/bojodimitrov/byfiri/util"
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
func AllocateFile(storage []byte, currentDirectory *structures.DirectoryIterator, name string, content string) int {
	if fileAlreadyExists(currentDirectory.DirectoryContent, name) {
		fmt.Println("allocate file: name already exists")
		return 0
	}
	if strings.ContainsAny(name, "\\:") {
		fmt.Println("allocate file: name cannot contain ", []string{"'\\'", "':'"})
		return 0
	}

	fsdata := ReadMetadata(storage)
	blocksGathered, err := allocateContent(storage, fsdata, content)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	inodeInfo := structures.Inode{Mode: 1, Size: uint32(len(content)), BlocksLocations: [12]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	addBlockIdsInInode(&inodeInfo, blocksGathered)

	inode := allocateInode(storage, fsdata, &inodeInfo)
	addFileToDirectory(storage, currentDirectory, &structures.DirectoryEntry{FileName: name, Inode: uint32(inode)})
	return inode
}

// AllocateDirectory writes a directory on storage
func AllocateDirectory(storage []byte, currentDirectory *structures.DirectoryIterator, name string) int {
	if fileAlreadyExists(currentDirectory.DirectoryContent, name) {
		fmt.Println("allocate directory: name already exists")
		return 0
	}

	fsdata := ReadMetadata(storage)
	content := []structures.DirectoryEntry{
		structures.DirectoryEntry{FileName: ".", Inode: 0},
		structures.DirectoryEntry{FileName: "..", Inode: uint32(currentDirectory.DirectoryInode)}}
	encoded, err := diracts.EncodeDirectoryContent(content)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	blocksGathered, err := allocateContent(storage, fsdata, encoded)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	inodeInfo := structures.Inode{Mode: 0, Size: uint32(len(encoded)), BlocksLocations: [12]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	addBlockIdsInInode(&inodeInfo, blocksGathered)

	inode := allocateInode(storage, fsdata, &inodeInfo)
	content[0].Inode = uint32(inode)

	encoded, err = diracts.EncodeDirectoryContent(content)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	updateContent(storage, fsdata, &inodeInfo, encoded)
	addFileToDirectory(storage, currentDirectory, &structures.DirectoryEntry{FileName: name, Inode: uint32(inode)})
	return inode
}

// AllocateAllStructures writes all basic structures on storage
func AllocateAllStructures(storage []byte, size int, blockSize int) *structures.DirectoryIterator {
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
	root := createRoot(storage)
	return root
}

// InitFsSpace creates storage array in memory
func InitFsSpace(size int) []byte {
	storage := make([]byte, size)
	return storage
}

func fileAlreadyExists(content []structures.DirectoryEntry, name string) bool {
	for _, entry := range content {
		if entry.FileName == name {
			return true
		}
	}
	return false
}

func addFileToDirectory(storage []byte, current *structures.DirectoryIterator, file *structures.DirectoryEntry) {
	content := ReadDirectory(storage, current.DirectoryInode)
	content = append(content, *file)
	UpdateDirectory(storage, current.DirectoryInode, content)
	current.DirectoryContent = append(current.DirectoryContent, *file)
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
	if numberOfRequiredBlocks >= 12 {
		return nil, fmt.Errorf("file is too long")
	}
	var gatheredBlocks []int
	for i := 0; i < numberOfRequiredBlocks; i++ {
		freeBlock, err := findFreeBitmapPosition(storage, fsdata, structures.Blocks, gatheredBlocks)
		if err != nil {
			return nil, err
		}
		gatheredBlocks = append(gatheredBlocks, freeBlock)
	}
	for i, value := range gatheredBlocks {
		blocksBeginning := int(fsdata.FirstBlock)
		offset := blocksBeginning + value*int(fsdata.BlockSize)
		Write(storage, cutContent(content, i, int(fsdata.BlockSize)), offset)
		markOnBitmap(storage, fsdata, true, value, structures.Blocks)
	}
	return gatheredBlocks, nil
}

func findFreeBitmapPosition(storage []byte, fsdata *structures.Metadata, bitmap structures.Bitmap, inMemoryTakenPositions []int) (int, error) {
	inodesBitmap := GetBitmap(storage, fsdata, bitmap)
	freePosition := -1
	for i, val := range inodesBitmap {
		if i != 0 && val == false && !util.Contains(inMemoryTakenPositions, i) {
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

func createRoot(storage []byte) *structures.DirectoryIterator {
	fmt.Println("creating root")
	content := []structures.DirectoryEntry{structures.DirectoryEntry{FileName: ".", Inode: 1}, structures.DirectoryEntry{FileName: "..", Inode: 0}}

	fsdata := ReadMetadata(storage)
	encoded, err := diracts.EncodeDirectoryContent(content)
	if err != nil {
		fmt.Println(err)
		panic("create root: could not encode root content")
	}
	blocksGathered, err := allocateContent(storage, fsdata, encoded)
	if err != nil {
		fmt.Println(err)
		panic("create root: could not allocate root")
	}
	inodeInfo := structures.Inode{Mode: 0, Size: uint32(len(encoded)), BlocksLocations: [12]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	addBlockIdsInInode(&inodeInfo, blocksGathered)
	allocateInode(storage, fsdata, &inodeInfo)
	//returns first directory iterator
	return &structures.DirectoryIterator{DirectoryInode: 1, DirectoryContent: content}
}

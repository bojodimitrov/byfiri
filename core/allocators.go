package core

import (
	"fmt"
	"math"
	"strings"

	"github.com/bojodimitrov/byfiri/diracts"
	"github.com/bojodimitrov/byfiri/logger"
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
func AllocateFile(storage []byte, currentDirectory *structures.DirectoryIterator, name string, content string) (int, error) {
	if fileAlreadyExists(currentDirectory.DirectoryContent, name) {
		return 0, fmt.Errorf("allocate file: name already exists")
	}
	if strings.ContainsAny(name, "\\:-") {
		return 0, fmt.Errorf("allocate file: name cannot contain %q", []string{"'\\'", "':'", "'-'"})
	}

	fsdata := ReadMetadata(storage)
	blocksGathered, err := allocateContent(storage, fsdata, content)
	if err != nil {
		return 0, err
	}
	inodeInfo := structures.Inode{Mode: 1, Size: uint32(len(content)), BlocksLocations: [12]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	addBlockIdsInInode(&inodeInfo, blocksGathered)

	inode, err := allocateInode(storage, fsdata, &inodeInfo)
	if err != nil {
		return 0, err
	}
	addFileToDirectory(storage, currentDirectory, &structures.DirectoryEntry{FileName: name, Inode: uint32(inode)})
	return inode, nil
}

// AllocateDirectory writes a directory on storage
func AllocateDirectory(storage []byte, currentDirectory *structures.DirectoryIterator, name string) (int, error) {
	if fileAlreadyExists(currentDirectory.DirectoryContent, name) {
		return 0, fmt.Errorf("allocate directory: name already exists")
	}

	fsdata := ReadMetadata(storage)
	content := []structures.DirectoryEntry{
		structures.DirectoryEntry{FileName: ".", Inode: 0},
		structures.DirectoryEntry{FileName: "..", Inode: uint32(currentDirectory.DirectoryInode)}}
	encoded, err := diracts.EncodeDirectoryContent(content)
	if err != nil {
		return 0, err
	}
	blocksGathered, err := allocateContent(storage, fsdata, encoded)
	if err != nil {
		return 0, err
	}
	inodeInfo := structures.Inode{Mode: 0, Size: uint32(len(encoded)), BlocksLocations: [12]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	addBlockIdsInInode(&inodeInfo, blocksGathered)

	inode, err := allocateInode(storage, fsdata, &inodeInfo)
	if err != nil {
		return 0, err
	}
	content[0].Inode = uint32(inode)

	encoded, err = diracts.EncodeDirectoryContent(content)
	if err != nil {
		return 0, err
	}
	updateContent(storage, fsdata, &inodeInfo, encoded)
	addFileToDirectory(storage, currentDirectory, &structures.DirectoryEntry{FileName: name, Inode: uint32(inode)})
	return inode, nil
}

// AllocateAllStructures writes all basic structures on storage
func AllocateAllStructures(storage []byte, size int, blockSize int) (*structures.DirectoryIterator, error) {
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
	root, err := createRoot(storage)
	if err != nil {
		return nil, err
	}
	return root, nil
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

func addFileToDirectory(storage []byte, current *structures.DirectoryIterator, file *structures.DirectoryEntry) error {
	content, err := ReadDirectory(storage, current.DirectoryInode)
	if err != nil {
		return err
	}
	content = append(content, *file)
	UpdateDirectory(storage, current.DirectoryInode, content)
	current.DirectoryContent = append(current.DirectoryContent, *file)
	return nil
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

func allocateInode(storage []byte, fsdata *structures.Metadata, inodeInfo *structures.Inode) (int, error) {
	freeInode, err := findFreeBitmapPosition(storage, fsdata, structures.Inodes, []int{})
	if err != nil {
		return 0, err
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
	return freeInode, nil
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
	inodesBitmap := getBitmap(storage, fsdata, bitmap)
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

func getInodeValue(storage []byte, fsdata *structures.Metadata, bitmap structures.Bitmap, position int) bool {
	inodesBitmap := getBitmap(storage, fsdata, bitmap)
	return inodesBitmap[position]
}

func addBlockIdsInInode(inodeInfo *structures.Inode, blocks []int) {
	for i, value := range blocks {
		inodeInfo.BlocksLocations[i] = uint32(value)
	}
}

func createRoot(storage []byte) (*structures.DirectoryIterator, error) {
	content := []structures.DirectoryEntry{structures.DirectoryEntry{FileName: ".", Inode: 1}, structures.DirectoryEntry{FileName: "..", Inode: 0}}

	fsdata := ReadMetadata(storage)
	encoded, err := diracts.EncodeDirectoryContent(content)
	if err != nil {
		logger.Log("create root: " + err.Error())
		return nil, fmt.Errorf("create root: could not create root content")
	}
	blocksGathered, err := allocateContent(storage, fsdata, encoded)
	if err != nil {
		logger.Log("create root: " + err.Error())
		return nil, fmt.Errorf("create root: could not allocate root")
	}
	inodeInfo := structures.Inode{Mode: 0, Size: uint32(len(encoded)), BlocksLocations: [12]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	addBlockIdsInInode(&inodeInfo, blocksGathered)
	allocateInode(storage, fsdata, &inodeInfo)
	//returns first directory iterator
	return &structures.DirectoryIterator{DirectoryInode: 1, DirectoryContent: content}, nil
}

package structures

const DefaultBlockSize = 4096

type Metadata struct { // Will be written in the beginning of the storage space
	StorageSize  uint64
	BlockSize    uint32
	BlockCount   uint32
	Root         uint32 // Location of the inode of root dir
	FreeSpaceMap uint32 // Beginning of free space bitmap, its len will be blockSize / 8
	FirstBlock   uint32 // location of beginning of first block
}

// So far the metadata struct will need 5*10 = 50 spaces to be fitted in the storage space
// Let us define that as a constant, which will be used to calculate required space for the metadata
const MetadataSize = 70

type Inode struct {
	Mode            uint8 // Determines wheter it is a file or directory: 0 is dir, 1 is file
	Size            uint32
	BlocksLocations [12]uint32
}

// So far the inode struct will need 3 + 10 + 12*10 = 133 spaces to be fitted in the storage space
// Let us define that as a constant, which will be used to calculate required space for the inodes array
const InodeSize = 133

// Any change in the structs means changing the corresponding constant

package structures

// DefaultBlockSize defined when block size is not provided
const DefaultBlockSize = 4096

// Metadata contains metadata of the whole file system
type Metadata struct {
	StorageSize  uint64
	InodesCount  uint32
	BlockSize    uint32
	BlockCount   uint32
	InodesMap    uint32 // Beginning of inode bitmap
	Root         uint32 // Location of the inode of root dir
	FreeSpaceMap uint32 // Beginning of free space bitmap, its len will be blockSize / 8
	FirstBlock   uint32 // location of beginning of first block
}

// MetadataSize is the size of the struct when written on storage space
const MetadataSize = 90

// Inode will contain metadata of a file
type Inode struct {
	Mode            uint8 // Determines wheter it is a file or directory: 0 is dir, 1 is file
	Size            uint32
	BlocksLocations [12]uint32
}

// InodeSize is the size of the struct when written on storage space
const InodeSize = 133

// Any change in the structs means changing the corresponding constant

// Bitmap enum
type Bitmap int

const (
	Inodes Bitmap = 0
	Blocks Bitmap = 1
)

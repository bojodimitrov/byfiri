package core

import (
	"strconv"
	"strings"

	"github.com/bojodimitrov/byfiri/util"

	"github.com/bojodimitrov/byfiri/errors"
	"github.com/bojodimitrov/byfiri/structures"
)

// Write writes in storage array
func Write(storage []byte, content string, offset int) {
	if len(content) == 0 {
		return
	}
	if offset > len(storage) || offset+len(content) > len(storage) {
		panic("operation out of bound")
	}
	for i := 0; i < len(content); i++ {
		storage[i+offset] = content[i]
	}
}

// WriteByte writes in storage array
func WriteByte(storage []byte, content []byte, offset int) {
	if len(content) == 0 {
		return
	}
	if offset > len(storage) || offset+len(content) > len(storage) {
		panic("operation out of bound")
	}
	for i := 0; i < len(content); i++ {
		storage[i+offset] = content[i]
	}
}

// Read reads from storage array
func Read(storage []byte, offset int, length int) string {
	if length < 0 {
		panic("length is negative")
	}
	if offset > len(storage) || offset+length > len(storage) {
		panic("operation out of bound")
	}
	return strings.Replace(string(storage[offset:length+offset]), "\x00", "", -1)
}

// ReadRaw reads from storage array and does not remove x00s
func ReadRaw(storage []byte, offset int, length int) string {
	if length < 0 {
		panic("length is negative")
	}
	if offset > len(storage) || offset+length > len(storage) {
		panic("operation out of bound")
	}
	return string(storage[offset : length+offset])
}

// ReadByte reads from storage array and return byte array
func ReadByte(storage []byte, offset int, length int) []byte {
	if length < 0 {
		panic("length is negative")
	}
	if offset > len(storage) || offset+length > len(storage) {
		panic("operation out of bound")
	}
	return storage[offset : length+offset]
}

// ReadMetadata read metadata information from storage and returns it
func ReadMetadata(storage []byte) *structures.Metadata {
	var input [8]int
	sizeStr := Read(storage, 0, 20)
	size, err := strconv.Atoi(sizeStr)
	errors.CorruptMetadata(err)
	input[0] = size
	for i := 1; i < 8; i++ {
		valueStr := Read(storage, (i+1)*10, 10)
		value, err := strconv.Atoi(valueStr)
		errors.CorruptMetadata(err)
		input[i] = value
	}

	fsdata := structures.Metadata{StorageSize: uint64(input[0]),
		InodesCount:  uint32(input[1]),
		BlockSize:    uint32(input[2]),
		BlockCount:   uint32(input[3]),
		InodesMap:    uint32(input[4]),
		Root:         uint32(input[5]),
		FreeSpaceMap: uint32(input[6]),
		FirstBlock:   uint32(input[7])}
	return &fsdata
}

// WriteBitmap writes the given value on the given index
func WriteBitmap(storage []byte, fsdata *structures.Metadata, bitmap structures.Bitmap, value byte, index int) {
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
		errors.IncorrectFormat("Bitmap option")
	}

	if index >= bitmapLength {
		panic("bitmap index exceeds bitmap length")
	}

	WriteByte(storage, []byte{value}, bitmapStart+index)
}

//GetBitmapIndex returns byte value transformed as binary array at index
func GetBitmapIndex(storage []byte, fsdata *structures.Metadata, bitmap structures.Bitmap, index int) []bool {
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
		errors.IncorrectFormat("Bitmap option")
	}

	if index >= bitmapLength {
		panic("bitmap index exceeds bitmap length")
	}

	byteValue := ReadByte(storage, bitmapStart+index, 1)
	boolOctet := util.ByteToBin(byteValue)
	// bitmapArray contains 8 bits that correspond to the byte index
	return boolOctet
}

//GetBitmap returns whole bitmap as binary array
func GetBitmap(storage []byte, fsdata *structures.Metadata, bitmap structures.Bitmap) []bool {
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
		errors.IncorrectFormat("Bitmap option")
	}

	bitmapHexStr := ReadByte(storage, bitmapStart, bitmapLength)
	bitmapArray := util.ByteToBin(bitmapHexStr)
	return bitmapArray
}

package core

import (
	"strings"

	"github.com/bojodimitrov/gofys/structures"
	"github.com/bojodimitrov/gofys/util"
)

func clearInode(storage []byte, fsdata *structures.Metadata, inode int) {
	inodesBeginning := int(fsdata.Root)
	offset := inodesBeginning + inode*structures.InodeSize

	Write(storage, strings.Repeat("\x00", structures.InodeSize), offset)
}

func clearFile(storage []byte, blocks [12]uint32, fsdata *structures.Metadata) {
	for _, block := range blocks {
		Write(storage, strings.Repeat("\x00", int(fsdata.BlockSize)), int(fsdata.FirstBlock+block*fsdata.BlockSize))
	}
}

func markOnBitmap(storage []byte, fsdata *structures.Metadata, value bool, location int, bitmap structures.Bitmap) {
	byteOctet := GetBitmapIndex(storage, fsdata, bitmap, location/8)
	byteOctet[location%8] = value
	val := util.BinToByteValue(byteOctet)
	WriteBitmap(storage, fsdata, bitmap, val, location/8)
}

func cutContent(content string, index int, blockSize int) string {
	contentLen := len(content)
	lowerBound := index * blockSize
	upperBound := (index + 1) * blockSize
	if upperBound > contentLen {
		return content[lowerBound:contentLen]
	}
	return content[lowerBound:upperBound]
}

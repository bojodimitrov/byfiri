package diracts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bojodimitrov/byfiri/errors"
	"github.com/bojodimitrov/byfiri/structures"
)

// EncodeDirectoryContent returns byte array containing all files info
func EncodeDirectoryContent(files []structures.DirectoryEntry) (string, error) {
	var content strings.Builder
	for _, value := range files {
		fileNameSize := len(value.FileName)
		if fileNameSize == 0 {
			return "", fmt.Errorf("file name cannot be empty")
		}
		if fileNameSize > 255 {
			return "", fmt.Errorf("file name too long")
		}
		inodeStr := fmt.Sprint(value.Inode)
		content.WriteString(inodeStr)
		content.WriteString(strings.Repeat("\x00", 10-len(inodeStr)))
		fileNameSizeStr := fmt.Sprint(uint8(fileNameSize))
		content.WriteString(fileNameSizeStr)
		content.WriteString(strings.Repeat("\x00", 3-len(fileNameSizeStr)))
		content.WriteString(value.FileName)
	}
	return content.String(), nil
}

// DecodeDirectoryContent receives all blocks content concatenated and returns array of DirectoryContent
func DecodeDirectoryContent(content string) []structures.DirectoryEntry {
	var filesInfo []structures.DirectoryEntry
	offset := 0
	contentLen := len(content)
	inodeStr := strings.Trim(content[offset:offset+10], "\x00")
	for inodeStr != "" {
		offset += 10
		inode, err := strconv.Atoi(inodeStr)
		errors.CorruptData(err, "inode")
		fileNameSize, err := strconv.Atoi(strings.Trim(content[offset:offset+3], "\x00"))
		offset += 3
		errors.CorruptData(err, "file size")
		fileName := content[offset : offset+fileNameSize]
		filesInfo = append(filesInfo, structures.DirectoryEntry{
			Inode:    uint32(inode),
			FileName: fileName})
		offset += fileNameSize
		if contentLen == offset {
			break
		}
		inodeStr = strings.Trim(content[offset:offset+10], "\x00")
	}
	return filesInfo
}

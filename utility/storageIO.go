package utility

import (
	"fmt"
	"strings"
)

// Write writes in storage array
func Write(storage []byte, content string, offset int) {
	if len(content) == 0 {
		return
	}
	for i := 0; i < len(content); i++ {
		storage[i+offset] = content[i]
	}
}

// Read reads from storage array
func Read(storage []byte, offset int, length int) (string, error) {
	if length < 0 {
		return "", fmt.Errorf("length is negative")
	}
	return strings.Replace(string(storage[offset:length+offset]), "\x00", "", -1), nil
}

// ReadRaw reads from storage array and does not remove x00s
func ReadRaw(storage []byte, offset int, length int) (string, error) {
	if length < 0 {
		return "", fmt.Errorf("length is negative")
	}
	return string(storage[offset : length+offset]), nil
}

package utility

import "fmt"

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
	return string(storage[offset : length+offset]), nil
}

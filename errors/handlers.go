package errors

import (
	"fmt"
)

// CorruptMetadata prints error returned from reading the metadata and stops execution
func CorruptMetadata(err error) {
	if err != nil {
		fmt.Println(err)
		panic("corrupt metadata")
	}
}

// CorruptInode prints error returned from reading an inode and stops execution
func CorruptInode(err error, inode int) {
	if err != nil {
		fmt.Println(err)
		panic("corrupt inode: " + string(inode))
	}
}

// CorruptBitmap prints error returned from reading the metadata and stops execution
func CorruptBitmap(err error, bitmap string) {
	if err != nil {
		fmt.Println(err)
		panic("corrupt " + bitmap + " bitmap")
	}
}

// CorruptData is generic error handler
func CorruptData(err error, message string) {
	if err != nil {
		fmt.Println(err)
		panic("corrupt data: " + message)
	}
}

// IncorrectFormat prints error from which the system cannot recover
func IncorrectFormat(message string) {
	panic("Incorrect format error: " + message)
}

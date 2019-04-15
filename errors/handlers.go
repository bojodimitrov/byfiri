package errors

import (
	"fmt"
	"os"
)

// CorruptMetadata prints error returned from reading the metadata and stops execution
func CorruptMetadata(err error) {
	if err != nil {
		fmt.Println("corrupt metadata")
		fmt.Println(err)
		os.Exit(3)
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
		fmt.Println("corrupt " + bitmap + " bitmap")
		fmt.Println(err)
		os.Exit(3)
	}
}

// CorruptData is generic error handler
func CorruptData(err error, message string) {
	if err != nil {
		fmt.Println("corrupt data: " + message)
		fmt.Println(err)
		os.Exit(3)
	}
}

// IncorrectFormat prints error from which the system cannot recover
func IncorrectFormat(message string) {
	fmt.Println("Incorrect format error: " + message)
	os.Exit(3)
}

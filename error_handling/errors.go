package error_handling

import (
	"fmt"
	"os"
)

// CorruptMetadata prints error returned from reading the metadata and stops execution
func CorruptMetadata(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Println("corrupt metadata")
		os.Exit(3)
	}
}

// CorruptBitmap prints error returned from reading the metadata and stops execution
func CorruptBitmap(err error, bitmap string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println("corrupt " + bitmap + " bitmap")
		os.Exit(3)
	}
}

// IncorrectFormat prints error from which the system cannot recover
func IncorrectFormat(message string) {
	fmt.Println("Incorrect format error: " + message)
	os.Exit(3)
}

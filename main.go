package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	size := getSize()
	storage := initFsSpace(size)
	sizeString := strconv.Itoa(size)
	write(storage, sizeString, 0)
	fmt.Println(read(storage, 0, -3))
	fmt.Println(size)
}

func initFsSpace(size int) []byte {
	storage := make([]byte, size)
	return storage
}

func getSize() int {
	requestedSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("incorrect size format:", os.Args[1])
		os.Exit(3)
	}
	sizeAbbreviation := os.Args[2]

	switch sizeAbbreviation {
	case "GB":
		requestedSize = requestedSize * 1073741824
	case "MB":
		requestedSize = requestedSize * 1048576
	}
	return requestedSize
}

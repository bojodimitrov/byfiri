package utility

import (
	"fmt"
	"strconv"
)

// ByteToBin transforms byte array to binary array
func ByteToBin(b []byte) []bool {
	var binary []bool
	for _, value := range b {
		binaryString := strconv.FormatInt(int64(value), 2)
		for _, rune := range binaryString {
			value, err := strconv.ParseBool(string(rune))
			if err != nil {
				fmt.Println("Value missed, incorrect bool")
			} else {
				binary = append(binary, value)
			}
		}
	}
	return binary
}

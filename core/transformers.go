package core

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// ByteToBin transforms byte array to binary array
func ByteToBin(b []byte) []bool {
	var binary []bool
	for _, value := range b {
		binaryString := strconv.FormatInt(int64(value), 2)
		binaryString = strings.Repeat("0", 8-len(binaryString)) + binaryString
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

// BinToByteValue transforms binary array to byte value
func BinToByteValue(b []bool) byte {
	var binBuf bytes.Buffer
	bitSet := int64(0)
	for _, v := range b {
		bitSet = 0
		if v {
			bitSet = 1
		}
		binBuf.WriteString(strconv.FormatInt(bitSet, 10))
	}
	binString := binBuf.String()
	val, _ := strconv.ParseInt(binString, 2, 64)
	return byte(val)
}

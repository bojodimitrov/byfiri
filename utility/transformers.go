package utility

import (
	"fmt"
	"strconv"
)

// StringToBin transforms hex string to binary array
func StringToBin(s string) []bool {
	binaryString := ""
	for _, c := range s {
		binaryString = fmt.Sprintf("%s%.8b", binaryString, c)
	}
	var binary []bool
	for _, rune := range binaryString {
		value, err := strconv.ParseBool(string(rune))
		if err != nil {
			fmt.Println("Value missed, incorrect bool")
		} else {
			binary = append(binary, value)
		}
	}
	return binary
}

package utility

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

func GetNoOfBitsAndBytes(value int64) (int, int) {
	// Example big integer
	bigNum := new(big.Int)
	bigNum.SetString(fmt.Sprintf("%d", value), 10)

	// Calculate the number of bits needed
	bitLen := bigNum.BitLen()

	// Calculate the number of bytes needed
	byteLen := (bitLen + 7) / 8
	return bitLen, byteLen
}

func FormatAsBinary(value, noOfBytes int) string {
	format := fmt.Sprintf("%%0%db", noOfBytes*8)
	result := fmt.Sprintf(format, value)

	re := regexp.MustCompile(".{1,8}")
	result = re.ReplaceAllStringFunc(result, func(s string) string {
		return s + " "
	})

	return fmt.Sprintf("[%s]", strings.TrimSpace(result))
}

func AbsInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

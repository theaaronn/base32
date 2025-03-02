package base32

import (
	"errors"
	"math"
	"strings"
)

var paddingEqualsMapEnc = map[int]int{
	1: 6,
	2: 4,
	3: 3,
	4: 1,
	5: 0,
}
var paddingBitsMapEnc = map[int]int{
	1: 2,
	2: 4,
	3: 1,
	4: 3,
	5: 0,
}
var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
var alphabetMap = map[rune]int{
	'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7,
	'I': 8, 'J': 9, 'K': 10, 'L': 11, 'M': 12, 'N': 13, 'O': 14, 'P': 15,
	'Q': 16, 'R': 17, 'S': 18, 'T': 19, 'U': 20, 'V': 21, 'W': 22, 'X': 23,
	'Y': 24, 'Z': 25, '2': 26, '3': 27, '4': 28, '5': 29, '6': 30, '7': 31,
}

func Encode32(initialString string) (finalString string, err error) {
	length := len(initialString)
	if length == 0 {
		return "", nil
	}

	var (
		inputBytes = []byte(initialString)
		start      = 0
		end        = 0
		buf        = 0
	)
	for end != length {
		start = end
		end = min(end+5, length)
		chunk := inputBytes[start:end]

		for _, b := range chunk {
			buf = (buf << 8) + int(b)
		}
		buf <<= paddingBitsMapEnc[len(chunk)]
		numChars := int(math.Ceil(float64(len(chunk)*8) / 5))
		for i := numChars - 1; i >= 0; i-- {
			val := (buf >> (i * 5)) & 0x1F
			finalString += string(alphabet[val])
		}
		// Pad equals
		for range paddingEqualsMapEnc[len(chunk)] {
			finalString += "="
		}
		buf = 0
	}
	return
}
func Decode32(encodedStr string) (decodedStr string, err error) {
	length := len(encodedStr)
	if length == 0 {
		return "", nil
	}

	encodedStr = strings.TrimRight(encodedStr, "=")

	var (
		buf      int
		bitsLeft int
		output   []byte
	)
	for _, char := range encodedStr {
		if val, ok := alphabetMap[char]; ok {
			buf = (buf << 5) | val
			bitsLeft += 5
			for bitsLeft >= 8 {
				bitsLeft -= 8
				b := (buf >> bitsLeft) & 0xFF
				output = append(output, byte(b))
			}
		} else {
			return "", errors.New("invalid character in input")
		}
	}
	return string(output), nil
}

package iabtcf

import (
	"fmt"
	"strings"
)

// //////////////////////////////////////////////////
// bits

type Bits struct {
	Length int
	Bytes  []byte
}

func (bits *Bits) HasBit(number int) bool {
	if number < 1 || number > bits.Length {
		// out of scope
		return false
	}
	byteIndex := (number - 1) / 8
	b := bits.Bytes[byteIndex]
	bitIndex := (number - 1) % 8
	mask := byte(1 << (7 - bitIndex))
	// fmt.Printf("[HasBit] number: %d / byteIndex: %d / b: %d / bitIndex:%d / mask:%d / &:%d > %t", number, byteIndex, b, bitIndex, mask, b&mask, b&mask == mask)
	return b&mask == mask
}

func (bits *Bits) ToBitString() string {
	if bits == nil {
		return ""
	}

	result := ""

	for i, b := range bits.Bytes {
		if i != 0 {
			result += " "
		}
		result += fmt.Sprintf("%08b", b)
	}

	trim := 8*len(bits.Bytes) - bits.Length
	if trim > 0 {
		result = result[:len(result)-trim]
	}

	return result
}

// //////////////////////////////////////////////////
// bit string helper

func BitStringToBits(value string) *Bits {
	return &Bits{
		Length: len(strings.ReplaceAll(string(value), " ", "")),
		Bytes:  BitStringToBytes(value),
	}
}

func BitStringToBytes(value string) []byte {
	bytes := make([]byte, 0, len(value)/8)

	position := 7
	var lastByte byte
	for i := 0; i < len(value); i++ {
		if value[i] == ' ' {
			continue
		}
		if value[i] == '1' {
			lastByte |= 1 << position
		}
		if position == 0 {
			position = 7
			bytes = append(bytes, lastByte)
			lastByte = 0
		} else {
			position--
		}
	}
	if position != 7 {
		bytes = append(bytes, lastByte)
	}
	return bytes
}

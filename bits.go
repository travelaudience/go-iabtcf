package iabtcf

import (
	"fmt"
	"strings"
	"time"
)

// //////////////////////////////////////////////////
// bits

// Bits represents a bitset with some helpers to read int, bool, string and time fields
type Bits struct {
	Length int
	Bytes  []byte
}

// HasBit checks if the bit number is set
func (b *Bits) HasBit(number int) bool {
	value, _ := b.ReadBoolField(number - 1)
	return value
}

// ReadInt64Field reads an int64 field of nbBits bits starting at offset
func (b *Bits) ReadInt64Field(offset, nbBits int) (int64, error) {
	if err := b.checkBounds(offset, nbBits); err != nil {
		return 0, err
	}
	var result int64
	byteIndex := offset / 8
	bitIndex := offset % 8
	for i := 0; i < nbBits; i++ {
		mask := byte(1 << (7 - bitIndex))
		if b.Bytes[byteIndex]&mask == mask {
			result |= 1 << (nbBits - 1 - i)
		}
		if bitIndex == 7 {
			byteIndex++
			bitIndex = 0
		} else {
			bitIndex++
		}
	}
	return result, nil
}

// ReadIntField reads an int field of nbBits bits starting at offset
func (b *Bits) ReadIntField(offset, nbBits int) (int, error) {
	value, err := b.ReadInt64Field(offset, nbBits)
	if err != nil {
		return 0, err
	}
	return int(value), nil
}

const (
	TimeNbBits = 36
)

// ReadTimeField reads a time field of 36 bits starting at offset
func (b *Bits) ReadTimeField(offset int) (time.Time, error) {
	ds, err := b.ReadInt64Field(offset, TimeNbBits)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(ds/dsPerSec, (ds%dsPerSec)*nsPerDs).UTC(), nil
}

const (
	CharacterNbBits = 6
)

// ReadStringField reads a string field of nbBits bits starting at offset
//
// note: each character is represented by 6 bits, so the number of bits must be a multiple of 6
// note: the characters are represented by the uppercase alphabet starting from 'A'
func (b *Bits) ReadStringField(offset, nbBits int) (string, error) {
	length := nbBits / CharacterNbBits
	if nbBits%CharacterNbBits != 0 {
		return "", ErrInvalidNbBitsMultiple(nbBits, CharacterNbBits)
	}
	var buf = make([]byte, 0, length)
	nextOffset := offset
	for i := 0; i < length; i++ {
		value, err := b.ReadInt64Field(nextOffset, CharacterNbBits)
		if err != nil {
			return "", err
		}
		buf = append(buf, byte(value)+'A')
		nextOffset += CharacterNbBits
	}
	return string(buf), nil
}

const (
	BoolNbBits = 1
)

// ReadBoolField reads a bool field of 1 bit starting at offset
func (b *Bits) ReadBoolField(offset int) (bool, error) {
	value, err := b.ReadInt64Field(offset, BoolNbBits)
	if err != nil {
		return false, err
	}
	return value == 1, nil
}

func (b *Bits) checkBounds(offset, nbBits int) error {
	if b == nil {
		return ErrNilBits
	}
	if offset < 0 {
		return ErrNegativeBitIndex(offset)
	}
	if offset+nbBits > b.Length {
		return ErrBitIndexHigherThanUpperBound(offset+nbBits, b.Length)
	}
	return nil
}

// ToBitString returns the bitset as a string of bits ( human readable 0s and 1s )
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

// BitStringToBits converts a bit string to a Bits struct
func BitStringToBits(value string) Bits {
	return Bits{
		Length: len(strings.ReplaceAll(string(value), " ", "")),
		Bytes:  BitStringToBytes(value),
	}
}

// BitStringToBytes converts a bit string to a byte slice
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

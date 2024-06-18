package iabtcf

import (
	"fmt"
	"time"
)

// //////////////////////////////////////////////////
// bits

// Bits represents a bitset with some helpers to read int, bool, string and time fields
//
// Bits are stored in a byte slice
// First byte will store the first 8 bits, second byte the next 8 bits, and so on
//
// note: the last byte may contain less than 8 bits. Those bits are left aligned.
type Bits []byte

// HasBit checks if the bit number is set
//
// note: number is not the index and it starts at 1.
func (b Bits) HasBit(number int) bool {
	value, _ := b.ReadBoolField(number - 1)
	return value
}

const (
	nbBitInByte  = 8
	lastBitIndex = nbBitInByte - 1
)

var (
	bitMasks = [nbBitInByte]byte{
		1 << 7,
		1 << 6,
		1 << 5,
		1 << 4,
		1 << 3,
		1 << 2,
		1 << 1,
		1,
	}
)

// ReadInt64Field reads an int64 field of nbBits bits starting at offset
func (b Bits) ReadInt64Field(offset, nbBits int) (int64, error) {
	if err := b.checkBounds(offset, nbBits); err != nil {
		return 0, err
	}
	var result int64
	byteIndex := offset / nbBitInByte
	bitIndex := offset % nbBitInByte
	for i := 0; i < nbBits; i++ {
		mask := bitMasks[bitIndex]
		if b[byteIndex]&mask == mask {
			result |= 1 << (nbBits - 1 - i)
		}
		if bitIndex == lastBitIndex {
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
		return "", fmt.Errorf("number of bits %d is not multiple of %d bits", nbBits, CharacterNbBits)
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

func (b Bits) checkBounds(offset, nbBits int) error {
	if b == nil {
		return fmt.Errorf("bits is nil")
	}
	if offset < 0 {
		return fmt.Errorf("negative bit index: %d", offset)
	}
	if offset+nbBits > len(b)*nbBitInByte {
		return fmt.Errorf("bit index %d is higher than upper bound %d", offset+nbBits, len(b)*nbBitInByte)
	}
	return nil
}

// ToBitString returns the bitset as a string of bits ( human readable 0s and 1s )
func (bits Bits) ToBitString() string {
	if bits == nil {
		return ""
	}

	result := ""

	for i, b := range bits {
		if i != 0 {
			result += " "
		}
		result += fmt.Sprintf("%08b", b)
	}

	return result
}

// //////////////////////////////////////////////////
// bit string helper

// BitStringToBits converts a bit string to a Bits struct
func BitStringToBits(value string) Bits {
	return Bits(BitStringToBytes(value))
}

// BitStringToBytes converts a bit string to a byte slice
func BitStringToBytes(value string) []byte {
	bytes := make([]byte, 0, len(value)/nbBitInByte)

	position := lastBitIndex
	var lastByte byte
	for i := 0; i < len(value); i++ {
		if value[i] == ' ' {
			continue
		}
		if value[i] == '1' {
			lastByte |= 1 << position
		}
		if position == 0 {
			position = lastBitIndex
			bytes = append(bytes, lastByte)
			lastByte = 0
		} else {
			position--
		}
	}
	if position != nbBitInByte-1 {
		bytes = append(bytes, lastByte)
	}
	return bytes
}

package iabtcf

import (
	"fmt"
	"time"

	"github.com/rupertchen/go-bits"
)

const (
	// deciseconds per second
	dsPerSec = 10
	// nanoseconds per decisecond
	nsPerDs = int64(time.Millisecond * 100)
)

type Reader struct {
	*bits.Reader
}

// NewReader returns a new Reader
func NewReader(src []byte) *Reader {
	return &Reader{bits.NewReader(bits.NewBitmap(src))}
}

// ReadInt reads the next n bits and converts them to an int.
func (r *Reader) ReadInt(n uint) (int, error) {
	b, err := r.ReadBits(n)
	if err != nil {
		return 0, fmt.Errorf("ReadBits failed: %s", err.Error())
	}

	return int(b), nil
}

// ReadTime reads the next 36 bits representing timestamp in deciseconds
// and converts it to time.Time
func (r *Reader) ReadTime() (time.Time, error) {
	b, err := r.ReadBits(36)
	if err != nil {
		return time.Time{}, fmt.Errorf("ReadBits failed: %s", err.Error())
	}
	ds := int64(b)
	return time.Unix(ds/dsPerSec, (ds%dsPerSec)*nsPerDs).UTC(), nil
}

// ReadString reads a string of length n
func (r *Reader) ReadString(length int) (string, error) {
	length = length / 6 // 6 bit per letter
	var buf = make([]byte, 0, length)
	for i := 0; i < length; i++ {
		if b, err := r.ReadBits(6); err != nil {
			return "", fmt.Errorf("ReadBits failed: %s", err.Error())
		} else {
			buf = append(buf, byte(b)+'A')
		}
	}
	return string(buf), nil
}

// ReadBitField reads the next n bits into a bit map
func (r *Reader) ReadBitField(length int) (Bits, error) {
	remaining := length % 8
	nb := (length - remaining) / 8
	bytes := make([]byte, 0, nb+1)
	for i := 0; i < nb; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return Bits{}, fmt.Errorf("ReadByte failed: %s", err.Error())
		}
		bytes = append(bytes, b)
	}
	if remaining > 0 {
		block, err := r.ReadBits(uint(remaining))
		if err != nil {
			return Bits{}, fmt.Errorf("ReadBits failed: %s", err.Error())
		}
		// note: bits are right aligned
		b := byte(block << (8 - remaining))
		bytes = append(bytes, b)
	}
	return Bits{Bytes: bytes, Length: length}, nil
}

// ReadRangeEntries reads a list of range entries
func (r *Reader) ReadRangeEntries(length int) ([]RangeEntry, error) {
	res := make([]RangeEntry, 0, length)
	var err error
	for i := 0; i < length; i++ {
		var isRange bool
		if isRange, err = r.ReadBool(); err != nil {
			return nil, fmt.Errorf("ReadBool failed: %s", err.Error())
		}
		var start, end int
		if start, err = r.ReadInt(16); err != nil {
			return nil, fmt.Errorf("ReadInt failed: %s", err.Error())
		}
		if isRange {
			if end, err = r.ReadInt(16); err != nil {
				return nil, fmt.Errorf("ReadInt failed: %s", err.Error())
			}
		} else {
			end = start
		}
		res = append(res, RangeEntry{StartOrOnlyVendorId: start, EndVendorID: end})
	}
	return res, nil
}

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
func (r *Reader) ReadString(n uint) (string, error) {
	n = n / 6 // 6 bit per letter
	var buf = make([]byte, 0, n)
	for i := uint(0); i < n; i++ {
		if b, err := r.ReadBits(6); err != nil {
			return "", fmt.Errorf("ReadBits failed: %s", err.Error())
		} else {
			buf = append(buf, byte(b)+'A')
		}
	}
	return string(buf), nil
}

// ReadBitField reads the next n bits into a map[int]bool
func (r *Reader) ReadBitField(n uint) (map[int]bool, error) {
	var m = make(map[int]bool)
	for i := uint(0); i < n; i++ {
		b, err := r.ReadBool()
		if err != nil {
			return nil, fmt.Errorf("ReadBool failed: %s", err.Error())
		}
		m[int(i)+1] = b
	}
	return m, nil
}

func (r *Reader) ReadRangeEntries(n uint) ([]*RangeEntry, error) {
	res := make([]*RangeEntry, 0, n)
	var err error
	for i := uint(0); i < n; i++ {
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
		res = append(res, &RangeEntry{StartOrOnlyVendorId: start, EndVendorID: end})
	}
	return res, nil
}

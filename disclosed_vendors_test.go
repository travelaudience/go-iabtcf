package iabtcf

import (
	"encoding/base64"
	"strconv"
	"strings"
	"testing"
)

func TestDisclosedVendors(t *testing.T) {

	tests := []struct {
		block                    string
		hasDisclosedVendorsBlock bool
		vendorID                 int
		hasVendorID              bool
	}{
		{ // case 0
			block:                    "0",
			hasDisclosedVendorsBlock: false,
			vendorID:                 123,
			hasVendorID:              false,
		},
		{ // case 1
			block:                    sprintb(1, 3) + sprintb(0, 16),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              false,
		},
		{ // case 2
			block:                    sprintb(1, 3) + sprintb(256, 16),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              false,
		},
		{ // case 3
			block:                    sprintb(1, 3) + sprintb(256, 16) + "0",
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              false,
		},
		{ // case 4
			block:                    sprintb(1, 3) + sprintb(256, 16) + "0" + strings.Repeat("0", 123),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              false,
		},
		{ // case 5
			block:                    sprintb(1, 3) + sprintb(256, 16) + "0" + strings.Repeat("0", 122) + "1",
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              true,
		},
		{ // case 6
			block:                    sprintb(1, 3) + sprintb(256, 16) + "1" + sprintb(0, 12),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              false,
		},
		{ // case 7
			block:                    sprintb(1, 3) + sprintb(256, 16) + "1" + sprintb(1, 12) + "0" + sprintb(123, 16),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              true,
		},
		{ // case 8
			block:                    sprintb(1, 3) + sprintb(256, 16) + "1" + sprintb(1, 12) + "0" + sprintb(124, 16),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              false,
		},
		{ // case 9
			block:                    sprintb(1, 3) + sprintb(256, 16) + "1" + sprintb(1, 12) + "1" + sprintb(120, 16) + sprintb(127, 16),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              true,
		},
		{ // case 10
			block:                    sprintb(1, 3) + sprintb(256, 16) + "1" + sprintb(1, 12) + "1" + sprintb(130, 16) + sprintb(137, 16),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              false,
		},
		{ // case 11
			block:                    sprintb(1, 3) + sprintb(256, 16) + "1" + sprintb(2, 12) + "1" + sprintb(130, 16) + sprintb(137, 16) + "0" + sprintb(123, 16),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              true,
		},
		{ // case 12
			block:                    sprintb(1, 3) + sprintb(256, 16) + "1" + sprintb(2, 12) + "0" + sprintb(124, 16) + "1" + sprintb(120, 16) + sprintb(123, 16),
			hasDisclosedVendorsBlock: true,
			vendorID:                 123,
			hasVendorID:              true,
		},
	}

	dummyCoreString := base64.RawURLEncoding.EncodeToString(sscanb("000010" + strings.Repeat("0", 224)))

	for i, tt := range tests {
		consent := dummyCoreString + "." + base64.RawURLEncoding.EncodeToString(sscanb(tt.block))
		lc, err := LazyParseCoreString(consent)
		if err != nil {
			t.Error(err)
			continue
		}
		if want, got := tt.hasDisclosedVendorsBlock, lc.HasDisclosedVendorsBlock(); want != got {
			t.Errorf("case %d: has disclosed vendor block: want %v, got %v", i, want, got)
		}
		if want, got := tt.hasVendorID, lc.IsVendorDisclosed(tt.vendorID); want != got {
			t.Errorf("case %d: includes vendor %d: want %v, got %v", i, tt.vendorID, want, got)
		}
	}
}

// sprintb returns a binary string representation of the number given.  The string is guarateed to be exactly bits in
// length.  In case of overflow, high order bits are truncated.  Bits must be ≤ 63.
func sprintb(number, bits int) string {
	x := strconv.FormatInt(int64(number), 2)
	switch {
	case len(x) > bits:
		return x[len(x)-bits:]
	case len(x) == bits:
		return x
	default:
		return strings.Repeat("0", bits-len(x)) + x
	}
}

func TestSprintb(t *testing.T) {
	tests := []struct {
		number int
		bits   int
		want   string
	}{
		{5, 4, "0101"},       // 5 in 4 bits: "0101"
		{5, 3, "101"},        // 5 in 3 bits: "101"
		{5, 2, "01"},         // 5 in 2 bits: "01" (truncated)
		{0, 4, "0000"},       // 0 in 4 bits: "0000"
		{15, 4, "1111"},      // 15 in 4 bits: "1111"
		{16, 4, "0000"},      // 16 in 4 bits: "0000" (truncated)
		{1, 1, "1"},          // 1 in 1 bit: "1"
		{2, 4, "0010"},       // 2 in 4 bits: "0010"
		{255, 8, "11111111"}, // 255 in 8 bits: "11111111"
		{256, 8, "00000000"}, // 256 in 8 bits: "00000000" (truncated)
	}
	for _, tt := range tests {
		got := sprintb(tt.number, tt.bits)
		if got != tt.want {
			t.Errorf("sprintb(%d, %d) = '%s', want '%s'", tt.number, tt.bits, got, tt.want)
		}
	}
}

// sscanb converts a string of 1s and 0s to a big-endian byte array.  That is the first digit in the string is
// mapped to bit 7 of the byte in position 0.  Right zero-padding is added as needed.
func sscanb(binary string) []byte {

	n := len(binary)
	if n%8 != 0 {
		binary += strings.Repeat("0", 8-n%8)
		n = len(binary)
	}
	out := make([]byte, n/8)
	for i := 0; i < n; i += 8 {
		var b byte
		for j := range 8 {
			if binary[i+j] == '1' {
				b |= 1 << (7 - j)
			}
		}
		out[i/8] = b
	}
	return out
}

func TestSscanb(t *testing.T) {
	tests := []struct {
		input string
		want  []byte
	}{
		{"", []byte{}},
		{"1", []byte{128}},
		{"0001", []byte{16}},
		{"10101", []byte{168}},
		{"11111111", []byte{255}},
		{"00000000", []byte{0}},
		{"1010101010101", []byte{170, 168}},
		{"1010101010101011", []byte{170, 171}},
		{"10101010101010110", []byte{170, 171, 0}},
		{"10101010101010111", []byte{170, 171, 128}},
	}
	for _, tt := range tests {
		got := sscanb(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("sscanb(%q) length = %d, want %d", tt.input, len(got), len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("sscanb(%q)[%d] = %d, want %d", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

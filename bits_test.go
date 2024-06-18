package iabtcf

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBits(t *testing.T) {

	wantHasBit := func(number int, expected []int) bool {
		for _, e := range expected {
			if number == e {
				return true
			}
		}
		return false
	}

	type TestCase struct {
		Base64        string
		WantBitString string
		WantHasBit    []int
	}

	values := map[string]*TestCase{
		"101": {
			Base64:        "oA",
			WantBitString: "10100000",
			WantHasBit:    []int{1, 3},
		},
		"00000001": {
			Base64:        "AQ",
			WantBitString: "00000001",
			WantHasBit:    []int{8},
		},
		"00000101": {
			Base64:        "BQ",
			WantBitString: "00000101",
			WantHasBit:    []int{6, 8},
		},
		"10000101": {
			Base64:        "hQ",
			WantBitString: "10000101",
			WantHasBit:    []int{1, 6, 8},
		},
		"00000001 00000101": {
			Base64:        "AQU",
			WantBitString: "00000001 00000101",
			WantHasBit:    []int{8, 14, 16},
		},
		"00000001 101": {
			Base64:        "AaA",
			WantBitString: "00000001 10100000",
			WantHasBit:    []int{8, 9, 11},
		},
		"00000001 00000000": {
			Base64:        "AQA",
			WantBitString: "00000001 00000000",
			WantHasBit:    []int{8},
		},
		"00000001 00000000 1": {
			Base64:        "AQCA",
			WantBitString: "00000001 00000000 10000000",
			WantHasBit:    []int{8, 17},
		},
		"00000001 0000001": {
			Base64:        "AQI",
			WantBitString: "00000001 00000010",
			WantHasBit:    []int{8, 15},
		},
	}

	for bitString, tc := range values {
		t.Run(bitString, func(t *testing.T) {
			t.Helper()
			fmt.Printf("\n[test] ---------- %s ---------- \n", bitString)
			var wantBytes, err = base64.RawURLEncoding.DecodeString(tc.Base64)
			fmt.Printf("[test] base64: %s >>> bytes: %v \n", tc.Base64, wantBytes)
			require.NoError(t, err, "unexpected base64 error")

			gotBits := BitStringToBits(bitString)
			fmt.Printf("[test] bits: %s >>> bytes: %v \n", bitString, gotBits)
			require.Equal(t, wantBytes, []byte(gotBits))

			fmt.Printf("[test] Bits: %v \n", gotBits)

			fmt.Printf("[test] bytes: %v >>> bits: %s \n", gotBits, gotBits.ToBitString())
			require.Equal(t, tc.WantBitString, gotBits.ToBitString())

			fmt.Printf("[test] WantHasBit: %v \n", tc.WantHasBit)
			length := len(strings.ReplaceAll(bitString, " ", ""))
			for number := 1; number <= length; number++ {
				gotHasBit := gotBits.HasBit(number)
				wantHasBit := wantHasBit(number, tc.WantHasBit)
				require.Equal(t, wantHasBit, gotHasBit)
			}
		})
	}
}

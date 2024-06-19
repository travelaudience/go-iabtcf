package iabtcf

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {

	type TestCase struct {
		Base64 string
	}

	values := map[string]*TestCase{
		"101": {
			Base64: "oA",
		},
		"00000001": {
			Base64: "AQ",
		},
		"00000101": {
			Base64: "BQ",
		},
		"10000101": {
			Base64: "hQ",
		},
		"00000001 00000101": {
			Base64: "AQU",
		},
		"00000001 101": {
			Base64: "AaA",
		},
		"00000001 00000000": {
			Base64: "AQA",
		},
		"00000001 00000000 1": {
			Base64: "AQCA",
		},
		"00000001 0000001": {
			Base64: "AQI",
		},
	}

	for bitString, tc := range values {
		t.Run(bitString, func(t *testing.T) {
			t.Helper()
			fmt.Printf("\n[test] ---------- %s ---------- \n", bitString)
			wantBytes, err := base64.RawURLEncoding.DecodeString(tc.Base64)
			fmt.Printf("[test] base64: %s >>> bytes: %v \n", tc.Base64, wantBytes)
			require.NoError(t, err, "unexpected base64 error")

			length := len(strings.ReplaceAll(bitString, " ", ""))
			reader := NewReader(wantBytes)
			gotBits, err := reader.ReadBitField(length)
			require.NoError(t, err, "unexpected reader error")
			fmt.Printf("[test] %s (%d) >>> bytes: %v >>> bitfield: %s \n", bitString, length, gotBits, gotBits.ToBitString())
		})
	}
}

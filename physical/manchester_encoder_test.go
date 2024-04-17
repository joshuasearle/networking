package physical_test

import (
	"networking/physical"
	"testing"
)

func BitStringToBitArray(bitString string) []physical.Bit {
	bits := make([]physical.Bit, len(bitString))
	for i, bit := range bitString {
		if bit == '0' {
			bits[i] = physical.Zero
		} else {
			bits[i] = physical.One
		}
	}
	return bits
}

var testCases = []struct {
	data             string
	expectedEncoding string
}{
	{
		"",
		"",
	},
	{
		"0",
		"10",
	},
	{
		"1",
		"01",
	},
	{
		"00",
		"1010",
	},
	{
		"01",
		"1001",
	},
	{
		"10",
		"0110",
	},
	{
		"11",
		"0101",
	},
}

func TestManchesterEncoder(t *testing.T) {
	en := physical.NewManchesterEncoder()

	for _, tc := range testCases {
		data := BitStringToBitArray(tc.data)
		expectedEncoding := BitStringToBitArray(tc.expectedEncoding)

		actualEncoding := en.Encode(data)

		if len(actualEncoding) != len(expectedEncoding) {
			t.Errorf("Expected %v, got %v", expectedEncoding, actualEncoding)
			continue
		}

		for i := 0; i < len(actualEncoding); i++ {
			if actualEncoding[i] != expectedEncoding[i] {
				t.Errorf("Expected %v, got %v", expectedEncoding, actualEncoding)
				break
			}
		}
	}
}

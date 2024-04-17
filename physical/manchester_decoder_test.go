package physical_test

import (
	"networking/physical"
	"testing"
)

func BitStringToBitArray2(bitString string) []physical.Bit {
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

var testCases2 = []struct {
	encodedBits         string
	expectedDecodedBits string
}{
	{
		"",
		"",
	},
	{
		"0",
		"",
	},
	{
		"1",
		"",
	},
	{
		"00",
		"",
	},
	{
		"01",
		"1",
	},
	{
		"10",
		"0",
	},
	{
		"11",
		"",
	},
	{
		"000",
		"",
	},
	{
		"001",
		"1",
	},
	{
		"010",
		"10",
	},
	{
		"011",
		"1",
	},
	{
		"100",
		"0",
	},
	{
		"101",
		"01",
	},
	{
		"110",
		"0",
	},
	{
		"111",
		"",
	},
	{
		"0000",
		"",
	},
	{
		"0001",
		"1",
	},
	{
		"0010",
		"10",
	},
	{
		"0011",
		"1",
	},
	{
		"0100",
		"10",
	},
	{
		"0101",
		"101",
	},
	{
		"0110",
		"10",
	},
	{
		"0111",
		"1",
	},
	{
		"1000",
		"0",
	},
	{
		"1001",
		"01",
	},
	{
		"1010",
		"010",
	},
	{
		"1011",
		"01",
	},
	{
		"1100",
		"0",
	},
	{
		"1101",
		"01",
	},
	{
		"1110",
		"0",
	},
	{
		"1111",
		"",
	},
	// 8 bit preamble
	{
		"0110011001100110",
		"10101010",
	},
	// 8 bit preamble with extra 10
	{
		"011001100110011010",
		"101010100",
	},
	// 7 bit preamble
	{
		"01100110011001",
		"1010101",
	},
	// 7 bit preamble with extra 0 (should go back to preamble)
	{
		"011001100110010",
		"10101010",
	},
}

type Handler struct {
	bits []physical.Bit
}

func NewHandler() *Handler {
	return &Handler{
		bits: make([]physical.Bit, 0),
	}
}

func (h *Handler) Handle(bit physical.Bit) {
	h.bits = append(h.bits, bit)
}

func TestManchesterDecodings(t *testing.T) {
	for _, tc := range testCases2 {
		h := NewHandler()
		d := physical.NewManchesterDecoder(h)

		encodedBits := BitStringToBitArray2(tc.encodedBits)
		expectedDecodedBits := BitStringToBitArray2(tc.expectedDecodedBits)

		for _, bit := range encodedBits {
			d.Handle(bit)
		}

		actualDecodedBits := h.bits

		if len(actualDecodedBits) != len(expectedDecodedBits) {
			t.Errorf("For %v, expected %v, got %v", tc.encodedBits, expectedDecodedBits, actualDecodedBits)
			continue
		}

		for i := 0; i < len(actualDecodedBits); i++ {
			if actualDecodedBits[i] != expectedDecodedBits[i] {
				t.Errorf("For %v, expected %v, got %v", tc.encodedBits, expectedDecodedBits, actualDecodedBits)
			}
			break
		}
	}
}

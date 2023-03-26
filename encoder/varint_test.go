package encoder

import (
	"bytes"
	"math"
	"testing"
)

var testCases = []struct {
	Uint64  uint64
	Uvarint []byte
}{
	{0, []byte{0}},
	{1, []byte{1}},
	{0b0111_1111, []byte{0b0111_1111}},
	{0b1111_1111, []byte{0b1111_1111, 0b0000_0001}},
	{0b1000_0000, []byte{0b1000_0000, 0b0000_0001}},
	{math.MaxUint64, []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 1}},
}

func TestWriteVarInt(t *testing.T) {
	for _, tc := range testCases {
		output := writeUvarint(tc.Uint64)
		if !bytes.Equal(output, tc.Uvarint) {
			t.Errorf("writeUvarint(%d)=%v, want=%v", tc.Uint64, output, tc.Uvarint)
		}
	}
}

func TestReadVarInt(t *testing.T) {
	for _, tc := range testCases {
		output, _ := readUvarint(tc.Uvarint)
		if output != tc.Uint64 {
			t.Errorf("readUvarint(%v)=%d, want=%d", tc.Uvarint, output, tc.Uint64)
		}
	}
}

var testCasesSkip2Bits = []struct {
	Uint64  uint64
	Uvarint []byte
}{
	{0, []byte{0}},
	{1, []byte{1}},
	{0b0001_1111, []byte{0b0001_1111}},
	{0b0011_1111, []byte{0b0011_1111, 0b0000_0001}},
	{math.MaxUint64, []byte{63, 255, 255, 255, 255, 255, 255, 255, 255, 7}},
}

func TestWriteVarIntSkip2Bits(t *testing.T) {
	for _, tc := range testCasesSkip2Bits {
		output := writeUvarintSkip2Bits(tc.Uint64)
		if !bytes.Equal(output, tc.Uvarint) {
			t.Errorf("writeUvarint(%d)=%v, want=%v", tc.Uint64, output, tc.Uvarint)
		}
	}
}

func TestReadVarIntSkip2Bits(t *testing.T) {
	for _, tc := range testCasesSkip2Bits {
		output, bytesRead := readUvarintSkip2Bits(tc.Uvarint)
		if output != tc.Uint64 {
			t.Errorf("readUvarint(%v)=%d, want=%d", tc.Uvarint, output, tc.Uint64)
		}
		if bytesRead != len(tc.Uvarint) {
			t.Errorf("readVarint(%v)=_,%d, want %d", output, bytesRead, len(tc.Uvarint))
		}
	}
}

var testCasesSigned = []int64{
	math.MinInt64,
	math.MinInt64 + 1,
	-257,
	-256,
	-255,
	-1,
	0,
	1,
	255,
	256,
	257,
}

func TestSignedVarint(t *testing.T) {
	for _, input := range testCasesSigned {
		output := writeVarint(input)
		backAgain, bytesRead := readVarint(output)
		if input != backAgain {
			t.Errorf("writeVarint(%[1]d)=%[2]v\readVarint(%[2]v)=%[1]d. expected to match", input, output)
		}
		if bytesRead != len(output) {
			t.Errorf("readVarint(%v)=_,%d, want %d", output, bytesRead, len(output))
		}
	}
}

func TestSignedVarintSkip2Bits(t *testing.T) {
	for _, input := range testCasesSigned {
		output := writeVarintSkip2Bits(input)
		backAgain, bytesRead := readVarintSkip2Bits(output)
		if input != backAgain {
			t.Errorf("writeVarintSkip2Bits(%[1]d)=%[2]v\nreadVarintSkip2Bits(%[2]v)=%[3]d. expected to match", input, output, backAgain)
		}
		if bytesRead != len(output) {
			t.Errorf("readVarintSkip32Bits(%v)=_,%d, want %d", output, bytesRead, len(output))
		}
	}
}

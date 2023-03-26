package encoder

import (
	"math"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestEncodeAndDecode(t *testing.T) {
	testCases := []any{
		"",
		"hello world",
		int64(42),
		int64(math.MaxInt64),
		int64(math.MinInt64),
		decimal.RequireFromString("12.34"),
		decimal.RequireFromString("-12.34"),
		[]any{"foo", int64(42), decimal.RequireFromString("12.34"), []any{"1", "2", "3"}},
	}

	for _, input := range testCases {
		output, err := EncodeToString(input)
		if err != nil {
			t.Fatal(err)
		}
		backAgain, err := DecodeString(output)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(input, backAgain) {
			t.Errorf("EncodeToString(%[1]v)=%[2]x\nDecodeString(%[2]x)=%[3]v. expected to match", input, []byte(output), backAgain)
		}
	}
}

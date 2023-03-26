package encoder

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
)

const (
	intType     = 0b00
	decimalType = 0b01
	stringType  = 0b10
	listType    = 0b11
)

var (
	magicHeader = [4]byte{0x79, 0xc7, 0x00, 0x01} // for integrity and versioning
)

func EncodeToString(toSend any) (string, error) {
	buf := append([]byte{}, magicHeader[:]...)
	if err := encodeAndAppend(&buf, toSend); err != nil {
		return "", err
	}
	return string(buf), nil
}

func DecodeString(s string) (any, error) {
	in := []byte(s)
	if !bytes.Equal(in[:len(magicHeader)], magicHeader[:]) {
		return nil, errors.New("unknown file-type or version")
	}

	result, _ := decode(in[len(magicHeader):])
	return result, nil
}

func encodeAndAppend(outp *[]byte, toSend any) error {
	out := *outp
	toSend = normalizeToWireType(toSend)
	switch toSend := toSend.(type) {
	case int64:
		v := writeVarintSkip2Bits(toSend)
		v[0] |= intType << 6
		out = append(out, v...)

	case decimal.Decimal:
		v := writeVarintSkip2Bits(int64(toSend.Exponent()))
		v[0] |= decimalType << 6
		out = append(out, v...)
		out = append(out, writeVarint(toSend.CoefficientInt64())...)

	case string:
		v := writeUvarintSkip2Bits(uint64(len(toSend)))
		v[0] |= stringType << 6
		out = append(out, v...)
		out = append(out, toSend...)

	case []any:
		v := writeUvarintSkip2Bits(uint64(len(toSend)))
		v[0] |= listType << 6
		out = append(out, v...)
		for _, x := range toSend {
			if err := encodeAndAppend(&out, x); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown type %v", reflect.TypeOf(toSend))
	}
	*outp = out
	return nil
}

func normalizeToWireType(toSend any) any {
	v := reflect.ValueOf(toSend)
	switch {
	case v.CanInt():
		return v.Int()
	case v.CanFloat():
		vf := v.Float()
		if vf == float64(int64(vf)) {
			return int64(vf)
		}
		return decimal.NewFromFloat(vf)
	default:
		return toSend
	}
}

func decode(in []byte) (any, int) {
	i := 0
	switch in[i] >> 6 {
	case intType:
		v, s := readVarintSkip2Bits(in[i:])
		i += s
		return v, i

	case decimalType:
		exponent, s := readVarintSkip2Bits(in[i:])
		i += s
		mantissa, s := readVarint(in[i:])
		i += s
		v := decimal.New(mantissa, int32(exponent))
		return v, i

	case stringType:
		size, s := readUvarintSkip2Bits(in[i:])
		i += s
		v := string(in[i : i+int(size)])
		i += int(size)
		return v, i

	case listType:
		size, s := readUvarintSkip2Bits(in[i:])
		i += s
		output := []any{}
		for j := uint64(0); j < size; j++ {
			v, s := decode(in[i:])
			i += s
			output = append(output, v)
		}
		return output, i

	default:
		// We have 4 types above for a 2 bit flag, so the above should be
		// exhaustive. So the line below should never be executed.
		panic("Could not parse data-type")
	}
}

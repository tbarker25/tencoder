package encoder

func readUvarint(data []byte) (result uint64, bytesRead int) {
	x := uint64(0)
	s := 0
	i := 0
	for ; data[i] >= 0b1000_0000; i++ {
		x |= uint64(data[i]&0b0111_1111) << s
		s += 7
	}

	return x | uint64(data[i])<<s, i + 1
}

func readVarint(data []byte) (result int64, bytesRead int) {
	x, bytesRead := readUvarint(data)
	if x&1 == 1 {
		return ^int64(x >> 1), bytesRead
	}
	return int64(x >> 1), bytesRead
}

func readUvarintSkip2Bits(data []byte) (result uint64, bytesRead int) {
	if data[0]&0b0010_0000 == 0 {
		return uint64(data[0] & 0b0001_1111), 1
	}

	i := 0
	x := uint64(data[i] & 0b0001_1111)
	s := 5
	i++
	for data[i] >= 0b1000_0000 {
		x |= uint64(data[i]&0b0111_1111) << s
		s += 7
		i++
	}

	return x | uint64(data[i])<<s, i + 1
}

func readVarintSkip2Bits(data []byte) (result int64, bytesRead int) {
	x, bytesRead := readUvarintSkip2Bits(data)
	if x&1 == 1 {
		return ^int64(x >> 1), bytesRead
	}
	return int64(x >> 1), bytesRead
}

func writeUvarint(x uint64) []byte {
	buf := [10]byte{}
	i := 0
	for x >= 0b1000_0000 {
		buf[i] = byte(x) | 0b1000_0000
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return buf[:i+1]
}

func writeVarint(x int64) []byte {
	var n uint64
	if x >= 0 {
		n = uint64(x) << 1
	} else {
		n = uint64(^x)<<1 | 1
	}
	return writeUvarint(n)
}

func writeUvarintSkip2Bits(x uint64) []byte {
	if x < 0b0010_0000 {
		return []byte{byte(x)}
	}

	buf := [10]byte{}
	i := 0
	buf[i] = byte(x)&0b0001_1111 | 0b0010_0000
	x >>= 5
	i++
	for x >= 0b1000_0000 {
		buf[i] = byte(x) | 0b1000_0000
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return buf[:i+1]
}

func writeVarintSkip2Bits(x int64) []byte {
	var n uint64
	if x >= 0 {
		n = uint64(x) << 1
	} else {
		n = uint64(^x)<<1 | 1
	}
	return writeUvarintSkip2Bits(n)
}

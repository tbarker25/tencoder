# tencoder

## Usage

build binary:
```sh
go build
```

### Encoding
```sh
echo '[123, "foo", 12.34, ["a", "b", "c"]]' | \
  ./tencoder -encode -base64=true
  
> eccAAcQ2B4Nmb29DpBPDgWGBYoFj
```

### Decoding
```sh
echo eccAAcQ2B4Nmb29DpBPDgWGBYoFj | \
  ./tencoder -decode -base64=true

> [123,"foo","12.34",["a","b","c"]]
```

## Wire protocol
```
magicHeader                = 0x79c70001 // for versioning and format identification
Datum                      = String | Int64 | Decimal | List<Datum>
encode(d: Datum)           = magicHeader ++ encodeItem(d)
encodeItem(n: Int64)       = 0b00 ++ varintskip2(n)
encodeItem(x: Decimal)     = 0b01 ++ varintskip2(x.exponent) ++ varint(x.mantissa)
encodeItem(s: String)      = 0b10 ++ varuintskip2(s.length) ++ s.chars
encodeItem(l: List<Datum>) = 0b11 ++ varuintskip2(l.length) ++ concat(encodeItem(d) for d in l)

varint: variable-length encoding for a 64-bit integer
varintskip2: variable-length encoding a 64-bit integer, ignoring the 2 high-order-bits on the first byte. This is used so we have space for a 2-bit type ID but can still fit a full 64-bit integer
```

## Design considerations
This format is designed for simplicity of implementation, compactness of output, and minimal CPU and memory consumption

Low processing time:
- On-wire format consistently aligns with byte-boundries to keep processing time low. Bit-packing would reduce output size, at the cost of CPU time.

Compact output:
- No padding or redundant fields
- Variable-length encoding to reduce overhead for small items:
  * A small integer uses 1-byte on-wire. A small string uses 1+length bytes. A small decimal uses 2-bytes.

One-pass encoding:
- The format is designed so its possible to encode in a single pass with constant memory usage
  * Caveat: this implementation loads the entire payload into memory for simplicity reasons

Few arbitrary restrictions:
- Can encode arbitrary strings up to 2^64-1 long
- Agnostic to character encoding and can accept arbitrary bytes in strings
- The full range of a 64-bit signed integer can be marshalled. Integers outside this boundry will be transparently transmitted as a BigDecimal.
- Non-integer numbers are transmitted as a Decimal to avoid rounding issues (appropriateness of this depends on business rules).

Limitations:
- The format is not particularly extensible, although there is a version number that could be used for backwards compatibility if the format were changed.
- There is no self-synchronizing mechanism
- There's no error checking. Its assumed this would be performed at in a different layer.

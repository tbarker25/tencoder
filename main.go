package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.org/tbarker25/tencoder/encoder"
)

var (
	encodeFlag = flag.Bool("encode", false, "reads JSON from STDIN and writes encoded data to STDOUT")
	decodeFlag = flag.Bool("decode", false, "reads encoded data from STDIN and writes JSON to STDOUT")
	base64Flag = flag.Bool("base64", true, "encode to base64 format. Set to false to write binary instead")
)

func main() {
	flag.Parse()

	switch {
	case *encodeFlag:
		var v any
		json.NewDecoder(os.Stdin).Decode(&v)

		s, err := encoder.EncodeToString(v)
		if err != nil {
			log.Fatal(err)
		}
		if *base64Flag {
			s = base64.StdEncoding.EncodeToString([]byte(s))
		}
		os.Stdout.WriteString(s)

	case *decodeFlag:
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		if *base64Flag {
			dst := make([]byte, base64.StdEncoding.DecodedLen(len(input)))
			n, err := base64.StdEncoding.Decode(dst, input)
			if err != nil {
				log.Fatal(err)
			}
			input = dst[:n]
		}

		output, err := encoder.DecodeString(string(input))
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(os.Stdout).Encode(output)

	default:
		log.Printf("Exactly one of 'encode' or 'decode' flags must be set")
	}
}

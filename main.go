package main

import (
	"encoding/base64"
	"fmt"

	"github.com/naus3a/ichingcode/encoding/baseIching"
)

func encode(data []byte, f func([]byte) string) {
	s := f(data)
	fmt.Printf(s)
	fmt.Printf("\n")
}

func encodeB64(data []byte) {
	encode(data, base64.StdEncoding.EncodeToString)
}

func encodeIching(data []byte) {
	encode(data, baseIching.Encode)
}

func decode(sInput string, savePth string, f func(string) ([]byte, error)) {
	data, err := f(sInput)
	if err != nil {
		fmt.Printf("ERROR: could not decode input string")
		return
	}
	err = SaveFile(data, savePth)
	if err != nil {
		fmt.Printf("ERROR: cannot save")
	}
}

func decodeB64(sInput string, savePth string) {
	decode(sInput, savePth, base64.StdEncoding.DecodeString)
}

func decodeIching(sInput string, savePth string) {
	decode(sInput, savePth, baseIching.Decode)
}

func genericEncode(cli Cli) {
	switch cli.curEncoding {
	case ENCODING_B64:
		encodeB64(cli.inputData)
	case ENCODING_ICHING:
		encodeIching(cli.inputData)
	}
}

func genericDecode(cli Cli) {
	switch cli.curEncoding {
	case ENCODING_B64:
		decodeB64(cli.inputString, cli.outputPath)
	case ENCODING_ICHING:
		decodeIching(cli.inputString, cli.outputPath)
	}
}

func main() {
	var cli Cli
	cli.Setup()
	cli.Parse()
	if cli.err != nil {
		return
	}
	switch cli.curMode {
	case MODE_ENCODE:
		genericEncode(cli)
	case MODE_DECODE:
		genericDecode(cli)
	}
}

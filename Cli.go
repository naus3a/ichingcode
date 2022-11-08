package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Mode is the tool mode; can be encode or decode
type Mode int

const (
	MODE_ENCODE = iota
	MODE_DECODE
)

// EncodingType is the encoding type; can be base64 or iching
type EncodingType int

const (
	ENCODING_B64 = iota
	ENCODING_ICHING
)

// Cli is the cli parse manager
type Cli struct {
	curMode     Mode
	curEncoding EncodingType
	bVerbose    bool
	inputData   []byte
	inputString string
	outputPath  string
	err         error
}

var flagVerbose *bool
var outPth *string

// Setup setups the cli
func (c *Cli) Setup() {
	c.err = nil
	c.curMode = MODE_ENCODE
	c.curEncoding = ENCODING_ICHING

	flag.Usage = func() {
		fmt.Printf("Usage: ichingcode [optional params] [input]\n\n")
		fmt.Printf("Mode:\n")
		fmt.Printf("\tencode (default): encode input \n")
		fmt.Printf("\tdecode: decode input \n")
		fmt.Printf("\tbase64: use base64 \n")
		fmt.Printf("\tiching (default): use iching \n")
		fmt.Printf("Params:\n")
		flag.PrintDefaults()
	}

	flagVerbose = flag.Bool("verbose", false, "print verbose output")
	outPth = flag.String("o", "output.bin", "output path")
}

// Parse does the actual parsing
func (c *Cli) Parse() {
	args := os.Args[1:]
	if len(args) < 1 {
		c.err = errors.New("Not enough parameters")
		flag.Usage()
		return
	}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "encode":
			c.curMode = MODE_ENCODE
		case "decode":
			c.curMode = MODE_DECODE
		case "base64":
			c.curEncoding = ENCODING_B64
		case "iching":
			c.curEncoding = ENCODING_ICHING
		}
	}
	flag.Parse()
	c.bVerbose = *flagVerbose
	c.outputPath = *outPth
	c.parseInputParameter(args[len(args)-1])
	if c.bVerbose {
		c.printVerboseSetupInfo()
	}
}

func (c *Cli) parseInputParameter(param string) {
	switch c.curMode {
	case MODE_ENCODE:
		c.inputData, c.err = LoadFile(param)
		if c.err != nil {
			fmt.Printf("ERROR: cannot load " + param)
			flag.Usage()
			return
		}
	case MODE_DECODE:
		c.inputString = param
	}
}

func (c *Cli) printVerboseSetupInfo() {
	fmt.Printf("ichingcode: ")
	switch c.curMode {
	case MODE_ENCODE:
		fmt.Printf("encoding binary data to ")
	case MODE_DECODE:
		fmt.Printf("decoding data from ")
	}
	switch c.curEncoding {
	case ENCODING_B64:
		fmt.Printf("base64")
	case ENCODING_ICHING:
		fmt.Printf("iching")
	}
	fmt.Printf("\n")
}

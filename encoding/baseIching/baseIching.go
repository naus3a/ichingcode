package baseIching

import (
	"errors"
	"fmt"
)

const hexagrams = "䷀䷁䷂䷃䷄䷅䷆䷇䷈䷉䷊䷋䷌䷍䷎䷏䷐䷑䷒䷓䷔䷕䷖䷗䷘䷙䷚䷛䷜䷝䷞䷟䷠䷡䷢䷣䷤䷥䷦䷧䷨䷩䷪䷫䷬䷭䷮䷯䷰䷱䷲䷳䷴䷵䷶䷷䷸䷹䷺䷻䷼䷽䷾䷿"
const trigrams = "☰☱☲☳☴☵☶☷"

///
/// iching parsing
///

func getUnicodeByIndex(idx int, dictionary string) string {
	from := idx * 3
	to := from + 3
	return string(dictionary[from:to])
}

func getIndexForUnicode(c string, dictionary string) (idx int, err error) {
	err = nil
	sz := len(dictionary) / 3
	for i := 0; i < sz; i++ {
		cur := getUnicodeByIndex(i, dictionary)
		if cur == c {
			idx = i
			return
		}
	}
	err = errors.New("Not found")
	return
}

func getTrigramByIndex(idx int) string {
	return getUnicodeByIndex(idx, trigrams)
}

func getHexagramByIndex(idx int) string {
	return getUnicodeByIndex(idx, hexagrams)
}

func getPaddingTrigram() string {
	return getTrigramByIndex(0)
}

func isPaddingTrigram(src *string, pos int) bool {
	if pos < 0 || pos >= len(*src) {
		return false
	}
	gram := getUnicodeByIndex(pos, *src)
	return gram == getPaddingTrigram()
}

func countPaddingTrigramsAtEnd(src *string) int {
	n := 0
	srcSz := len(*src)
	if srcSz > 0 {
		srcSz /= 3
		for i := srcSz - 1; i >= 0; i-- {
			if isPaddingTrigram(src, i) {
				n++
			} else {
				return n
			}
		}
	}
	return n
}

func getIndexForHexagram(c string) (idx int, err error) {
	return getIndexForUnicode(c, hexagrams)
}

func hexagram2byte(src *string, idx int) (val byte, err error) {
	hexagram := getUnicodeByIndex(idx, *src)
	var idxHexagram int
	idxHexagram, err = getIndexForHexagram(hexagram)
	val = byte(idxHexagram)
	return
}

///
/// bit logic
///

func makeRuneShift(pos int, bits int, stride int) int {
	return stride - (bits * pos)
}

func make8bitsRuneShift(pos int) int {
	return makeRuneShift(pos, 8, 16)
}

func make6bitsRuneShift(pos int) int {
	return makeRuneShift(pos, 6, 18)
}

func packByteInRune(val byte, pos int) rune {
	if pos < 0 {
		pos = 0
	} else if pos > 2 {
		pos = 2
	}
	r := rune(val)
	var shift int
	shift = make8bitsRuneShift(pos)
	r = r << shift
	return r
}

func pack3bytesInRune(src []byte, startPos int) (r rune, err error) {
	if len(src) < 3 {
		err = errors.New("buffer is smaller than 3 bytes")
		return
	}
	if startPos < 0 || startPos > (len(src)-3) {
		err = errors.New("start position out of range")
		return
	}
	rr := make([]rune, 3)
	for i := 0; i < 3; i++ {
		rr[i] = packByteInRune(src[startPos+i], i)
		r = r | rr[i]
	}
	return
}

func pack6bitsInRune(bits byte, pos int) rune {
	if pos < 0 {
		pos = 0
	} else if pos > 3 {
		pos = 3
	}
	r := rune(bits)
	var shift int
	shift = make6bitsRuneShift(pos)
	r = r << shift
	return r
}

func pack24bitsInRune(bits []byte) rune {
	rr := make([]rune, 4)
	var r rune
	r = 0
	for i := 0; i < 4; i++ {
		rr[i] = pack6bitsInRune(bits[i], i)
		r = r | rr[i]
	}
	return r
}

func extractPackedByteFromRune(r rune, pos int) byte {
	if pos < 0 {
		pos = 0
	} else if pos > 2 {
		pos = 2
	}
	var b byte
	shift := make8bitsRuneShift(pos)
	b = byte(r >> shift)
	return b
}

///
/// public
///

// Encode a binary buffer to iching
func Encode(src []byte) string {
	var err error
	s := ""
	if len(src) < 1 {
		return s
	}

	idxHexagrams := make([]rune, 4)
	numTriplets := len(src) / 3
	remain := len(src) % 3
	idxTripletStart := 0

	for iT := 0; iT < numTriplets; iT++ {
		// convert 3x 8bit source bytes into 4 bytes
		var val rune
		val, err = pack3bytesInRune(src, idxTripletStart)
		if err != nil {
			fmt.Println(err)
			return s
		}

		// for every byte get the encoding vocabulary index
		for iI := 0; iI < len(idxHexagrams); iI++ {
			idxHexagrams[iI] = val >> make6bitsRuneShift(iI) & 0x3F
			// and append hexagrams
			s = s + getHexagramByIndex(int(idxHexagrams[iI]))
		}

		idxTripletStart += 3
	}

	// if we got remain, we will encode less than 4 hexagrams...
	if remain > 0 {
		//... since we're encoding triplets, remain can be 1 or 2
		// meaning we have at least 1 byte to parse...
		val := packByteInRune(src[idxTripletStart], 0)
		if remain == 2 {
			//... and we have a 2nd byte in case it's 2
			val |= packByteInRune(src[idxTripletStart+1], 1)
		}
		//since we have at least 1 byte, we'll encode at least 2 hexagrams...
		for iI := 0; iI < 2; iI++ {
			idxHexagrams[iI] = val >> make6bitsRuneShift(iI) & 0x3F
			s = s + getHexagramByIndex(int(idxHexagrams[iI]))
		}
		//...finally...
		switch remain {
		case 1:
			//... all bytes encoded; just fill in with 2 trigrams
			s += getPaddingTrigram()
			s += getPaddingTrigram()
		case 2:
			//... a 3rd hexagram is needed...
			idxHexagrams[2] = val >> make6bitsRuneShift(2) & 0x3F
			s = s + getHexagramByIndex(int(idxHexagrams[2]))
			//... and finally a trigram
			s += getPaddingTrigram()
		}
	}

	return s
}

func hasPaddingTrigramAtIndex(src string, idx int) bool {
	if idx >= len(src)/3 {
		return false
	}
	c := getUnicodeByIndex(idx, src)
	return c == getPaddingTrigram()
}

func Decode(src string) (dst []byte, err error) {
	rawSrcSize := len(src)
	if rawSrcSize == 0 {
		err = errors.New("Input string cannot be empty")
		return
	}
	//we're parsing unicode characters, which take 3 bytes each
	if rawSrcSize%3 != 0 {
		err = errors.New("Unsupported encoding")
		return
	}

	// let's figure out the encoded file size and evenutual remain
	var ichingSrcSize int
	var outputBitSize int
	var outputByteSize int
	var numQuadruplets int
	var numRemainHexagrams int
	var numRemainBytes int
	ichingSrcSize = rawSrcSize / 3

	if ichingSrcSize%4 != 0 {
		err = errors.New("String size is not multile of 4; unsupported encoding")
		return
	}

	numPaddingTrigrams := countPaddingTrigramsAtEnd(&src)
	numRemainHexagrams = 0
	numRemainBytes = 0
	if numPaddingTrigrams > 0 {
		if numPaddingTrigrams > 2 {
			err = errors.New("More than 2 padding characters")
			return
		}
		numRemainHexagrams = 4 - numPaddingTrigrams
		switch numRemainHexagrams {
		case 2:
			numRemainBytes = 1
		case 3:
			numRemainBytes = 2
		}
	}
	numHexagrams := ichingSrcSize - numPaddingTrigrams
	numQuadruplets = numHexagrams / 4
	outputBitSize = (numQuadruplets * 24) + (numRemainBytes * 8)
	outputByteSize = outputBitSize / 8
	dst = make([]byte, outputByteSize)

	//now let's parse quadruplets
	valHexagrams := make([]byte, 4)
	idxHexagram := 0
	idxByte := 0
	for iQ := 0; iQ < numQuadruplets; iQ++ {
		//get hexagrams in groups of 4
		// 1 hexagram is 6 bit, so a group of 4 is 24 bit, aka 3 bytes
		for iH := 0; iH < 4; iH++ {
			valHexagrams[iH], err = hexagram2byte(&src, idxHexagram)
			if err != nil {
				fmt.Println(err)
				return
			}
			idxHexagram++
		}

		//let's make a 32bit buffer to put our (6*4 bits) idx in
		buf := pack24bitsInRune(valHexagrams)

		//then let's unpack the 3 bytes we obtained into the dst buffer
		for iB := 0; iB < 3; iB++ {
			if idxByte < len(dst) {
				dst[idxByte] = extractPackedByteFromRune(buf, iB)
			} else {
				err = errors.New("something went wrong: out of boundaries")
				fmt.Println(err)
				return
			}
			idxByte++
		}
	}

	//if we have a remain...
	if numRemainBytes > 0 {
		//... we have at least 1 extra byte to encode, aka 2 hexagrams...
		for iH := 0; iH < 2; iH++ {
			valHexagrams[iH], err = hexagram2byte(&src, idxHexagram)
			if err != nil {
				fmt.Println(err)
				return
			}
			idxHexagram++
		}
		//... if we have a 3rd extra hex, let's get its index too...
		if numRemainHexagrams > 2 {
			valHexagrams[2], err = hexagram2byte(&src, idxHexagram)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		//...finally extract bytes...
		buf := pack24bitsInRune(valHexagrams)
		switch numRemainBytes {
		case 1:
			dst[outputByteSize-1] = extractPackedByteFromRune(buf, 0)
		case 2:
			dst[outputByteSize-2] = extractPackedByteFromRune(buf, 0)
			dst[outputByteSize-1] = extractPackedByteFromRune(buf, 1)
		}
	}

	return
}

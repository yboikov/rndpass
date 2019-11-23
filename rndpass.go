package rndpass

import (
	"bytes"
	cr "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
)

type chars []byte

// Config struct
// some info
type Config struct {
	// Expected password length
	Length   int
	Symbols  int
	Numbers  int
	Lower    int
	Upper    int
	NoRepeat bool
	Cons     bool //consecutive
	Exclude  string
}

func New(l int) Config {
	c := Config{
		Length:  l,
		Upper:   1,
		Lower:   1,
		Numbers: 1,
		Symbols: 1,
	}

	return c
}

func (c Config) Gen() (string, error) {
	var p, r []byte

	if c.Length > 0 && c.Lower+c.Upper+c.Symbols+c.Numbers > c.Length {
		return "", fmt.Errorf("Required password length is less than minimum required characters")
	}

	if c.Length == 0 {
		c.Length = c.Lower + c.Upper + c.Symbols + c.Numbers
	}

	if c.Upper > 0 {
		ch, rest := pick(c.Upper, c.NoRepeat, c.Cons, c.Exclude, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
		p = append(p, ch...)
		r = append(r, rest...)
	}

	if c.Lower > 0 {
		ch, rest := pick(c.Lower, c.NoRepeat, c.Cons, c.Exclude, []byte("abcdefghijklmnopqrstuvwxyz"))
		p = append(p, ch...)
		r = append(r, rest...)
	}

	if c.Numbers > 0 {
		ch, rest := pick(c.Numbers, c.NoRepeat, c.Cons, c.Exclude, []byte("0987654321"))
		p = append(p, ch...)
		r = append(r, rest...)
	}

	if c.Symbols > 0 {
		ch, rest := pick(c.Symbols, c.NoRepeat, c.Cons, c.Exclude, []byte(" !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"))
		p = append(p, ch...)
		r = append(r, rest...)
	}

	if c.Length > c.Lower+c.Upper+c.Symbols+c.Numbers {
		rest, _ := pick(c.Length-(c.Lower+c.Upper+c.Symbols+c.Numbers), c.NoRepeat, c.Cons, c.Exclude, r)
		p = append(p, rest...)

	}
	p, _ = pick(c.Length, true, c.Cons, c.Exclude, p)
	if len(p) == 0 {
		return "", fmt.Errorf("There are no characters to pickup from!")
	}
	return string(p), nil
}

func pick(n int, r, c bool, e string, chars []byte) ([]byte, []byte) {
	var tempChars []byte
	var pos int
	var b [8]byte
	pickUp := make([]byte, n)
	cr.Read(b[:])
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	// remove exluded characters from chars
	for _, v := range e {
		pos = bytes.IndexByte(chars, byte(v))
		if pos > -1 {
			chars = removeByte(chars, pos)
		}
	}

	if len(chars) == 0 {
		return []byte{}, []byte{}
	}

	for k := 0; k < n; k++ {
		if c && k > 0 {
			tempChars = consGroup(chars, rune(pickUp[k-1]))
			pos = rand.Intn(len(tempChars))
			pickUp[k] = tempChars[pos]
		} else {
			pos = rand.Intn(len(chars))
			pickUp[k] = chars[pos]
		}

		// if no repeat is set remove picked up char
		if r {
			chars = removeByte(chars, bytes.IndexByte(chars, pickUp[k]))

		}
		if len(chars) == 0 {
			return pickUp[:], chars[:]
		}
	}
	return pickUp[:], chars[:]
}

func consGroup(b []byte, ch rune) []byte {
	var newGroup []byte

	for _, v := range b {
		if ch >= 'a' && ch <= 'z' {
			if v < 'a' || v > 'z' {
				newGroup = append(newGroup, byte(v))
				continue
			}
		}
		if ch >= 'A' && ch <= 'Z' {
			if v < 'A' || v > 'Z' {
				newGroup = append(newGroup, byte(v))
				continue
			}
		}
		if ch >= '0' && ch <= '9' {
			if v < '0' || v > '9' {
				newGroup = append(newGroup, byte(v))
				continue
			}
		}
		if (ch < '0' || ch > '9') && (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') {
			if (v >= '0' && v <= '9') || (v >= 'A' && v <= 'Z') || (v >= 'a' && v <= 'z') {
				newGroup = append(newGroup, byte(v))
			}
		}
	}
	if len(newGroup) > 0 {
		return newGroup
	}
	return b

}

func removeByte(b []byte, pos int) []byte {
	b[pos] = b[len(b)-1]
	b[len(b)-1] = 0
	b = b[:len(b)-1]
	return b[:]
}

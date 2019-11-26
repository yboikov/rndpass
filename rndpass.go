package rndpass

import (
	"bytes"
	cr "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"regexp"
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
	var pass []byte
	pass, err := c.GenBytes()
	if err != nil {
		return "", err
	}
	return string(pass), nil
}

func (c Config) GenBytes() ([]byte, error) {
	var p, rest []byte
	if c.Length > 0 && c.Lower+c.Upper+c.Symbols+c.Numbers > c.Length {
		return []byte{}, fmt.Errorf("Required password length is less than minimum required characters")
	}

	if c.Lower+c.Upper+c.Symbols+c.Numbers == 0 {
		// default
		c.Lower = 1
		c.Upper = 1
		c.Numbers = 1
		c.Symbols = 1
	}

	if c.Length == 0 {
		c.Length = c.Lower + c.Upper + c.Symbols + c.Numbers
	}

	if c.Upper > 0 {
		ch, r := pick(c.Upper, c.NoRepeat, c.Cons, c.Exclude, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
		p = append(p, ch...)
		rest = append(rest, r...)
	}

	if c.Lower > 0 {
		ch, r := pick(c.Lower, c.NoRepeat, c.Cons, c.Exclude, []byte("abcdefghijklmnopqrstuvwxyz"))
		p = append(p, ch...)
		rest = append(rest, r...)

	}

	if c.Numbers > 0 {
		ch, r := pick(c.Numbers, c.NoRepeat, c.Cons, c.Exclude, []byte("0987654321"))
		p = append(p, ch...)
		rest = append(rest, r...)

	}

	if c.Symbols > 0 {
		ch, r := pick(c.Symbols, c.NoRepeat, c.Cons, c.Exclude, []byte(" !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"))
		p = append(p, ch...)
		rest = append(rest, r...)

	}

	if c.Length > c.Lower+c.Upper+c.Symbols+c.Numbers {
		freeCount := c.Length - (c.Lower + c.Upper + c.Symbols + c.Numbers)
		rest, _ = pick(freeCount, c.NoRepeat, c.Cons, c.Exclude, rest)
		p = append(p, rest...)
	}

	p, _ = pick(c.Length, true, c.Cons, c.Exclude, p)
	if c.Cons {
		p = consByte(p)
	}
	if len(p) == 0 {
		return []byte{}, fmt.Errorf("There are no characters to pickup from!")
	}
	return p, nil
}

func pick(n int, r, c bool, e string, chars []byte) ([]byte, []byte) {
	var pos int
	var b [8]byte
	var pickUp []byte
	cr.Read(b[:])
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	// remove exluded characters from chars
	i := bytes.IndexAny(chars, e)
	for i > -1 {
		chars = removeByte(chars, i)
		i = bytes.IndexAny(chars, e)
	}
	if len(chars) == 0 {
		return []byte{}, []byte{}
	}

	for k := 0; k < n; k++ {
		pos = rand.Intn(len(chars))
		pickUp = append(pickUp, chars[pos])

		// if no repeat is set remove picked up char
		if r {
			chars = removeByte(chars, pos)

		}
		if len(chars) == 0 {
			return pickUp[:], chars[:]
		}
	}
	return pickUp[:], chars[:]
}

func consByte(a []byte) []byte {
	var tmp []byte
	var randSeed [8]byte
	var re *regexp.Regexp
	cr.Read(randSeed[:])
	rand.Seed(int64(binary.LittleEndian.Uint64(randSeed[:])))
	tmp = append(tmp, a[:1]...)
	a = removeByte(a, 0)

	for k, _ := range a {
		if k > 0 {
			tempChars := consGroup(a, rune(tmp[k-1]))
			pos := rand.Intn(len(tempChars))
			tmp = append(tmp, tempChars[pos])
			a = removeByte(a, bytes.IndexByte(a, tempChars[pos]))
		}
	}

	tmp = append(tmp, a...)
	badValuesRe := regexp.MustCompile(`([a-z]{2,}|[A-Z]{2,}|[0-9]{2,})`)
	badValues := badValuesRe.Find(tmp)
	if len(badValues) > 0 {
		for k := len(badValues) - 1; k >= 0; k-- {
			//for _, v := range badValues {
			if tmp[k] >= 'A' && tmp[k] <= 'Z' {
				re = regexp.MustCompile(`[^A-Z]{2}`)
			}
			if tmp[k] >= 'a' && tmp[k] <= 'z' {
				re = regexp.MustCompile(`[^a-z]{2}`)
			}
			if tmp[k] >= '0' && tmp[k] <= '9' {
				re = regexp.MustCompile(`[^0-9]{2}`)
			}
			if (tmp[k] < '0' || tmp[k] > '9') && (tmp[k] < 'A' || tmp[k] > 'Z') && (tmp[k] < 'a' || tmp[k] > 'z') {
				re = regexp.MustCompile(`[^A-Za-z0-9]{2,}`)
			}
			loc := re.FindIndex(tmp)
			if loc == nil {
				continue
			}
			tmp = append(tmp, 0)
			copy(tmp[loc[1]:], tmp[loc[1]-1:])
			tmp[loc[1]-1] = tmp[k]

			tmp[len(tmp)-1] = 0
			tmp = tmp[:len(tmp)-1]
		}
	}

	return tmp
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

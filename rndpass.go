package rndpass

import (
	cr "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type Config struct {
	// Expected password length
	Length   int
	Symbols  int
	Numbers  int
	Lower    int
	Upper    int
	NoRepeat bool
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
	// length := 10
	// uCount := 4
	// lCount := 4
	var reg *regexp.Regexp
	var excludeRex string
	if len(c.Exclude) > 0 {
		exclude := moveToEnd([]byte(c.Exclude), 45)
		excludeRex = fmt.Sprintf("[%s]+", regexp.QuoteMeta(string(*exclude)))
	}
	reg = regexp.MustCompile(excludeRex)

	upperSet := []byte(reg.ReplaceAllString("ABCDEFGHIJKLMNOPQRSTUVWXYZ", ""))
	lowerSet := []byte(reg.ReplaceAllString("abcdefghijklmnopqrstuvwxyz", ""))
	numberSet := []byte(reg.ReplaceAllString("0987654321", ""))
	symbolSet := []byte(reg.ReplaceAllString(" !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", ""))

	var allBytes []byte
	var rndBytes []byte

	if c.Upper > 0 {

		rndBytes = append(rndBytes, c.getRndChars(&upperSet, c.Upper)...)
		allBytes = append(allBytes, upperSet...)

	}
	if c.Lower > 0 {
		rndBytes = append(rndBytes, c.getRndChars(&lowerSet, c.Lower)...)
		allBytes = append(allBytes, lowerSet...)

	}
	if c.Numbers > 0 {
		rndBytes = append(rndBytes, c.getRndChars(&numberSet, c.Numbers)...)
		allBytes = append(allBytes, numberSet...)

	}
	if c.Symbols > 0 {
		rndBytes = append(rndBytes, c.getRndChars(&symbolSet, c.Symbols)...)
		allBytes = append(allBytes, symbolSet...)

	}

	if c.NoRepeat {
		sorted := moveToEnd([]byte(rndBytes), 45)
		s := fmt.Sprintf("[%s]+", regexp.QuoteMeta(string(*sorted)))

		reg := regexp.MustCompile(s)
		// remove already used bytes
		allBytes = []byte(reg.ReplaceAllString(string(allBytes), ""))
	}

	if len(rndBytes) < c.Length {
		rndBytes = append(rndBytes, c.getRndChars(&allBytes, c.Length-len(rndBytes))...)
	}

	// return shuffle(rndBytes), nil
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(rndBytes), func(i, j int) {
		rndBytes[i], rndBytes[j] = rndBytes[j], rndBytes[i]
	})
	return rndBytes, nil

}

func (cfg Config) getRndChars(chrs *[]byte, countChars int) []byte {
	if len(*chrs) == 0 {
		return *chrs
	}

	var b = make([]byte, countChars+8)
	var rndChrs []byte
	cr.Read(b[:])
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	for c := 0; c < countChars; c++ {
		idx := int(b[c]) % len(*chrs)
		rndChrs = append(rndChrs, (*chrs)[idx])
		if cfg.NoRepeat {
			*chrs = append((*chrs)[:idx], (*chrs)[idx+1:]...)
			if len(*chrs) == 0 {
				return rndChrs
			}
		}
	}
	return rndChrs
}

func moveToEnd(b []byte, mB byte) *[]byte {
	newB := []byte{}
	for _, v := range b {
		if v == mB {
			defer func(bb *[]byte, v byte) {
				*bb = append(*bb, v)
				//fmt.Println(bb)
			}(&newB, v)
			continue
		}
		newB = append(newB, v)
	}
	return &newB
}

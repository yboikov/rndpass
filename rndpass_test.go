package rndpass

import (
	"bytes"
	"testing"
)

var result []byte

func testGen(le, l, u, n, s int, e string, noRepeat bool, t *testing.T) {
	c := &Config{Length: le, Numbers: n, Lower: l, Upper: u, Symbols: s, Exclude: e, NoRepeat: noRepeat}
	a, err := c.GenBytes()
	if err != nil {
		t.Errorf("%s", err)
	}
	var lowerCase, upperCase, numbers, symbols int
	for _, v := range a {
		if v >= 'a' && v <= 'z' {
			lowerCase++
			continue
		}
		if v >= 'A' && v <= 'Z' {
			upperCase++
			continue
		}
		if v >= '0' && v <= '9' {
			numbers++
			continue
		}
		symbols++
	}
	if lowerCase < l && !noRepeat {
		t.Errorf("Want minimum %d lowercase letters have %d %s", l, lowerCase, a)
	}
	if upperCase < u && !noRepeat {
		t.Errorf("Want minimum %d uppercase letters have %d", u, upperCase)
	}
	if numbers < n && !noRepeat {
		t.Errorf("Want minimum %d numbers have %d", n, numbers)
	}
	if symbols < s && !noRepeat {
		t.Errorf("Want minimum %d symbols have %d %s", s, symbols, a)
	}

	if e != "" {
		if bytes.IndexAny(a, e) > -1 {
			t.Errorf("Found excluded character %s in %s", e, a)
		}
	}

	if noRepeat {
		for k, v := range a {
			if bytes.IndexByte(a[k+1:], v) > -1 {
				t.Errorf("found duplicate character %s in %s", string(v), a)
			}
		}
	}

}

func TestGenLowerCase(t *testing.T)     { testGen(40, 20, 0, 0, 0, "", false, t) }
func TestGenLowerCaseOver(t *testing.T) { testGen(40, 40, 0, 0, 0, "", false, t) }
func TestGenUpperCase(t *testing.T)     { testGen(40, 0, 26, 0, 0, "", false, t) }
func TestGenUpperCaseOver(t *testing.T) { testGen(40, 0, 40, 0, 0, "", false, t) }
func TestGenNumbers(t *testing.T)       { testGen(40, 0, 0, 26, 0, "", false, t) }
func TestGenNumbersOver(t *testing.T)   { testGen(40, 0, 0, 26, 0, "", false, t) }
func TestGenSymbols(t *testing.T)       { testGen(40, 0, 0, 0, 13, "", false, t) }
func TestGenSymbolsOver(t *testing.T)   { testGen(40, 0, 0, 0, 36, "", false, t) }

func TestGenNoRepeat(t *testing.T)        { testGen(40, 30, 30, 30, 30, "", true, t) }
func TestGenNoRepeatExclude(t *testing.T) { testGen(40, 30, 30, 30, 30, "#@$234SDFdsf", true, t) }

func benchmark_Gen(size int, b *testing.B) {
	var r []byte
	c := &Config{
		Length:  size,
		Exclude: "@#$ASD",
	}
	for i := 0; i <= b.N; i++ {
		r, _ = c.GenBytes()
	}
	result = r
}

func BenchmarkGen10(b *testing.B) { benchmark_Gen(10, b) }
func BenchmarkGen20(b *testing.B) { benchmark_Gen(20, b) }
func BenchmarkGen30(b *testing.B) { benchmark_Gen(30, b) }
func BenchmarkGen40(b *testing.B) { benchmark_Gen(40, b) }
func BenchmarkGen50(b *testing.B) { benchmark_Gen(50, b) }
func BenchmarkGen60(b *testing.B) { benchmark_Gen(60, b) }

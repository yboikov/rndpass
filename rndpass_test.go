package rndpass

import (
	"strings"
	"testing"
)

var result string

// const (
// 	reNums    = `^([^0-9]*[0-9]){%d}[^0-9]*`
// 	reUpper   = `^([^A-Z]*[A-Z]){%d}[^A-Z]*`
// 	reLower   = `^([^a-z]*[a-z]){%d}[^a-z]*`
// 	reSymbols = `^([A-Za-z0-9]*[^A-Za-z0-9]){%d}[A-Za-z0-9]*`
// )

func testGen(le, l, u, n, s int, e string, noRepeat bool, t *testing.T) {
	c := &Config{Length: le, Numbers: n, Lower: l, Upper: u, Symbols: s, Exclude: e, NoRepeat: noRepeat}
	a, err := c.Gen()
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
	if lowerCase < l {
		t.Errorf("Want minimum %d lowercase letters have %d %s", l, lowerCase, a)
	}
	if upperCase < u {
		t.Errorf("Want minimum %d uppercase letters have %d", u, upperCase)
	}
	if numbers < n {
		t.Errorf("Want minimum %d numbers have %d", n, numbers)
	}
	if symbols < s {
		t.Errorf("Want minimum %d symbols have %d %s", s, symbols, a)
	}

	if e != "" {
		for _, v := range e {
			if strings.Index(a, string(v)) > -1 {
				t.Errorf("Found excluded character %s in %s", string(v), a)
			}
		}
	}

	if noRepeat {
		for k, v := range a {
			i := strings.Index(a[k+1:], string(v))
			if i > -1 {
				t.Errorf("Found duplicate character %s at %d in %s", string(v), i, a)
			}
		}
	}

	// re := regexp.MustCompile(fmt.Sprintf(reNums, n))
	// if !re.MatchString(a) {
	// 	t.Errorf("Have less than %d numbers", n)
	// }

}

func TestGenLowerCase(t *testing.T) { testGen(40, 20, 0, 0, 0, "", false, t) }
func TestGenUpperCase(t *testing.T) { testGen(40, 10, 20, 0, 10, "", false, t) }
func TestGenNumbers(t *testing.T)   { testGen(40, 0, 10, 20, 10, "", false, t) }
func TestGenSymbols(t *testing.T)   { testGen(40, 10, 10, 0, 20, "", false, t) }

func TestGenExclude(t *testing.T) { testGen(50, 20, 10, 0, 20, "ABC!@#", false, t) }
func TestGenRepeat(t *testing.T)  { testGen(50, 20, 10, 0, 20, "ABC!@#", true, t) }

func Test_pick(t *testing.T) {
	a, b := pick(10, false, false, "", []byte("test"))
	if string(a) == "test" && string(b) != "" {
		t.Errorf("pick return wrong %s != \"test\"", string(a))
	}
}

func Test_consGroup(t *testing.T) {
	a := consGroup([]byte("1234567890abcdefghijklmnopqrstuvwxyz"), '1')
	if string(a) != "abcdefghijklmnopqrstuvwxyz" {
		t.Errorf("expects %s = %s", string(a), "abcdefghijklmnopqrstuvwxyz")
	}
	a = consGroup([]byte("1234567890abcdefghijklmnopqrstuvwxyz"), 'a')
	if string(a) != "1234567890" {
		t.Errorf("expects %s = %s", string(a), "1234567890")
	}
}

func benchmark_Gen(size int, b *testing.B) {
	var r string
	c := &Config{
		Length:  size,
		Exclude: "@#$ASD",
	}
	for i := 0; i <= b.N; i++ {
		r, _ = c.Gen()
	}
	result = r
}

func BenchmarkGen10(b *testing.B) { benchmark_Gen(10, b) }
func BenchmarkGen20(b *testing.B) { benchmark_Gen(20, b) }
func BenchmarkGen30(b *testing.B) { benchmark_Gen(30, b) }
func BenchmarkGen40(b *testing.B) { benchmark_Gen(40, b) }
func BenchmarkGen50(b *testing.B) { benchmark_Gen(50, b) }
func BenchmarkGen60(b *testing.B) { benchmark_Gen(60, b) }

package rndpass

import (
	"bytes"
	"regexp"
	"testing"
)

var result []byte

// const (
// 	reNums    = `^([^0-9]*[0-9]){%d}[^0-9]*`
// 	reUpper   = `^([^A-Z]*[A-Z]){%d}[^A-Z]*`
// 	reLower   = `^([^a-z]*[a-z]){%d}[^a-z]*`
// 	reSymbols = `^([A-Za-z0-9]*[^A-Za-z0-9]){%d}[A-Za-z0-9]*`
// )

func testGen(le, l, u, n, s int, e string, noRepeat, cons bool, t *testing.T) {
	c := &Config{Length: le, Numbers: n, Lower: l, Upper: u, Symbols: s, Exclude: e, NoRepeat: noRepeat, Cons: cons}
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

	if c.Cons {
		re := regexp.MustCompile(`([a-z]{2,}|[A-Z]{2,}|[0-9]{2,})`)
		if re.Match(a) {
			t.Errorf("Found consecutive characters %s in %s", string(re.Find(a)), a)
		}
	}

}

// func TestGenLowerCase(t *testing.T) { testGen(40, 26, 0, 0, 0, "", false, false, t) }
// func TestGenUpperCase(t *testing.T) { testGen(40, 0, 26, 0, 0, "", false, false, t) }
// func TestGenNumbers(t *testing.T)   { testGen(40, 0, 0, 26, 10, "", false, false, t) }
// func TestGenSymbols(t *testing.T)   { testGen(40, 0, 0, 0, 26, "", false, false, t) }

// func TestGenExclude(t *testing.T) { testGen(50, 20, 10, 0, 20, "ABC!@#", false, false, t) }
// func TestGenRepeat(t *testing.T)  { testGen(50, 20, 10, 0, 20, "ABC!@#", true, false, t) }

// func TestGenCons10(t *testing.T) {
// 	// var per int
// 	//TODO: calc percent of bad passwords
// 	for x := 0; x <= 1; x++ {
// 		testGen(32, 1, 1, 1, 1, "", true, true, t)
// 	}
// }
func Test_pick(t *testing.T) {
	a, b := pick(10, false, false, "", []byte("test"))
	if string(a) == "test" && string(b) != "" {
		t.Errorf("pick return wrong %s != \"test\"", string(a))
	}
}

func Test_consGroup(t *testing.T) {
	a := consGroup([]byte("01239abcdxyz!@#|\""), '1')
	if string(a) != "abcdxyz!@#|\"" {
		t.Errorf("expects %s = %s", string(a), "abcdxyz!@#|\"")
	}
	a = consGroup([]byte("01239abcdxyz!@#|\""), 'a')
	if string(a) != "01239!@#|\"" {
		t.Errorf("expects %s = %s", string(a), "01239!@#|\"")
	}
	a = consGroup([]byte("01239abcdxyz"), '%')
	if string(a) != "01239abcdxyz" {
		t.Errorf("expects %s = %s", string(a), "01239abcdxyz")
	}

}

func Test_consByte(t *testing.T) {

	// a := []byte("jM~H1|Y0Tt(m)rFyJs^RfCDVG")
	// fmt.Println("Orig: ", string(a))

	// b := consByte(a)
	// fmt.Println("Cons: ", string(b))
}

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

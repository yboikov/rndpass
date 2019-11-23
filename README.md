# rndpass

Random password generator library.

### Configs
- Disallow repeat characters
- Disallow consecutive characters from same type ( sagdfge, 3124353, DARGW, @#$%^&)
- Exclude charaters list


## Getting Started

```go get -u github.com/yboikov/rndpass```


## Examples

basic
```
package main

import (
	"fmt"

	"github.com/yboikov/rndpass"
)

func main() {
	f := rndpass.New(12)
	p, err := f.Gen()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
```


generate password with length 21 characters with only Upper case, Lower case and numbers
```
package main

import (
	"fmt"

	"github.com/yboikov/rndpass"
)

func main() {
	f := &rndpass.Config{
		Length: 21,
		Numbers:  1, // have 10 numbers if possible ( with noRepeat=true and exclude="12345" only 67890 will be use )
		Lower:    1, // have 10 lower case letters  
		Upper:    1, // have 10 upper case letters 
	}
	p, err := f.Gen()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
}
```


advanced
```
package main

import (
	"fmt"

	"github.com/yboikov/rndpass"
)

func main() {
	f := &rndpass.Config{
		// Length: 50 // set lenght greater than numbers,lower,upper and symbols count sum. otherwise passwowrd will be 40 characters or less if we have noRepat =true and exlude too many characters 
		Numbers:  10, // have 10 numbers if possible ( with noRepeat=true and exclude="12345" only 67890 will be use )
		Lower:    10, // have 10 lower case letters  
		Upper:    10, // have 10 upper case letters 
		Symbols:  10, // have 10 symbols
		NoRepeat: true, // do not repeat characters in password
		Cons:     true, // disallow consecutive characters ( abf,128,!@* etc)
		Exclude:  "@\"12345", // exclude list - what charcters to exlude 
	}
	p, err := f.Gen()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
}
```

```
&rndpass.Config{
		Length:   22,
		Numbers:  1, // have 10 numbers if possible ( with noRepeat=true and exclude="12345" only 67890 will be use )
		Lower:    1, // have 10 lower case letters
		Upper:    1, // have 10 upper case letters
		Symbols:  1,
		Exclude:  "@\"",
		NoRepeat: true,
		Cons:     true,
}
```
output:
```
wDe!H^Z3u+z&7Xc#T\Qm2}
3d{w}9O.I,u1qT_UaW(X&:
t.J9V%I5{xW<L6Mu}Nb\:!
[9!sH>aUjTh7]3xV_t0 Z%
+Pr;a5hV'J6j]zX.0c#Sw(
8i_7OgAw1,a:0rTf2Kt)Re
2;DnHi?N=sMuV:0wJ-1>%\
d3f6kT:yK}2)0uOo(7\I*>
!Wk%Hq2C5a1jO=z,Ph6V(<
+s&K:kRy7F3u0Jz|W%C}xj
Ze`r7tT uEo)h9{v}H8,X]
h+F|5~Q3z8g{I[E>Dc'7l*
U;aV3K8Tu]zNr.hCb^n<_}
m}M0bB)G>U.Y#7-d(Py%lq
M(1]j<2%6f^XyB!AaK:8hJ
4KeN|E.1]q Ha\pYfF7-u~
iYa>7Q_mG!e}u)0c2J&W9T
```


## Todo
- better readme
- docs
- more features?

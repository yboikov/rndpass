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

## Todo
- better readme
- docs
- more features?

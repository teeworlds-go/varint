# varint

varint is a simple variable-length integer encoding. It is a way to store integers in a space-efficient manner.
This variant of varint is space efficient for small integers and is used in the Teeworlds network protocol.

Additionally this linrary also provides functions that operate on 64bit integers which is out of scope of the Teeeworlds protocol.
These varants may be used for security research or other purposes.

```text
/ Format: ESDDDDDD EDDDDDDD EDD... Extended, Sign, Data,
// E: is next byte part of the current integer
// S: Sign of integer 0 = positive, 1 = negative
// Data, Integer bits that follow the sign
```

## Installation
```shell
// for latest tagged release
go get github.com/teeworlds-go/varint@latest

// for bleeding edge version
go get github.com/teeworlds-go/varint@master
```

## Example

```go
package main

import (
	"fmt"

	"github.com/teeworlds-go/varint"
)

func main() {
	buf := make([]byte, varint.MaxVarintLen32)
	written := varint.PutVarint(buf, 33)
	out, read := varint.Varint(buf)
	fmt.Println("written:", written)
	fmt.Println("read:", read)
	fmt.Println("value:", out)
	// Output:
	// written: 1
	// read: 1
	// value: 33
}
```


## BitArray for Golang [![Build Status](https://travis-ci.org/damnever/bitarray.svg?branch=master)](https://travis-ci.org/damnever/bitarray) [![GoDoc](https://godoc.org/github.com/damnever/bitarray?status.svg)](https://godoc.org/github.com/damnever/bitarray)

### Installation

```
go get github.com/damnever/bitarray
```

### Example

```go
package main

import "github.com/damnever/bitarray"

func main() {
    bits := bitarray.New(64)

    bits.Put(8, 1)  // set value of the bit to 1 by index 16
    bits.Get(8)     // get value of the bit by index 16
    bits.Set(9, 16) // set all bits to 1 between 9 and 16
    bits.Count()   // get the count of bit 1
    // Clear/Not/Eq/Leq/Lt/ToArray ...
}
```

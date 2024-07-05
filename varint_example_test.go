package varint_test

import (
	"fmt"

	"github.com/teeworlds-go/varint"
)

func Example() {
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

func ExamplePutVarint() {
	buf := make([]byte, varint.MaxVarintLen32)
	written := varint.PutVarint(buf, 63)
	fmt.Println(written)
	fmt.Printf("%b\n", buf[:written])
	// Output:
	// 1
	// [111111]
}

func ExampleVarint() {
	// 0b1xxxxxxx - extend bit set
	// 0bx0xxxxxx - positive sign
	buf := []byte{0b10111111, 0b00000001}
	out, read := varint.Varint(buf)
	fmt.Println("read:", read)
	fmt.Println("value:", out)
	// Output:
	// read: 2
	// value: 127
}

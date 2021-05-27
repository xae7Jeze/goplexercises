/*
Exercise 4.1: Write a function that counts the number of bits that are
different in two SHA256 hashes. (See PopCount from Section 2.6.2.)
*/
package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <string1> <string2>\n", os.Args[0])
		os.Exit(1)
	}
	c1 := sha256.Sum256([]byte(os.Args[1]))
	c2 := sha256.Sum256([]byte(os.Args[2]))
	fmt.Fprintf(os.Stdout, "c1: %x\nc2: %x\nDifferent bits: %v\n", c1, c2, DiffBits(&c1, &c2))
	os.Exit(0)
}

func DiffBits(x *[32]uint8, y *[32]uint8) int {
	var cnt int
	for i := 0; i < 32; i++ {
		for v := x[i] ^ y[i]; v != 0; v >>= 1 {
			cnt += int(v & 1)
		}
	}
	return cnt
}

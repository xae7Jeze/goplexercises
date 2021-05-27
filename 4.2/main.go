/*
Exercise 4.2: Write a program that prints the SHA256 hash of its
standard input by default but supports a command-line flag to print
the SHA384 or SHA512 hash instead.
*/
package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}
	var h hash.Hash
	switch os.Args[1] {
	case "sha256":
		h = sha256.New()
	case "sha384":
		h = sha512.New384()
	case "sha512":
		h = sha512.New()
	default:
		usage()
		os.Exit(1)
	}
	io.Copy(h, os.Stdin)
	fmt.Fprintf(os.Stdout, "%x\n", h.Sum(nil))
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %[1]s {sha256|sha384|sha512}\n"+
		"%[1]s reads from stdin and computes given hash\n",
		os.Args[0])
}

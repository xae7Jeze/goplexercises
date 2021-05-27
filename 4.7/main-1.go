package main

import (
	"errors"
	"fmt"
	"os"
	"unicode/utf8"
)

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}
	list := []byte(os.Args[1])
	if e := rreverse(list); e != nil {
		fmt.Fprintf(os.Stderr, "Oops: An Error occured: %v\n", e)
	}
	fmt.Fprintf(os.Stdout, "%s\n", list)
}

/*
* rotates l n positions, left if n < 0 and right if > 0
 */

func rotate(l []byte, n int) {
	var dir int8 = 1
	l_len := len(l)
	if l_len <= 1 {
		return
	}
	if n < 0 {
		dir = -1
		n *= -1
	}
	n = n % l_len
	if n == 0 {
		return
	}
	if dir > 0 {
		for i, j := 0, l_len-1; i < j; i, j = i+1, j-1 {
			l[i], l[j] = l[j], l[i]
		}
	}
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	for i, j := n, l_len-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	if dir < 0 {
		for i, j := 0, l_len-1; i < j; i, j = i+1, j-1 {
			l[i], l[j] = l[j], l[i]
		}
	}
}

/*
* reverses byte array with utf-8 chars by
* - read length of left utf char
* - rotate array lenght steps left to move first
*   character to end of array
* - shorten array by length
* - Stop if there is only one character
 */

func reverse(u []byte) error {
	var l int
	var r rune
	for uw := u; len(uw) > l; uw = uw[:len(uw)-l] {
		r, l = utf8.DecodeRune(uw[:])
		if r == utf8.RuneError && l == 0 {
			return errors.New("E_INVALID_UTF8")
		}
		rotate(uw, -1*l)
	}
	return nil
}

/*
* The same function using recursion
 */

func rreverse(u []byte) error {
	r, l := utf8.DecodeRune(u[:])
	if r == utf8.RuneError && l == 0 {
		return errors.New("E_INVALID_UTF8")
	}
	if len(u) <= l {
		return nil
	}
	rotate(u, -1*l)
	return rreverse(u[:len(u)-l])
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s <String>\n"+
		"Takes <String> and reverses it\n\n",
		os.Args[0])
}

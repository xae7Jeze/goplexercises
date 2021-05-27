/*
Exercise 4.7:
Modify reverse to reverse the characters of a []byte slice that represents a
UTF-8-encoded string , in place. Can you do it without allocating new memory?
*/
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
	fmt.Fprintf(os.Stdout, "%s => ", []byte(list))
	if e := reverse(list); e != nil {
		fmt.Fprintf(os.Stderr, "%s: OOPS: An error occured: %v\n", os.Args[0], e)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "%s\n", list)
}

func reverse(u []byte) error {
	var rns = make([]rune, 0, len(u))
	var i, j, l int
	var r rune
	for i, j = 0, 0; true; j++ {
		r, l = utf8.DecodeRune(u[i:])
		if l < 1 {
			break
		}
		if r == utf8.RuneError {
			return errors.New("E_INVALID_UTF8")
		}
		rns = append(rns, r)
		i += l
	}
	for i, j = 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	for i, j = 0, 0; j < len(rns); j++ {
		r, l = utf8.DecodeRuneInString(string(rns[j]))
		if r == utf8.RuneError {
			return errors.New("E_INVALID_UTF8")
		}
		utf8.EncodeRune(u[i:], r)
		i += l
	}
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s <String>\n"+
		"Takes <String> and reverses it\n\n",
		os.Args[0])
}

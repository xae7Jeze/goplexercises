/*
Exercise 4.6: Write an in-place function that squashes each run of adjacent
Unicode spaces (see unicode.IsSpace) in a UTF-8-encoded []byte slice into a
single ASCII space.
*/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

const (
	max_in = 1024 * 1024 * 1024
)

func main() {
	if len(os.Args) != 1 {
		usage()
		os.Exit(1)
	}
	buffer := make([]byte, 1024)
	list := make([]byte, 0)
	var err error
	var cnt int
	reader := bufio.NewReader(os.Stdin)
	for err != io.EOF {
		cnt, err = reader.Read(buffer)
		list = append(list, buffer[:cnt]...)
		if len(list) > max_in {
			err = errors.New("E_INPUT_SIZE_TOO_BIG")
			list = list[0:0]
			break
		}
	}
	if !(err == nil || err == io.EOF) {
		fmt.Fprintf(os.Stderr, "%s: An error occured: '%v'\n", os.Args[0], err)
		os.Exit(1)
	}
	list, err = squash_spcs(list)
	if err == nil {
		fmt.Printf("%s", list)
	} else {
		fmt.Fprintf(os.Stderr, "%s: An error occured: '%v'\n", os.Args[0], err)
	}
}

func squash_spcs(b []byte) ([]byte, error) {
	var l, i, sp_i, rm_sp int
	sp_i = -1
	var r rune
	for true {
		r, l = utf8.DecodeRune(b[i:])
		i += l
		if l < 1 {
			break
		}
		if r == utf8.RuneError {
			return b, errors.New("E_INVALID_UTF8")
		}
		if unicode.IsSpace(r) {
			rm_sp += l
			if sp_i < 0 {
				sp_i = i - l
			}
			continue
		}
		if sp_i >= 0 {
			utf8.EncodeRune(b[sp_i:], ' ')
			rm_sp--
			copy(b[sp_i+1:], b[i-l:])
			b = b[:(len(b) - rm_sp)]
			rm_sp = 0
			i = sp_i + 1
			sp_i = -1
		}
	}
	if sp_i >= 0 {
		utf8.EncodeRune(b[sp_i:], ' ')
		rm_sp--
	}
	b = b[:(len(b) - rm_sp)]
	return b, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s\n"+
		"Reads data from <Stdin>, squashes multiple utf8-ws to a single ascii space\n\n",
		os.Args[0])
}

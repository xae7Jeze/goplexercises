/*
Exercise 7.1: Using the ide as from ByteCounter, implement counters
for words and for lines. You will find bufio.ScanWords useful.
*/
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

var maxBytes = 100 * 1024 * 1024

type WordCounter int

func (w *WordCounter) Write(r io.Reader) (int, error) {
	in := bufio.NewScanner(r)
	in.Split(bufio.ScanWords)
	for in.Scan() {
		(*w)++
	}
	return int(*w), nil
}

type LineCounter int

func (l *LineCounter) Write(r io.Reader) (int, error) {
	in := bufio.NewScanner(r)
	in.Split(bufio.ScanLines)
	for in.Scan() {
		(*l)++
	}
	return int(*l), nil
}

func main() {
	b := make([]byte, maxBytes)
	r := bufio.NewReader(os.Stdin)
	rv := 0
	switch br, err := r.Read(b); {
	case err != nil:
		fmt.Printf("Ooops: %v\n", err)
		os.Exit(2)
	case br >= maxBytes:
		fmt.Printf("Ooops: Input is possibly truncated\n")
		rv = 1
	}
	buf := bytes.NewBuffer(b)
	var l LineCounter
	var w WordCounter
	l.Write(buf)
	fmt.Printf("Input contains %v lines\n", l)
	buf = bytes.NewBuffer(b)
	w.Write(buf)
	fmt.Printf("Input contains %v words\n", w)
	b = []byte{}
	os.Exit(rv)
}

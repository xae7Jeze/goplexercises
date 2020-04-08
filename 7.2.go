/*
Exercise 7.2: Write a function CountingWriter with the signature below that,
given an io.Writer , returns a new Writer that wraps the original, and a
pointer to an int64 variable that at any moment contains the number of bytes
written to the new Writer.
func CountingWriter(w io.Writer) (io.Writer, *int64)

Code based on https://github.com/4hel/gopl/blob/master/chap07/b_exercise-7.2/main.go
Couldn't find a basically different solution
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type counter struct {
	// 'count' counts all what will be written to 'writer'
	count  int64
	writer io.Writer
}

type c int64

func (cnt *c) Write(p []byte) (n int, err error) {
	c += len(p)
}

// make counter implement io.Writer
// original writer is in writer and sum up the bytes written to it
func (cnt *counter) Write(p []byte) (n int, err error) {
	n, err = cnt.writer.Write(p)
	cnt.count += int64(n)
	return n, err
}

// init new counter variable with original counter
// and return reference to it and its count element
func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var cnt counter
	cnt = counter{0, w}
	return &cnt, &cnt.count
}

func main() {
	cw, n := CountingWriter(os.Stdout)
	b := make([]byte, 1024)
	r := bufio.NewReader(os.Stdin)
	for {
		br, err := r.Read(b)
		if br == 0 && err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Ooops: %v\n", err)
			os.Exit(2)
		}
		fmt.Fprintf(cw, "%s", b[:br])
		fmt.Fprintf(os.Stderr, "%d\n", *n)
	}
	os.Exit(0)
}

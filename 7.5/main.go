/*
* Exercise 7.5: The LimitReader function in the io package accepts an io.Reader
* r and a number of bytes n , and returns another Reader that reads from r but
* reports an end-of-file condition after n bytes. Implement it.
* func LimitReader(r io.Reader, n int64) io.Reader
Read(p []byte) (n int, err error)
*/

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type readerN struct {
	r io.Reader
	n int64
}

func (lr *readerN) Read(b []byte) (n int, err error) {
	l := int64(len(b))
	err = io.EOF
	var e error = nil
	switch {
	case lr.n <= 0:
		n = 0
	case lr.n >= l:
		n, e = lr.r.Read(b)
	default:
		n, e = lr.r.Read(b[:lr.n])
	}
	if !(e == io.EOF || e == nil) {
		err = e
	}
	return

}
func LimitReader(r io.Reader, n int64) io.Reader {
	lr := &readerN{r, n}
	return lr
}
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <bytes_to_read>\n", os.Args[0])
		os.Exit(1)
	}
	n, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Oops: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	buf := make([]byte, 4096)
	r := LimitReader(os.Stdin, n)
	n64, err := r.Read(buf)
	fmt.Printf("%s", buf[:n64])
}

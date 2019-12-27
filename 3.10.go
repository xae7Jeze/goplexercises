/*
Exercise 3.10: Write a non-recursive version of comma,
using bytes.Buffer instead of string concatenation.
*/
package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%v: Ooops - missing arg\n", os.Args[0])
		os.Exit(1)
	}
	var err error
	var s string
	if s, err = comma(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "%v: Ooops: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", s)
}

func comma(s string) (string, error) {
	var b bytes.Buffer
	var s1, s2, sign, dot string
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return s, err
	}
	if len(s) < 1 {
		return s, errors.New("NaN")
	}
	if s[0] == '-' || s[0] == '+' {
		sign = s[0:1]
		s = s[1:]
	}
	if strings.IndexByte(s, '.') >= 0 {
		sa := strings.Split(s, ".")
		s1, s2 = sa[0], sa[1]
		dot = "."
	} else if i := strings.IndexByte(strings.ToLower(s), 'e'); i >= 0 {
		s1 = s[:i]
		s2 = s[i:]
	} else {
		s1 = s
	}
	l := len(s1)
	if l < 4 {
		return (sign + s1 + dot + s2), nil
	}
	if fc := l % 3; fc > 0 {
		b.WriteString(s1[:fc])
		b.WriteByte(',')
		s1 = s1[fc:]
		l -= fc
	}
	for ; l > 3; l -= 3 {
		b.WriteString(s1[0:3])
		s1 = s1[3:]
		b.WriteByte(',')
	}
	b.WriteString(s1[0:])
	return sign + b.String() + dot + s2, nil
}

func c_book(s string) string {
	n := len(s)
	if n < 4 {
		return s
	}
	return c_book(s[:n-3]) + "," + s[n-3:]
}

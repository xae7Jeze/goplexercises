/*
Exercise 4.5: Write an in-place function to eliminate adjacent duplicates
in a []string slice.
*/
package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}
	s, err := strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
		os.Exit(1)
	}
	list := make([]string, 0, 100)
	for _, arg := range os.Args[2:] {
		list = append(list, arg)
	}
	fmt.Fprintf(os.Stdout, "%v => ", list)
	rotate(list, s)
	fmt.Fprintf(os.Stdout, "%v\n", list)
}

func rotate(l []string, n int) {
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

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s <steps> Arg1 ... ArgN*\n"+
		"   <steps> : INT: If < 0 rotates <steps> left\n"+
		"                  If > 0 rotates <steps> right\n\n"+
		"* The remaining args building the array to rotate\n\n",
		os.Args[0])
}

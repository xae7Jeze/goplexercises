/*
Exercise 2.5: The expression x&(x-1) clears the rightmost non-zero bit of x.
Write a version of PopCount that counts bits by using this fact, and assess
its performance.
*/
package main

import (
	"fmt"
	"os"
	"popcount"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: <%s> <list_of_unsigned_integers>\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Printf("%-20s %-8s %-s\n", "Number", "PopCount", "PopCountR")
	for _, arg := range os.Args[1:] {
		if u, err := strconv.ParseUint(arg, 10, 64); err == nil {
			fmt.Printf("%-20d %-8d %-d\n", u, popcount.PopCount(u), popcount.PopCountR(u))
		} else {
			fmt.Fprintf(os.Stderr, "Invalid arg: '%s'\n", arg)
		}
	}
}

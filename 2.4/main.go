/*
Exercise 2.4: Write a version of PopCount that counts bits by shifting
its argument through 64 bit position s, testing the rightmost bit each
time. Compare its performance to the table-lookup version.
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
	fmt.Printf("%-20s %-8s %-s\n", "Number", "PopCount", "PopCountS")
	for _, arg := range os.Args[1:] {
		if u, err := strconv.ParseUint(arg, 10, 64); err == nil {
			fmt.Printf("%-20d %-8d %-d\n", u, popcount.PopCount(u), popcount.PopCountS(u))
		} else {
			fmt.Fprintf(os.Stderr, "Invalid arg: '%s'\n", arg)
		}
	}
	os.Exit(0)
}

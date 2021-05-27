/*
Exercise 2.3: Rewrite PopCount to use a loop instead of a single expression.
Compare the performance of the two versions. (Section 11.4 shows how to
compare the perfor mance of different implementations systematically.)
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
	fmt.Printf("%-20s %-8s %-s\n", "Number", "PopCount", "PopCountL")
	for _, arg := range os.Args[1:] {
		if u, err := strconv.ParseUint(arg, 10, 64); err == nil {
			fmt.Printf("%-20d %-8d %-d\n", u, popcount.PopCount(u), popcount.PopCountL(u))
		} else {
			fmt.Fprintf(os.Stderr, "Invalid arg: '%s'\n", arg)
		}
	}
}

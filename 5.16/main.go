/*
Exercise 5.16: Write a variadic version of strings.Join .
*/
package main

import (
	"fmt"
	"os"
)

func join(sep string, elements ...string) string {
	result := ""
	if len(elements) < 1 {
		return result
	}
	result = elements[0]
	if len(elements) == 1 {
		return result
	}
	for _, s := range elements[1:] {
		result += sep + s
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	if len(os.Args) == 2 {
		fmt.Printf("\n")
	} else {
		fmt.Printf("%s\n", join(os.Args[1], os.Args[2:]...))
	}
	os.Exit(0)

}

func usage() {
	fmt.Fprintf(os.Stderr,
		"\nUsage: %s <separator> <strings to join with separator>\n",
		os.Args[0])
}

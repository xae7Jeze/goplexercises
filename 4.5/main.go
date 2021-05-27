/*
Exercise 4.4: Write a version of rotate that operates in a single pass.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		usage()
		os.Exit(1)
	}
	list := make([]string, 0, 100)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: An error occured: '%v'\n", os.Args[0], err)
	}
	list = remove_dups(list)
	for _, s := range list {
		fmt.Fprintf(os.Stdout, "%s\n", s)
	}
}

func remove_dups(l []string) []string {
	l_len := len(l)
	if l_len < 1 {
		return l
	}
	for i, cs := 1, l[0]; i < l_len; {
		if l[i] == cs {
			copy(l[i:l_len-1], l[i+1:l_len])
			l_len--
			l = l[:l_len]
		} else {
			cs = l[i]
			i++
		}
	}
	return l
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s\n"+
		"Reads lines from <Stdin>, removes adjacent dups and prints them\n\n",
		os.Args[0])
}

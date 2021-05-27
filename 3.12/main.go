/*
Exercise 3.12: Write a function that reports whether two strings are
anagrams of each other, that is, they contain the same letters in a
different order.
*/
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type RuneSlice []rune

func (p RuneSlice) Len() int           { return len(p) }
func (p RuneSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p RuneSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %v <string1> <string2>\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Printf("Anagram ? %s <=> %s -> %v\n", os.Args[1], os.Args[2], is_ana(os.Args[1], os.Args[2]))
}

func is_ana(s1, s2 string) bool {
	l1 := len(s1)
	l2 := len(s2)
	if l1 != l2 {
		return false
	}
	if l1 == 0 {
		return true
	}
	rns1 := make([]rune, 0, len(s1))
	rns2 := make([]rune, 0, len(s1))
	for _, r := range s1 {
		rns1 = append(rns1, r)
	}
	for _, r := range s2 {
		rns2 = append(rns2, r)
	}
	sort.Sort(RuneSlice(rns1))
	sort.Sort(RuneSlice(rns2))

	return strings.ToLower(string(rns1)) == strings.ToLower(string(rns2))
}

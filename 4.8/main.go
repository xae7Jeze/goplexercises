/*
Exercise 4.8:
Modify charcount to count letters, digits, and so on in their Unicode categories,
using functions like unicode.IsLetter .
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"unicode"
)

func main() {
	if len(os.Args) != 1 {
		usage()
		os.Exit(1)
	}
	cats := make(map[string]int)
	total := 0
	cats["Invalid"] = 0
	in := bufio.NewReader(os.Stdin)
	for {
		r, l, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Oops: An Error occured: %v\n", err)
			os.Exit(1)
		}
		total++
		if r == unicode.ReplacementChar && l == 1 {
			cats["Invalid"]++
			continue
		}
		if unicode.IsControl(r) {
			cats["Controls"]++
			continue
		}
		if unicode.IsDigit(r) {
			cats["Digits"]++
			continue
		}
		if unicode.IsLetter(r) {
			cats["Letters"]++
			continue
		}
		if unicode.IsMark(r) {
			cats["Marks"]++
			continue
		}
		if unicode.IsNumber(r) {
			cats["Numbers"]++
			continue
		}
		if unicode.IsPunct(r) {
			cats["Puncts"]++
			continue
		}
		if unicode.IsSpace(r) {
			cats["Spaces"]++
			continue
		}
		if unicode.IsSymbol(r) {
			cats["Symbols"]++
			continue
		}
		cats["Other"]++
	}
	keys := []string{}
	ml := 0
	for k, _ := range cats {
		keys = append(keys, k)
		if len(k) > ml {
			ml = len(k)
		}
	}
	ml++
	ll := ml * len(keys)
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%[1]*[2]s", ml, k)
	}
	fmt.Printf("\n")
	for i := 0; i < ll; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("\n")
	for _, k := range keys {
		fmt.Printf("%[1]*[2]v", ml, cats[k])
	}
	fmt.Printf("\n\nTotal: %v (%v)\n\n", total, sumup(cats))

}

func sumup(m map[string]int) int {
	s := 0
	for _, v := range m {
		s += v
	}
	return (s)
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s\n"+
		"Reads from Stdin and prints statistic of chars read.\n\n",
		os.Args[0])
}

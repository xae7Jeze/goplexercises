/*
Exercise 4.9:
Write a program wordfreq to report the frequency of each word in an input text
file. Call input.Split(bufio.ScanWords) before the first call to Scan to break the input into
words instead of lines.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	if len(os.Args) != 1 {
		usage()
		os.Exit(1)
	}
	words := make(map[string]int)
	total, ml := 0, 0
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanWords)
	for in.Scan() {
		total++
		w := in.Text()
		words[w]++
		if len(w) > ml {
			ml = len(w)
		}
	}
	if err := in.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops, an error occured: %v\n", err)
		os.Exit(2)
	}
	if ml > 60 {
		ml = 60
	}

	for _, v := range SortByValsN(words, true) {
		fmt.Printf("%-[1]*[2]v %[3]v\n", ml, v.word, v.count)
	}
	fmt.Printf("\n\nTotal: %v\n\n", total)
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s\n"+
		"Reads from Stdin and prints statistic of chars read.\n\n",
		os.Args[0])
}

type scnts struct {
	word  string
	count int
}

type scnts_a []scnts

func (s scnts_a) Len() int {
	return len(s)
}
func (s scnts_a) Less(i, j int) bool {
	return s[i].count < s[j].count
}
func (s scnts_a) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func SortByValsN(m map[string]int, reverse bool) scnts_a {
	vals := make(scnts_a, 0, len(m))
	for k, v := range m {
		vals = append(vals, scnts{word: k, count: v})
	}
	if reverse {
		sort.Sort(sort.Reverse(vals))
	} else {
		sort.Sort(vals)
	}
	return vals
}

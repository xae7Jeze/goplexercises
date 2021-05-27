/*
Exercise 1.4: Modify dup2 to print the names of all files in which each
duplicated line occurs.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	fnames := make(map[string]map[string]bool)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts, fnames)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, fnames)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%-3d (Found in %v): `%s`\n", n, mapKeys(fnames[line], ", "), line)
		}
	}
}

func countLines(f *os.File, counts map[string]int,
	fnames map[string]map[string]bool) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		txt := input.Text()
		fname := f.Name()
		counts[txt]++
		if fnames[txt] == nil {
			fnames[txt] = make(map[string]bool)
		}
		fnames[txt][fname] = true
	}
}

func mapKeys(m map[string]bool, s string) string {
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, s)
}

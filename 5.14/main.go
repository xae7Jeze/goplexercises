/*
Exercise 5.14: Use the breadthFirst function to explore a different structure.
For example, you could use the course dependencies from the topoSort example
(a directed graph), the file system hierarchy on your computer (a tree), or a
list of bus or subway routes downlo ade d from your city government’s website
(an undirected graph).
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

//!+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) ([]string, error), worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				if wl, err := f(item); err != nil {
					log.Printf("ERROR: %v\n", err.Error())
				} else {
					worklist = append(worklist, wl...)
				}
			}
		}
	}
}

//!-breadthFirst

//!+crawl
func crawl(start string) ([]string, error) {
	start = strings.TrimRight(start, string(os.PathSeparator))
	stat, err := os.Lstat(start)
	if err != nil {
		return []string{}, err
	}
	fmt.Println(start)
	if stat.IsDir() == false {
		return []string{}, nil
	}
	f, err := os.Open(start)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()
	list, err := f.Readdirnames(-1)
	if err != nil {
		return []string{}, err
	}
	for i, n := range list {
		if start == "/" {
			list[i] = start + n
		} else {
			list[i] = start + string(os.PathSeparator) + n
		}
	}
	return list, nil
}

//!-crawl

//!+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}

//!-main

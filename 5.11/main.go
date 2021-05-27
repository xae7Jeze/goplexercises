/*
Exercise 5.11: The instructor of the linear algebra course decides that
calculus is now a prerequisite. Extend the topoSort function to report cycles.
*/
package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type depMap map[string][]string

type pathTrace struct {
	np   map[string]bool
	path []*string
}

//!+table
// prereqs maps computer science courses to their prerequisites.
var prereqs = depMap{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math", "compilers"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

//!-table

//!+main
func main() {
	if l, err := topoSort(prereqs); err != nil {
		fmt.Fprintf(os.Stderr, "%s: An error occured: %v\n", os.Args[0], err)
		os.Exit(1)
	} else {
		for i, course := range l {
			fmt.Printf("%d:\t%s\n", i+1, course)
			os.Exit(0)
		}
	}
}

func TestTopoSort(t *testing.T) {
	oM := make(map[string]int)
	sortetCourses, _ := topoSort(prereqs)
	for i, course := range sortetCourses {
		fmt.Printf("%d:\t%s. Depends on: %q\n", i, course, prereqs[course])
		oM[course] = i
	}
	for _, course := range sortetCourses {
		for _, pC := range prereqs[course] {
			if oM[pC] > oM[course] {
				t.Errorf("Order invalid: %v depends on %v", course, pC)
			}
		}
	}
}

func topoSort(m depMap) ([]string, error) {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(string, pathTrace) error
	var pt pathTrace
	visitAll = func(item string, pt pathTrace) error {
		if _, exists := pt.np[item]; exists == true {
			pt.path = append(pt.path, &item)
			path := []string{}
			var begin bool = false
			for _, s := range pt.path {
				if *s == item {
					begin = true
				}
				if begin {
					path = append(path, *s)
				}
			}
			return fmt.Errorf("Loop detected: %s\n", strings.Join(path, " -> "))
		} else {
			pt.np[item] = true
			pt.path = append(pt.path, &item)
		}

		if seen[item] {
			return nil
		}
		seen[item] = true
		if s, ok := m[item]; ok != true {
			order = append(order, item)
			return nil
		} else {
			for _, i := range s {
				if err := visitAll(i, pt); err != nil {
					return err
				}
			}
		}
		order = append(order, item)
		return nil
	}
	/*
		var keys []string
		for key := range m {
			keys = append(keys, key)
		}

		sort.Strings(keys)
	*/
	for key := range m {
		pt = *(new(pathTrace))
		pt.np = make(map[string]bool)
		if err := visitAll(key, pt); err != nil {
			return []string{}, err
		}
	}
	return order, nil
}

//!-main

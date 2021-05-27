/*
Exercise 5.10: Rewrite topoSort to use maps instead of slices and eliminate
the initial sort. Verify that the results, though nondeterministic, are valid
topological orderings.
*/
package main

import (
	"fmt"
	"testing"
	//	"sort"
)

type depMap map[string][]string

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

	"data structures":       {"discrete math"},
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
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func TestTopoSort(t *testing.T) {
	oM := make(map[string]int)
	sortetCourses := topoSort(prereqs)
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

func topoSort(m depMap) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(string)

	visitAll = func(item string) {
		if seen[item] {
			return
		}
		seen[item] = true
		if s, ok := m[item]; ok != true {
			order = append(order, item)
			return
		} else {
			for _, item := range s {
				visitAll(item)
			}
		}
		order = append(order, item)
	}
	/*
		var keys []string
		for key := range m {
			keys = append(keys, key)
		}

		sort.Strings(keys)
	*/
	for key := range m {
		visitAll(key)
	}
	return order
}

//!-main

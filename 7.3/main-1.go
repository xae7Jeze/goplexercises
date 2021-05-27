/*
Exercise 7.3: Write a String method for the *tree type in gopl.io/ch4/treesort
(§4.4) that reveals the sequence of values in the tree.
*/
package main

import (
	"fmt"
	"math/rand"
	"os"
	//"strings"
	"sort"
	"time"
)

type tree struct {
	value       int
	left, right *tree
}

type treeWithPd struct {
	tree
	predessor *tree
}

type t2s struct {
	s       string
	sortkey int
}
type t2sA []t2s

const (
	L = -1
	M = 0
	R = 1
)

func (items t2sA) Len() int {
	return len(items)
}
func (items t2sA) Less(i, j int) bool {
	return items[i].sortkey < items[j].sortkey
}
func (items t2sA) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (t *tree) String() string {
	//values := []string{}
	depth := 0
	index := 0
	var ts t2sA
	var f func(t *tree, ts *t2sA, lr int)
	s := ""
	f = func(t *tree, ts *t2sA, lr int) {
		depth++
		if t == nil {
			depth--
			return
		}
		index++
		if t.left != nil {
			f(t.left, ts, L)
		}
		pfx := ""
		for i := 0; i < depth; i++ {
			pfx += "-"
		}
		if pfx == "-" {
			pfx = "*"
		}
		*ts = append(*ts, t2s{s: fmt.Sprintf("%-8s%3d\n", pfx, t.value), sortkey: lr * depth})
		//*s += fmt.Sprintf("%c %-8s%3d\n", lr, pfx, t.value)
		if t.right != nil {
			f(t.right, ts, R)
		}
		depth--
		return
	}
	f(t, &ts, M)
	/*for _, v := range v {
		fmt.Fprintf(os.Stderr, "%s\n", v)
	}*/
	sort.Stable(t2sA(ts))
	for _, i := range ts {
		s += i.s
	}
	return fmt.Sprintf("%+v", s)
}

func add(t *tree, value int) *tree {
	if t == nil {
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

func main() {
	var t *tree = nil
	data := make([]int, 10)
	rand.Seed(time.Now().UnixNano())
	for i := range data {
		data[i] = rand.Int() % 100
		if t == nil {
			t = add(t, data[i])
		} else {
			add(t, data[i])
		}
	}
	fmt.Fprintf(os.Stderr, "%v\n", data)
	fmt.Fprintf(os.Stdout, "%s\n", t.String())
	//t.String()
}

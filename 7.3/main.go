/*
Exercise 7.3: Write a String method for the *tree type in gopl.io/ch4/treesort
(§4.4) that reveals the sequence of values in the tree.
*/
package main

import (
	"fmt"
	"math/rand"
	"os"
	//"sort"
	"time"
)

type tree struct {
	value       int
	left, right *tree
}

var (
	tBlack     = fmt.Sprintf("%c%c%c%c%c", 033, 0133, 063, 060, 0155)
	tRed       = fmt.Sprintf("%c%c%c%c%c", 033, 0133, 063, 061, 0155)
	tGreen     = fmt.Sprintf("%c%c%c%c%c", 033, 0133, 063, 062, 0155)
	tYellow    = fmt.Sprintf("%c%c%c%c%c", 033, 0133, 063, 063, 0155)
	tBlue      = fmt.Sprintf("%c%c%c%c%c", 033, 0133, 063, 064, 0155)
	tMagenta   = fmt.Sprintf("%c%c%c%c%c", 033, 0133, 063, 065, 0155)
	tCyan      = fmt.Sprintf("%c%c%c%c%c", 033, 0133, 063, 066, 0155)
	tWhite     = fmt.Sprintf("%c%c%c%c%c", 033, 0133, 063, 067, 0155)
	tUnderline = fmt.Sprintf("%c%c%c%c", 033, 0133, 064, 0155)
	tBold      = fmt.Sprintf("%c%c%c%c", 033, 0133, 061, 0155)
	tReset     = fmt.Sprintf("%c%c%c%c", 033, 0133, 060, 0155)
)

func (t *tree) elementString(indent int) string {
	la := t.left
	lv := 0
	indentS := fmt.Sprintf("%*s", indent, "")
	if la != nil {
		lv = la.value
	}
	ra := t.right
	rv := 0
	if ra != nil {
		rv = ra.value
	}
	me := fmt.Sprintf(" %3v (%010p) ", t.value, t)
	l := fmt.Sprintf(" %3v (%010p) ", lv, la)
	r := fmt.Sprintf(" %3v (%010p) ", rv, ra)
	el := len(me)
	line := "+"
	for i := 0; i <= 2*el; i++ {
		line += "-"
	}
	line += "+"
	s := indentS + line + "\n"
	s += indentS + "|" + fmt.Sprintf("%*s%s%*s", el/2, "", me, (el/2)+1, "") + "|\n"
	s += indentS + line + "\n"
	s += indentS + "|" + l + "|" + r + "|\n"
	s += indentS + line
	//fmt.Fprintf(os.Stderr, "\nDEBUG:\n%s\n", s)
	return s
}

func (t *tree) String() string {
	var f func(t *tree)
	s := ""
	m := make(map[*tree][2]int)
	minx := 0
	maxx, maxy := 0, 0
	startspread := 8
	depth := -1
	f = func(t *tree) {
		depth++
		if m[t][0] < minx {
			minx = m[t][0]
		}
		if m[t][0] > maxx {
			maxx = m[t][0]
		}
		if m[t][1] > maxy {
			maxy = m[t][1]
		}
		spread := 1
		if depth < startspread {
			spread = startspread - depth
		}
		if t.left != nil {
			m[t.left] = [2]int{(m[t][0] - spread), m[t][1] + 1}
			f(t.left)
		}
		if t.right != nil {
			m[t.right] = [2]int{(m[t][0] + spread), m[t][1] + 1}
			f(t.right)
		}
		depth--
		return
	}
	m[t] = [2]int{0, 0}
	f(t)
	if minx < 0 {
		for k, _ := range m {
			m[k] = [2]int{m[k][0] - minx, m[k][1]}
		}
		maxx += (-1 * minx)
	}
	mr := make(map[[2]int][]int)
	for k, v := range m {
		mr[v] = append(mr[v], k.value)
	}
	ml := 0
	for _, v := range m {
		if l := len(v); l > ml {
			ml = l
		}
	}
	for col := 0; col <= maxx; col++ {
		s += fmt.Sprintf("____")
	}
	s += fmt.Sprintf("\n")
	for row := 0; row <= maxy; row++ {
		for col := 0; col <= maxx; col++ {
			if v, ok := mr[[2]int{col, row}]; ok == true {
				sTmp := "[__]"
				if l := len(v); l > 1 {
					for i := 0; i < l; i++ {
						if i == l-1 {
							sTmp = ""
						}
						s += fmt.Sprintf("[%[2]s%2[1]d%[3]s]%[4]s", v[i], tRed, tReset, sTmp)
					}
				} else {
					s += fmt.Sprintf("[%[2]s%2[1]d%[3]s]", v[0], tRed, tReset)
				}
			} else {
				s += fmt.Sprintf("[__]")
			}
		}
		s += fmt.Sprintf("\n")
	}
	return s
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
	fmt.Fprintf(os.Stdout, "\nRandom input data: %v gives unbalanced tree:\n\n", data)
	fmt.Fprintf(os.Stdout, "%s\n", t.String())
}

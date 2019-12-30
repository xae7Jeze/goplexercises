/*
Exercise 6.1: Implement these additional methods:
func (*IntSet) Len() int // return the number of elements
func (*IntSet) Remove(x int) // remove x from the set
func (*IntSet) Clear() // remove all elements from the set
func (*IntSet) Copy() *IntSet // return a copy of the set
*/

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var x IntSet
	var n []int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		n = append(n, rand.Intn(10000))
	}
	fmt.Printf("Adding %v to %v = ", n, x.String())
	x.AddAll(n...)
	fmt.Printf("%v\n", x.String())
}

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Using algo from lesson 2.5: (x & (x - 1) clears leftmost bit
func (s *IntSet) Len() uint {
	var r uint
	for _, x := range s.words {
		for x != 0 {
			x = x & (x - 1)
			r++
		}
	}
	return r
}

// Removes the non-negative value x to the set.
// In addition it cleans up unneeded words before
// returning
func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)
	defer func() {
		for (len(s.words) > 0) && (s.words[len(s.words)-1] == 0) {
			s.words = s.words[:len(s.words)-1]
		}
	}()
	if len(s.words) <= word {
		return
	}
	s.words[word] &= ^(1 << bit)
}

// clears set by setting to empty set
func (s *IntSet) Clear() {
	s.words = []uint64{}
}

// Copy
func (s *IntSet) Copy() *IntSet {
	c := new(IntSet)
	c.words = make([]uint64, len(s.words))
	copy(c.words, s.words)
	return c
}

// AddAll adds a set of values to the set.
func (s *IntSet) AddAll(x ...int) {
	for _, x := range x {
		word, bit := x/64, uint(x%64)
		for word >= len(s.words) {
			s.words = append(s.words, 0)
		}
		s.words[word] |= 1 << bit
	}
}

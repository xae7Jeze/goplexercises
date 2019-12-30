/*
Exercise 6.3: (*IntSet).UnionWith computes the union of two sets using |,
the word-parallel bitwise OR operator. Implement methods for IntersectWith,
DifferenceWith , and SymmetricDifference for the corresponding set operations.
(The symmetr ic dif ference of two sets contains the elements present in one
set or the other but not both.)
*/

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

func printRed(s string) {
	fmt.Printf("%c%c%c%c%c%s%c%c%c%c",
		033, 0133, 063, 061, 0155,
		s,
		033, 0133, 060, 0155)
}

func main() {
	x := &IntSet{}
	y := &IntSet{}
	var n []int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		n = append(n, rand.Intn(1000))
	}
	x.AddAll(n...)
	for i := 0; i < 5; i++ {
		n[i] = rand.Intn(1000)
	}
	y.AddAll(n...)
	n = n[:3]
	for i := 0; i < 3; i++ {
		n[i] = rand.Intn(1000)
	}
	fmt.Printf("X: %v\nY; %v\nCommon Elements: %v\n", x.String(), y.String(), n)
	x.AddAll(n...)
	y.AddAll(n...)
	y.AddAll(9064)
	xcp := x.Copy()
	fmt.Printf("%v \u2229 %v = ", x.String(), y.String())
	x.IntersectWith(y)
	fmt.Printf("%v\n", x.String())
	x = xcp.Copy()
	fmt.Printf("%v \u2216 %v = ", x.String(), y.String())
	x.DifferenceWith(y)
	fmt.Printf("%v\n", x.String())
	x = xcp.Copy()
	fmt.Printf("%v \u2206 %v = ", x.String(), y.String())
	x.SymmetricDifference(y)
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

// IntersectWith sets s to the intersect of s and t.
func (s *IntSet) IntersectWith(t *IntSet) {
	if len(s.words) > len(t.words) {
		s.words = s.words[:len(t.words)]
	}
	for i, _ := range s.words {
		s.words[i] &= t.words[i]
	}
}

// DifferentWith, sets s to s \ t.
func (s *IntSet) DifferenceWith(t *IntSet) {
	max := len(s.words)
	if len(t.words) < max {
		max = len(t.words)
	}
	for i := 0; i < max; i++ {
		s.words[i] &= ^t.words[i]
	}
}

// SymmetricDifferent sets s to s ∆ t
func (s *IntSet) SymmetricDifference(t *IntSet) {
	max := len(s.words)
	if len(t.words) < max {
		max = len(t.words)
	}
	for i := 0; i < max; i++ {
		s.words[i] ^= t.words[i]
	}
	if len(t.words) > max {
		s.words = append(s.words, t.words[max:]...)
	}
}

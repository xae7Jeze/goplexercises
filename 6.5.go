/*
Exercise 6.5: The type of each word used by IntSet is uint64 , but 64-bit arithmetic may be
inefficient on a 32-bit platform. Modify the program to use the uint type, which is the most
efficient unsigned integer type for the platform. Instead of dividing by 64, define a constant
holding the effective size of uint in bits, 32 or 64. You can use the perhaps too clever expres-
sion 32 << (^uint(0) >> 63) for this purpose.
*/

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

const (
	wordLen = 32 << (^uint(0) >> 63)
)

func main() {
	x := &IntSet{}
	var n []int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		n = append(n, rand.Intn(1<<16))
	}
	x.AddAll(n...)
	fmt.Printf("IntSet presented as:\n")
	fmt.Printf("String: %v\n", x.String())
	fmt.Printf("[]uint: %v\n", x.Elems())
}

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/wordLen, uint(x%wordLen)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/wordLen, uint(x%wordLen)
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
		for j := 0; j < wordLen; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", wordLen*i+j)
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
	word, bit := x/wordLen, uint(x%wordLen)
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
	s.words = []uint{}
}

// Copy
func (s *IntSet) Copy() *IntSet {
	c := new(IntSet)
	c.words = make([]uint, len(s.words))
	copy(c.words, s.words)
	return c
}

// AddAll adds a set of values to the set.
func (s *IntSet) AddAll(x ...int) {
	for _, x := range x {
		word, bit := x/wordLen, uint(x%wordLen)
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

// Elems: returns elemants of IntSet as slice of uints
func (s *IntSet) Elems() []uint {
	r := []uint{}
	for i, w := range s.words {
		for b := uint(0); b < wordLen; b++ {
			if (w & (1 << uint(b))) > 0 {
				r = append(r, uint(i)*wordLen+b)
			}
		}
	}
	return r
}

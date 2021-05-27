/*
Exercise 4.3: Rewrite reverse to use an array pointer instead of a slice.
*/
package main

import (
	"fmt"
	"os"
)

func main() {
	list := []int{1, 2, 3, 4, 5}
	fmt.Fprintf(os.Stdout, "%v => ", list)
	reverse(&list)
	fmt.Fprintf(os.Stderr, "%v\n", list)
}

func reverse(l *[]int) {
	for i, j := 0, len(*l)-1; i < j; i, j = i+1, j-1 {
		(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
	}
}

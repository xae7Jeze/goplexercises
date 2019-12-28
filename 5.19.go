/*
Exercise 5.18: Without changing its behavior, rewrite the fetch function
to use defer to close the writable file.
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func panicAndRecoverSquared(x int) (q int) {
	defer func() {
		switch p := recover(); p {
		case nil:
		case x:
			q = x * x
		default:
			panic(p)
		}
	}()
	panic(x)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(42)
	fmt.Printf("%v^2 = %v\n", x, panicAndRecoverSquared(x))
}

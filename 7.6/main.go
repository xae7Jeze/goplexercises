/*
* Exercise 7.6: Add support for Kelvin temperatures to tempflag.
 */
package main

import (
	"flag"
	"fmt"
	"tempconv"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}

//!-

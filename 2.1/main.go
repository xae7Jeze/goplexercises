/*
Exercise 2.1: Add types, constants, and functions to tempconv for
processing temperatures in the Kelvin scale, where zero Kelvin is
−273.15°C and a difference of 1K has the same magnitude as 1°C.
*/
package main

import (
	"fmt"
	"tempconv"
)

func main() {
	var k tempconv.Kelvin
	fmt.Printf("%8s %10s %8s\n", "Celsius", "Fahrenheit", "Kelvin")
	for k = 0; k <= tempconv.Kelvin(tempconv.AbsZeroC*(-1)+100); k += 1 {
		fmt.Printf("%8.2f %10.2f %8.2f\n", tempconv.KtoC(k), tempconv.KtoF(k), k)
	}
}

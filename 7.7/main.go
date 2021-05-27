/*
* Exercise 7.7: Explain why the help message contains °C when the default
* value of 20.0 does not.
 */

package main

import (
	"fmt"
)

func main() {
	question :=
		`
Q: Explain why the help message contains °C when the default
value of 20.0 does not.
`
	answer :=
		`
A: The help message contains the unit, because it's output is formatted by
CelsiusFlag's String method, inherited from type Celsius
`
	fmt.Printf("%s\n%s\n", question, answer)
}

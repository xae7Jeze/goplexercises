/*
Exercise 2.2: Write a general-purpose unit-conversion program analogous
to cf that reads numbers from its command-line arguments or from the
standard input if there are no arguments, and converts each number into
units like temperature in Celsius and Fahrenheit, lengt h in feet and
meters, weight in pounds and kilograms, and the like.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	AbsZeroC = -273.15
	P2C      = 0.45359237
	F2M      = 0.3048
)

func main() {
	var readFromStdin = false
	var mode = ""
	switch len(os.Args) {
	case 1:
		usage()
		os.Exit(1)
	case 2:
		readFromStdin = true
	}
	mode = strings.ToLower(os.Args[1])
	if readFromStdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if u, err := strconv.ParseFloat(scanner.Text(), 64); err != nil {
				fmt.Fprintf(os.Stderr, "%s: Invalid  input: '%s'\n", os.Args[0], scanner.Text())
				os.Exit(2)
			} else {
				dispatch(mode, u)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error reading standard input: %v\n", os.Args[0], err)
			os.Exit(2)
		}
	} else {
		for _, arg := range os.Args[2:] {
			if u, err := strconv.ParseFloat(arg, 64); err != nil {
				fmt.Fprintf(os.Stderr, "%s: Invalid  input: '%s'\n", os.Args[0], arg)
				os.Exit(2)
			} else {
				dispatch(mode, u)
			}
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <INPUT_UNIT> N1 N2 N3\n"+
		"Usage: %s <INPUT_UNIT> [READ FROM STDIN]\n"+
		"<INPUT_UNIT> must be one of 'kelvin', 'celsius','fahrenheit',"+
		"'feet','meter','pound','kilo'\n", os.Args[0], os.Args[0])
}

func dispatch(mode string, u float64) {
	switch mode {
	case "kelvin":
		Kelvin(u)
	case "celsius":
		Celsius(u)
	case "fahrenheit":
		Fahrenheit(u)
	case "feet":
		Feet(u)
	case "meter":
		Meter(u)
	case "pound":
		Pound(u)
	case "kilo":
		Kilo(u)
	default:
		usage()
		os.Exit(1)
	}
}

func CToK(c float64) float64 {
	return (c - AbsZeroC)
}

func CToF(c float64) float64 {
	return (c*9/5 + 32)
}

func Celsius(c float64) {
	fmt.Fprintf(os.Stdout, "%6.2f° C = %6.2f° K = %6.2f° F\n", c, CToK(c), CToF(c))
}

func KToC(k float64) float64 {
	return (k + AbsZeroC)
}

func KToF(k float64) float64 {
	return CToF(KToC(k))
}

func Kelvin(k float64) {
	fmt.Fprintf(os.Stdout, "%6.2f° K = %6.2f° C = %6.2f° F\n", k, KToC(k), KToF(k))
}

func FToC(f float64) float64 {
	return ((f - 32) * 5 / 9)
}

func FToK(f float64) float64 {
	return CToK(FToC(f))
}

func Fahrenheit(f float64) {
	fmt.Fprintf(os.Stdout, "%6.2f° F = %6.2f° K = %6.2f° C\n", f, FToK(f), FToC(f))
}

func PToK(p float64) float64 {
	return (P2C * p)
}
func Pound(p float64) {
	fmt.Fprintf(os.Stdout, "%6.2f P = %6.2f K\n", p, PToK(p))
}
func KToP(k float64) float64 {
	return (k / P2C)
}

func Kilo(k float64) {
	fmt.Fprintf(os.Stdout, "%6.2f K = %6.2f P\n", k, KToP(k))
}

func MToF(m float64) float64 {
	return (m / F2M)
}
func Meter(m float64) {
	fmt.Fprintf(os.Stdout, "%6.2f M = %6.2f F\n", m, MToF(m))
}
func FToM(f float64) float64 {
	return (f * F2M)
}
func Feet(f float64) {
	fmt.Fprintf(os.Stdout, "%6.2f ft = %6.2f m\n", f, FToM(f))
}

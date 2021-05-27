/*
Exercise 5.15: Write variadic functions max and min, analogous to sum .
What should these functions do when called with no arguments? Write
variants that require at least one argument.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func max(vals ...int64) (int64, error) {
	if len(vals) < 1 {
		return 0, fmt.Errorf("Error: at least one value needs to be supplied")
	}
	mx := vals[0]
	for _, val := range vals {
		if val > mx {
			mx = val
		}
	}
	return mx, nil
}

func min(vals ...int64) (int64, error) {
	if len(vals) < 1 {
		return 0, fmt.Errorf("Error: at least one value needs to be supplied")
	}
	mn := vals[0]
	for _, val := range vals {
		if val < mn {
			mn = val
		}
	}
	return mn, nil
}

func main() {
	actionFs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	maxF := actionFs.Bool("max", false, "prints the maximum value of remaining args")
	minF := actionFs.Bool("min", false, "prints the minimum value of remaining args")
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	actionFs.Parse(os.Args[1:])
	if *maxF == *minF {
		usage()
		os.Exit(1)
	}
	ints, err := argsToInt(actionFs.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: ERROR: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	var f func(...int64) (int64, error)
	switch {
	case *maxF:
		f = max
	case *minF:
		f = min
	}
	if result, err := f(ints...); err != nil {
		fmt.Fprintf(os.Stderr, "%s: ERROR: %v\n", os.Args[0], err)
		os.Exit(1)
	} else {
		fmt.Printf("Result: %v\n", result)
		os.Exit(0)
	}
}

func argsToInt(a []string) ([]int64, error) {
	r := []int64{}
	for _, s := range a {
		if i, err := strconv.ParseInt(s, 10, 64); err != nil {
			return []int64{}, err
		} else {
			r = append(r, i)
		}
	}
	return r, nil
}

func usage() {
	fmt.Fprintf(os.Stderr,
		"\nUsage: %s <-max|-min> <List of numbers>\n",
		os.Args[0])
}

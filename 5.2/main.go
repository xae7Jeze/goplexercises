/*
* Exercise 5.2: Write a function to populate a mapping from element names
* — p , div , span , and so on — to the number of elements with that name
* in an HTML document tree.
 */

package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
)

type ecount = map[string]int

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	e := make(ecount)
	walkChildren(doc, &e)
	for k, v := range e {
		fmt.Printf("%-10s : %d\n", k, v)
	}
}

func walkChildren(n *html.Node, e *ecount) {
	if n == nil {
		return
	}
	if n.Type == html.ElementNode {
		(*e)[n.Data]++
	}
	walkChildren(n.FirstChild, e)
	walkChildren(n.NextSibling, e)
}

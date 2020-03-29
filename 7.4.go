/*
Exercise 7.4: The strings.NewReader function returns a value that satisfies the
io.Reader interface (and others) by reading from its argument, a string .
Implement a simple version of NewReader yourself, and use it to make the HTML
parser (§5.2) take input from a string.

reads html from stdin, saves it to a string and pretty prints it

*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type byteSlice []byte

func NewReader(s string) *byteSlice {
	r := byteSlice(s)
	return &r
}

func (r *byteSlice) Read(b []byte) (n int, err error) {
	copied := copy(b, *r)
	*r = (*r)[copied:]
	if copied <= 0 {
		err = io.EOF
	}
	return copied, err
}

func main() {
	b := make([]byte, 4096)
	r := bufio.NewReader(os.Stdin)
	var s string
	for {
		br, err := r.Read(b)
		if br == 0 && err == io.EOF {
			break
		}
		s += string(b[:br])
	}
	err := prettyPrint(s, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oops: an error occured: %v\n", err)
		os.Exit(1)
	}
}

// PrettyPrinter slighly modified from lesson 5.7

var depth int = -1

func prettyPrint(s string, out io.Writer) error {
	sr := NewReader(s)
	doc, err := html.Parse(sr)
	if err != nil {
		return err
	}

	forEachNode(doc, startElement, endElement, out)
	return nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node, out io.Writer), out io.Writer) {
	if pre != nil {
		pre(n, out)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post, out)
	}

	if post != nil {
		post(n, out)
	}
}

func startElement(n *html.Node, out io.Writer) {
	//printNode(n)
	switch n.Type {
	case html.ElementNode:
		depth++
		fmt.Fprintf(os.Stdout, "%*s<%s", depth*2, "", n.Data)
		printAttributes(&n.Attr, html.ElementNode, out)
		if n.FirstChild != nil {
			fmt.Fprintf(out, ">\n")
		} else {
			fmt.Fprintf(out, "/>\n")
		}
	case html.TextNode:
		//printNode(n)
		s := strings.TrimSpace(n.Data)
		if s != "" {
			if len(s)+2*depth > 80 {
				for _, s := range strings.Split(strings.ReplaceAll(s, "\r", ""), "\n") {
					fmt.Fprintf(out, "%*s%s\n", depth*2+depth, "", strings.TrimSpace(s))
				}
			} else {
				fmt.Fprintf(out, "%*s%s\n", depth*2+depth, "", strings.TrimSpace(s))
			}
		}
	case html.DocumentNode:
	case html.CommentNode:
		fmt.Fprintf(out, "%*s <!--%s-->\n", depth*2, "", n.Data)
	case html.DoctypeNode:
		fmt.Fprintf(out, "<!DOCTYPE %s", n.Data)
		printAttributes(&n.Attr, html.DoctypeNode, out)
		fmt.Fprintf(out, ">\n")
	default:
		fmt.Fprintf(out, "%*sOOPS %s", depth*2, "", n.Data)
	}
}

func endElement(n *html.Node, out io.Writer) {
	switch n.Type {
	case html.ElementNode:
		if n.FirstChild != nil {
			fmt.Fprintf(out, "%*s</%s>\n", depth*2, "", n.Data)
		}
		depth--
	case html.TextNode:
	case html.DocumentNode:
	case html.CommentNode:
	case html.DoctypeNode:
	case html.ErrorNode:
	default:
		fmt.Fprintf(out, "%*sOOPS %s", depth*2, "", n.Data)
	}
}

func printAttributes(as *[]html.Attribute, nt html.NodeType, out io.Writer) {
	l := len(*as)
	if l < 1 {
		return
	}
	fmt.Fprintf(out, " ")
	for i, a := range *as {
		if nt == html.DoctypeNode {
			if a.Key == "public" {
				fmt.Fprintf(out, "PUBLIC \"%s\"", a.Val)
			}
			if a.Key == "system" {
				fmt.Fprintf(out, "\"%s\"", a.Val)
			}
		} else {
			fmt.Fprintf(out, "%s=\"%s\"", a.Key, a.Val)
		}
		if i < (l - 1) {
			fmt.Fprintf(out, " ")
		}
	}
}

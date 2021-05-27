/*
Exercise 5.12: The startElement and endElement functions in
gopl.io/ch5/outline2 (§5.5) share a global variable, depth .
Turn them into anonymous functions that share a variable local
to the outline function.
*/
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	//	"unicode"

	"golang.org/x/net/html"
)

var out io.Writer = os.Stdout

func main() {
	for _, url := range os.Args[1:] {
		outline(url)
	}
}

var fDepth func(i int) int

func outline(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var depth int = -1
	fDepth = func(i int) int {
		depth += i
		return depth
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	forEachNode(doc, startElement, endElement)
	return nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

func startElement(n *html.Node) {
	//printNode(n)
	switch n.Type {
	case html.ElementNode:
		fDepth(1)
		fmt.Fprintf(out, "%*s<%s", fDepth(0)*2, "", n.Data)
		printAttributes(&n.Attr, html.ElementNode)
		if n.FirstChild != nil {
			fmt.Fprintf(out, ">\n")
		} else {
			fmt.Fprintf(out, "/>\n")
		}
	case html.TextNode:
		//printNode(n)
		s := strings.TrimSpace(n.Data)
		if s != "" {
			if len(s)+2*fDepth(0) > 80 {
				for _, s := range strings.Split(strings.ReplaceAll(s, "\r", ""), "\n") {
					fmt.Fprintf(out, "%*s%s\n", fDepth(0)*2+fDepth(0), "", strings.TrimSpace(s))
				}
			} else {
				fmt.Fprintf(out, "%*s%s\n", fDepth(0)*2+fDepth(0), "", strings.TrimSpace(s))
			}
		}
	case html.DocumentNode:
	case html.CommentNode:
		fmt.Fprintf(out, "%*s <!--%s-->\n", fDepth(0)*2, "", n.Data)
	case html.DoctypeNode:
		fmt.Fprintf(out, "<!DOCTYPE %s", n.Data)
		printAttributes(&n.Attr, html.DoctypeNode)
		fmt.Fprintf(out, ">\n")
	default:
		fmt.Fprintf(out, "%*sOOPS %s", fDepth(0)*2, "", n.Data)
	}
}

func endElement(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		if n.FirstChild != nil {
			fmt.Fprintf(out, "%*s</%s>\n", fDepth(0)*2, "", n.Data)
		}
		fDepth(-1)
	case html.TextNode:
	case html.DocumentNode:
	case html.CommentNode:
	case html.DoctypeNode:
	case html.ErrorNode:
	default:
		fmt.Fprintf(out, "%*sOOPS %s", fDepth(0)*2, "", n.Data)
	}
}

func printNode(n *html.Node) {
	var m = []string{"ErrorNode", "TextNode", "DocumentNode", "ElementNode", "CommentNode", "DoctypeNode"}
	fmt.Fprintf(out,
		"\n\n-------- Node: %p: ----------------\n"+
			"  Parent: %p\n"+
			"  FirstChild %p\n"+
			"  LastChild %p\n"+
			"  PrevSibling %p\n"+
			"  NextSibling %p\n"+
			"  Type: %s\n"+
			"  DataAtom: '%v'\n"+
			"  Data: '%s'\n"+
			"  Attributes:\n\n",
		n, n.Parent, n.FirstChild, n.LastChild, n.PrevSibling, n.NextSibling, m[n.Type], n.DataAtom, n.Data)

	for _, a := range n.Attr {
		fmt.Fprintf(out, "    '%s' = '%s'", a.Key, a.Val)
		if a.Namespace != "" {
			fmt.Fprintf(out, " [Namespace: '%s']", a.Namespace)
		}
		fmt.Fprintf(out, "\n")
	}
	fmt.Fprintf(out, "---------------------------------------------------\n\n")

}

func printAttributes(as *[]html.Attribute, nt html.NodeType) {
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

/* ErrorNode NodeType = iota
   TextNode
   DocumentNode
   ElementNode
   CommentNode
   DoctypeNode
*/

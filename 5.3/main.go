/*
Exercise 5.3: Write a function to print the contents of all text nodes in an
HTML document tree. Do not descend into <script> or <style> elements, since
their contents are not visible in a web browser.
*/
package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	walkChildren(doc)
}

func walkChildren(n *html.Node) {
	if n == nil {
		return
	}
	if n.Data == "script" {
		return
	}
	if n.Type == html.TextNode {
		fmt.Printf("%s", n.Data)
	}
	walkChildren(n.FirstChild)
	walkChildren(n.NextSibling)
}

/*
type Node struct {
        Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

        Type      NodeType
        DataAtom  atom.Atom
        Data      string
        Namespace string
        Attr      []Attribute
}
*/

/*
const (
        ErrorNode NodeType = iota
        TextNode
        DocumentNode
        ElementNode
        CommentNode
        DoctypeNode
)
*/

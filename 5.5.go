/*
Exercise 5.5: Implement countWordsAndImages.
(See Exercise 4.9 for word-splitting)
*/
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	for _, uri := range os.Args[1:] {
		if w, img, err := CountWordsAndImages(uri); err != nil {
			fmt.Fprintf(os.Stderr, "%s: Oops. Error processing '%s: %v'\n", os.Args[0], uri, err)
			continue
		} else {
			fmt.Printf("Found %d words and %d images in '%s'\n", w, img, uri)
		}
	}
	os.Exit(0)

}

func CountWordsAndImages(uri string) (words, images int, err error) {
	var r *http.Response
	if r, err = http.Get(uri); err != nil {
		return
	}
	doc, err := html.Parse(r.Body)
	r.Body.Close()
	if err != nil {
		return
	}
	words, images = countWordsAndImages(doc)
	return words, images, err
}

func countWordsAndImages(n *html.Node) (words, images int) {
	if n == nil {
		return
	}
	/* fixes parse failure of golang.org/x/net/html for '<noscript>' parts */
	if n.Data == "noscript" {
		s := n.FirstChild.Data
		r := strings.NewReader(s)
		nn, _ := html.ParseFragment(r, nil)
		for _, nnn := range nn {
			w, i := countWordsAndImages(nnn)
			words += w
			images += i
		}
		return
	}

	if n.Type == html.TextNode {
		in := bufio.NewScanner(bytes.NewBufferString(n.Data))
		in.Split(bufio.ScanWords)
		for in.Scan() {
			words++
		}
	} else if n.Type == html.ElementNode && n.Data == "img" {
		images++
	}
	w, i := countWordsAndImages(n.FirstChild)
	words += w
	images += i
	w, i = countWordsAndImages(n.NextSibling)
	words += w
	images += i
	return
}

var m = []string{"ErrorNode", "TextNode", "DocumentNode", "ElementNode", "CommentNode", "DoctypeNode"}

func printNode(n *html.Node) {
	fmt.Printf(
		"Node: %p:\n"+
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
		fmt.Printf("    '%s' = '%s' [Namespace: '%s']\n", a.Key, a.Val, a.Namespace)
	}

}

/*
type NodeType int32

const (
        ErrorNode NodeType = iota
        TextNode
        DocumentNode
        ElementNode
        CommentNode
        DoctypeNode
)

type Attribute struct {
        Namespace, Key, Val string
}
    An Attribute is an attribute namespace-key-value triple. Namespace is
    non-empty for foreign attributes like xlink, Key is alphabetic (and hence
    does not contain escapable characters like '&', '<' or '>'), and Val is
    unescaped (it looks like "a<b" rather than "a&lt;b").

    Namespace is only used by the parser, not the tokenizer.

type Node struct {
        Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

        Type      NodeType
        DataAtom  atom.Atom
        Data      string
        Namespace string
        Attr      []Attribute
}
*/

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s <http_uri> [ ... <http_uri>]\n"+
		"Open each http_uri in turn and prints statistics of words and\n"+
		"images in document.\n\n",
		os.Args[0])
}

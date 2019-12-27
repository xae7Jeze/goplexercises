/*
Exercise 5.8: Modify forEachNode so that the pre and post functions return a
boolean result indicating whether to continue the traversal. Use it to write
a function ElementByID with the following signature that finds the first HTML
element with the specified id attribute. The function should stop the traversal
as soon as a match is found.

func ElementByID(doc *html.Node, id string) *html.Node
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	//	"unicode"

	"golang.org/x/net/html"
)

var out io.Writer = os.Stdout

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(1)
	}
	resp, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatal("Oops, getting URI '" + os.Args[1] + "' failed")
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if e := ElementById(doc, os.Args[2]); e != nil {
		printNode(e)
	}
}

func ElementById(doc *html.Node, id string) *html.Node {
	if nf := forEachNode(doc, &id, startElement); nf != nil {
		return nf
	} else {
		return nil
	}
}

func forEachNode(n *html.Node, id *string, pre func(n *html.Node, id *string) bool) *html.Node {
	if pre != nil && pre(n, id) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if n := forEachNode(c, id, pre); n != nil {
			return n
		}
	}
	return nil
}

func startElement(n *html.Node, id *string) bool {
	if n.Type == html.ElementNode || n.Type == html.DoctypeNode {
		if searchAttributes(&n.Attr, html.ElementNode, id) {
			return true
		}
	}
	return false
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

func searchAttributes(as *[]html.Attribute, nt html.NodeType, id *string) bool {
	l := len(*as)
	if l < 1 {
		return false
	}
	for _, a := range *as {
		if a.Key == "id" && a.Val == *id {
			return true
		}
	}
	return false
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s <http_uri> <id_to_find>\n"+
		"Find element with <id_to_find> in <http_uri>\n\n",
		os.Args[0])
}

/* ErrorNode NodeType = iota
   TextNode
   DocumentNode
   ElementNode
   CommentNode
   DoctypeNode
*/

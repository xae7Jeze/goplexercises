/*
Exercise 5.17: Write a variadic function ElementsByTagName that, given
an HTML node tree and zero or more names, returns all the elements that
match one of those names. Here are two example calls:

func ElementsByTagName(doc *html.Node, name ...string) []*html.Node
images := ElementsByTagName(doc, "img")
headings := ElementsByTagName(doc, "h1", "h2", "h3", "h4")
*/

package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}
	pluralS := "s"
	if len(os.Args) == 3 {
		pluralS = ""
	}
	resp, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatal("Oops, getting URI '" + os.Args[1] + "' failed")
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	found := ElementByTagName(doc, os.Args[2:]...)
	fmt.Fprintf(os.Stdout, "Found %v nodes with tag%s %v\n", len(found), pluralS, os.Args[2:])
	for _, n := range found {
		printNode(n)
	}
}

func ElementByTagName(doc *html.Node, name ...string) []*html.Node {
	if len(name) < 1 {
		return []*html.Node{}
	}
	names := make(map[string]bool)
	nodes := []*html.Node{}
	for _, n := range name {
		names[n] = true
	}
	nodes = *(forEachNode(doc, names, &nodes))
	return nodes
}

func forEachNode(n *html.Node, names map[string]bool, nodes *[]*html.Node) *[]*html.Node {
	if _, ok := names[strings.ToLower(n.Data)]; ok {
		*nodes = append(*nodes, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, names, nodes)
	}
	return nodes
}

func printNode(n *html.Node) {
	var m = []string{"ErrorNode", "TextNode", "DocumentNode", "ElementNode", "CommentNode", "DoctypeNode"}
	fmt.Fprintf(os.Stdout,
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
		fmt.Fprintf(os.Stdout, "    '%s' = '%s'", a.Key, a.Val)
		if a.Namespace != "" {
			fmt.Fprintf(os.Stdout, " [Namespace: '%s']", a.Namespace)
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
	fmt.Fprintf(os.Stdout, "---------------------------------------------------\n\n")

}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s <http_uri> <tags_to find>...\n"+
		"Find element with <tags_to_find> in <http_uri>\n\n",
		os.Args[0])
}

/* ErrorNode NodeType = iota
   TextNode
   DocumentNode
   ElementNode
   CommentNode
   DoctypeNode
*/

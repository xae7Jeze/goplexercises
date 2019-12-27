/*
Exercise 5.7: Develop startElement and endElement into a general HTML pretty-printer.
Print comment nodes, text nodes, and the attribut es of each element ( <a href='...'> ). Use
short forms like <img/> instead of <img></img> when an element has no children. Write a
test to ensure that the out put can be parsed successfully. (See Chapter 11.)
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"
	//	"unicode"

	"golang.org/x/net/html"
)

func TestOutline(t *testing.T) {
	f, err := os.Open("./topsites.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oops: Opening ./topsites.txt failed: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	var sites = []string{}
	scn := bufio.NewScanner(f)
	scn.Split(bufio.ScanLines)
	for scn.Scan() {
		sites = append(sites, scn.Text())
	}
	for _, site := range sites {
		out = new(bytes.Buffer)
		if err := outline("http://" + site + "/"); err != nil {
			t.Errorf("Fetching http://%s/ failed: %v", site, err)
			continue
		}
		/* if you want to print outline()'s output, use this */
		//txt := []byte(fmt.Sprintf("%s", out.(*bytes.Buffer).String()))
		//r := bytes.NewReader(txt)
		/* othewise this */
		// Type assertion "out is of type *bytes.Buffer"
		r := bytes.NewReader(out.(*bytes.Buffer).Bytes())
		if _, err := html.Parse(r); err != nil {
			t.Errorf("Parsing pretty version of http://%s/ FAIL: %v", site, err)
		} else {
			t.Logf("Parsing pretty version of http://%s/: OK", site)
		}
		//t.Logf("%s", txt)
	}
}

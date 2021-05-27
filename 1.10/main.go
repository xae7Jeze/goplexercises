/*
Exercise 1.10: Find a web site that produces a large amount of data.
Investigate caching by running fetchall twice in succession to see
whether the rep orted time changes much. Do you get the same content
each time? Modify fetchall to print its out put to a file so it can
be examined.
*/
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type fetchInfo struct {
	elapsed int64
	uri     string
	nbytes  int64
	file    string
}

func main() {
	args := []string{"-wBu"}
	fmt.Printf("Fetching '%v' twice:\n", os.Args[1])
	fI, err := fetch(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: An error occured: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	args = append(args, fI.file)
	defer os.Remove(fI.file)
	fmt.Printf("First run: Elapsed: %v ms, Bytes: %v\n", fI.elapsed, fI.nbytes)
	fI, err = fetch(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: An error occured: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	fmt.Printf("Second run: Elapsed: %v ms, Bytes: %v\n", fI.elapsed, fI.nbytes)
	args = append(args, fI.file)
	defer os.Remove(fI.file)
	cmd := exec.Command("diff", args...)
	outS := ""
	if out, err := cmd.Output(); (err != nil) && (cmd.ProcessState.ExitCode() > 1) {
		fmt.Fprintf(os.Stderr, "%s: An error occured: %v\n", os.Args[0], err)
		os.Exit(1)
	} else {
		outS = string(out)
	}
	if len(outS) == 0 {
		outS = " no diff"
	} else {
		outS = "\n" + outS
	}
	fmt.Printf("Diff:%s\n", outS)
}

func fetch(url string) (fetchInfo, error) {
	var fI fetchInfo
	var err error
	fI.uri = url
	c := &http.Client{
		Timeout: 30 * time.Second,
	}
	start := time.Now()
	r, err := c.Get(url)
	if err != nil {
		return fI, err
	}
	var f *os.File
	for _, d := range []string{"", "/tmp"} {
		if f, err = ioutil.TempFile(d, ".fetched"); err == nil {
			break
		}
	}
	if err != nil {
		return fI, err
	}
	nbytes, err := io.Copy(f, r.Body)
	f.Close()
	fI.nbytes = nbytes
	r.Body.Close()
	if err != nil {
		os.Remove(f.Name())
		return fI, err
	}
	fI.elapsed = time.Since(start).Milliseconds()
	fI.file = f.Name()
	return fI, nil
}

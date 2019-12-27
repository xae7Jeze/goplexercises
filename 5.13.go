/*
Exercise 5.13: Modify crawl to make local copies of the pages it finds,
creating directories as necessary. Don’t make copies of pages that come
from a different domain. For example, if the original page comes from
golang.org , save all files from there, but exclude ones from vimeo.com .
*/
package main

import (
	"fmt"
	"gopl.io/ch5/links"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	//	"syscall"
)

//!+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) ([]string, error), worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				if wl, err := f(item); err != nil {
					log.Print(err)
					continue
				} else {
					worklist = append(worklist, wl...)
				}
			}
		}
	}
}

//!-breadthFirst

//!+crawl
var alreadyFetched = make(map[string]bool)

func fetch(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("ERROR: Parsing '%v' failed: %v\n", uri, err.Error())
	}
	if !(u.Scheme == "http" || u.Scheme == "https") {
		return fmt.Errorf("WARN: Unsupported Scheme: %v\n", u.Scheme)
	}
	u.User = nil
	u.Host = strings.ToLower(u.Host)
	host := u.Host
	if u.Path == "" {
		u.Path = "/"
	}
	if u.Path[len(u.Path)-1] == os.PathSeparator {
		u.Path += "index.html"
	}
	rP := strings.NewReplacer(".", "", "-", "")
	if rP.Replace(host) == "" {
		return fmt.Errorf("ERROR: Invalid Hostname: %v\n", host)
	}
	localPath := strings.Replace(u.String(), u.Scheme+"://"+host, ".", 1)
	if strings.HasSuffix(localPath, "/..") {
		localPath = strings.TrimSuffix(localPath, "/..") + url.PathEscape("/..")
	}
	if strings.HasPrefix(localPath, "../") || (strings.Index(localPath, "/../") > -1) {
		localPath = strings.TrimPrefix(localPath, "../") + url.PathEscape("../")
	}
	if strings.Index(localPath, "/../") > -1 {
		localPath = strings.ReplaceAll(localPath, "/../", url.PathEscape("/../"))
	}
	if alreadyFetched[localPath] {
		return nil
	}
	if err := os.MkdirAll(host, 0777); err != nil {
		return fmt.Errorf("ERROR: Creating dir '" + host + "' failed: " + err.Error())
	}
	if err := os.Chdir(host); err != nil {
		return fmt.Errorf("ERROR: Changing dir to '" + host + "' failed: " + err.Error())
	}
	defer os.Chdir("..")
	dir, file := path.Split(localPath)
	if dir == "" || dir == "." || dir == ".." {
		return fmt.Errorf("ERROR: Invalid directory path... Shouldn't happen %v", dir)
	}
	if file == "." || file == ".." || file == "" {
		return fmt.Errorf("ERROR: Invalid file name '%v' ... Shouldn't happen", file)
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("ERROR: Creating dir '" + dir + "' failed: " + err.Error())
	}
	log.Printf("INFO: Fetching '%v' to '%v'\n", uri, localPath)
	f, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return fmt.Errorf("ERROR: Creating file '" + file + "' failed: " + err.Error())
	}
	defer f.Close()
	res, err := http.Get(uri)
	if err != nil {
		return fmt.Errorf("ERROR: Fetching uri '" + uri + "' failed: " + err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("ERROR: Fetching uri '" + uri + "' failed: " + http.StatusText(res.StatusCode))
	}
	var r io.LimitedReader
	r.R = res.Body
	r.N = 50 * 1024 * 1024
	for buf := make([]byte, 64*1024); ; {
		n, err := r.Read(buf)
		if !(err == nil || err == io.EOF) {
			return fmt.Errorf("ERROR: Fetching uri '" + uri + "' failed: " + err.Error())
		}
		if n <= 0 {
			break
		}
		f.Write(buf[:n])
	}
	f.Close()
	if r.N <= 0 {
		return fmt.Errorf("ERROR: Fetching uri '" + uri + "' failed: FileTooBig")
	}
	alreadyFetched[localPath] = true
	return nil
}

func crawl(uri string) ([]string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return []string{}, fmt.Errorf("ERROR: Parsing '%v' failed: %v\n", uri, err.Error())
	}
	if !(u.Scheme == "http" || u.Scheme == "https") {
		return []string{}, fmt.Errorf("WARN: Unsupported Scheme: %v\n", u.Scheme)
	}
	u.Host = strings.ToLower(u.Host)
	host := u.Host
	if err := fetch(uri); err != nil {
		log.Print(err)
	}
	list := []string{}
	l, err := links.Extract(uri)
	if err != nil {
		log.Print(err)
	}
	for _, i := range l {
		i = strings.TrimSpace(i)
		lower := strings.ToLower(i)
		if !(strings.HasPrefix(lower, "https://"+host+"/") || strings.HasPrefix(lower, "http://"+host+"/")) {
			continue
		}
		list = append(list, i)
	}
	return list, nil
}

//!-crawl

//!+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}

//!-main

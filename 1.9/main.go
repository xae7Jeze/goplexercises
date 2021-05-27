/*
Exercise 1.9: Modify fetch to also print the HTTP status code,
found in resp.Status.
*/
package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var c http.Client
	var tr = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	c.Transport = tr
	// don't follow redirects. Show first response
	c.CheckRedirect = cr
	for _, url := range os.Args[1:] {
		if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
			url = "http://" + url
		}
		r, err := c.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error fetching %s: %v\n", os.Args[0], url, err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", r.Status)
		_, err = io.Copy(os.Stdout, r.Body)
		r.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error reading from %s: %v\n", os.Args[0], url, err)
			os.Exit(1)
		}
	}
}
func cr(r *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

/*
Exercise 1.7: The function call io.Copy(dst, src) reads from src and writes
to dst. Use it instead of ioutil.ReadAll to copy the response body to
os.Stdout without requiring a buffer large enoug hto hold the entire stream.
Be sure to check the error result of io.Copy .
*/
package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
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
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
	}
	c.Transport = tr
	for _, url := range os.Args[1:] {
		r, err := c.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error fetching %s: %v\n", os.Args[0], url, err)
			os.Exit(1)
		}
		_, err = io.Copy(os.Stdout, r.Body)
		r.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error reading from %s: %v\n", os.Args[0], url, err)
			os.Exit(1)
		}
	}
}

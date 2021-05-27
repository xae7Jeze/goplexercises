/*
Exercise 4.14: Create a web server that queries GitHub once and then allows navigation of the
list of bug reports, milestones, and users.
*/
package main

import (
	"fmt"
	"github"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	maxItems = 200
)
const templ = `<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>Milestone</th>
  <th>User</th>
  <th>Title</th>
</tr>
{{range .Items}}
<tr style='text-align: left'>
  <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td><a href='{{.Milestone | msGetHTMLURL}}'>{{.Milestone | msGetNumber}}</a></td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
  <td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>	
{{end}}
</table>`

var report = template.Must(template.New("issuesList").
	Funcs(template.FuncMap{"msGetHTMLURL": msGetHTMLURL, "msGetNumber": msGetNumber}).
	Parse(templ))

type param struct {
	q string
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	http.HandleFunc("/", http_handler)
	log.Print("Opening Listening port 8080")
	go http.ListenAndServe("127.0.0.1:8080", nil)
	//log.Fatal(go http.ListenAndServe("127.0.0.1:8080", nil))
	log.Fatal(http.ListenAndServe("[::1]:8080", nil))
}

func http_handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print("oops: parsing url failed")
		return
	}
	if r.Method != "GET" {
		log.Print("answering only to get")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	//pal = rand_palette(ps)
	result, err := github.SearchIssues(append(os.Args[1:], []string{"sort:created"}...), maxItems)
	if err != nil {
		log.Println(err)
		return
	}
	if err := report.Execute(w, result); err != nil {
		log.Println(err)
		return
	}
}

func usage() {
	fmt.Fprintf(os.Stderr,
		"\nUsage: %s <IssueSearchString>\n",
		os.Args[0])
}

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func msGetHTMLURL(m *github.Milestone) string {
	if m != nil {
		return m.HTMLURL
	} else {
		return ""
	}
}
func msGetNumber(m *github.Milestone) string {
	if m != nil {
		return strconv.Itoa(m.Number)
	} else {
		return ""
	}
}

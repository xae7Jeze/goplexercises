package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	BaseURL         = "https://api.github.com/"
	IssuesSearchURL = BaseURL + "search/issues"
	ReposURL        = BaseURL + "repos/"
	HTTPUser        = "_GITHUB_USERNAME_"
	HTTPToken       = "__GITHUB_ACCESS_TOKEN_"
	perPage         = "100"
)

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number        int
	HTMLURL       string `json:"html_url"`
	RepositoryURL string `json:"repository_url"`
	Title         string
	State         string
	User          *User
	Milestone     *Milestone
	CreatedAt     time.Time `json:"created_at"`
	Body          string
}

type NewIssue struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Milestone int      `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	State     string   `json:"state,omitempty"`
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Milestone struct {
	Number      int
	Title       string
	Description string
	HTMLURL     string `json:"html_url"`
}

/*
repos/xae7Jeze/hello-world/issues

{
  "title": "Found a bug",
  "body": "I'm having a problem with this.",
  "assignees": [
    "octocat"
  ],
  "milestone": 1,
  "labels": [
    "bug"
  ]
}

*/

func UpdateOrCreateIssue(repo string, new_issue NewIssue, issue int) (*Issue, error) {
	Issue := new(Issue)
	var create bool = false
	if issue <= 0 {
		create = true
	}
	if i := strings.Index(repo, "/"); i == -1 || i != strings.LastIndex(repo, "/") {
		return nil, fmt.Errorf("Invalid Repo %s. Must be <user>/<repo>", repo)
	}
	if new_issue.Title == "" && new_issue.State != "closed" {
		return nil, fmt.Errorf("Missing Issue Title")
	}
	if new_issue.Milestone < 0 {
		return nil, fmt.Errorf("Invalid Milestone Number")
	}
	json_data, err := json.Marshal(new_issue)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(json_data)
	fmt.Fprintf(os.Stderr, "DEBUG: %s\n", json_data)
	//return nil, fmt.Errorf("Bla")

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	c := &http.Client{Jar: jar}
	var req *http.Request
	if create {
		req, err = http.NewRequest("POST", ReposURL+repo+"/issues", buf)
		fmt.Fprintf(os.Stderr, "DEBUG: %s\n", ReposURL+repo+"/issues")
	} else {
		req, err = http.NewRequest("PATCH", ReposURL+repo+"/issues/"+strconv.Itoa(issue), buf)
		fmt.Fprintf(os.Stderr, "DEBUG: %s\n", ReposURL+repo+"/issues/"+strconv.Itoa(issue))
	}
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(HTTPUser, HTTPToken)
	req.Header.Add("Accept", "application/vnd.github.symmetra-preview+json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if create {
		if res.StatusCode != http.StatusCreated {
			return nil, fmt.Errorf("couldn't create issue: %s", res.Status)
		}
	} else {
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("couldn't update issue: %s", res.Status)
		}
	}
	if err = json.NewDecoder(res.Body).Decode(&Issue); err != nil {
		return nil, err
	}
	return Issue, nil
}

func ReadIssue(repo string, issue int) (*Issue, error) {
	Issue := new(Issue)
	if issue < 0 {
		return nil, fmt.Errorf("Invalid Issue Number")
	}
	if i := strings.Index(repo, "/"); i == -1 || i != strings.LastIndex(repo, "/") {
		return nil, fmt.Errorf("Invalid Repo %s. Must be <user>/<repo>", repo)
	}
	if user_repo := strings.Split(repo, "/"); len(user_repo) == 2 {
		repo = strings.Join(user_repo, "/")
	} else {
		return nil, fmt.Errorf("Invalid Repo %s. Must be <user>/<repo>", repo)
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	c := &http.Client{Jar: jar}
	req, err := http.NewRequest("GET", ReposURL+repo+"/issues/"+strconv.Itoa(issue), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(HTTPUser, HTTPToken)
	req.Header.Add("Accept", "application/vnd.github.symmetra-preview+json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("couldn't read issue: %s", res.Status)
	}
	if err = json.NewDecoder(res.Body).Decode(&Issue); err != nil {
		return nil, err
	}
	return Issue, nil
}

func SearchIssues(terms []string, maxResults int) (*IssuesSearchResult, error) {
	var (
		tResult     IssuesSearchResult
		page              = 1
		resGet            = 0
		rlReset     int64 = 0
		rlRemaining int64 = 0
	)
	tResult.TotalCount = -1
	q := url.QueryEscape(strings.Join(terms, " "))
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	c := &http.Client{Jar: jar}
	for ; true; page++ {
		sResult := new(IssuesSearchResult)
		req, err := http.NewRequest("GET", IssuesSearchURL+"?per_page="+perPage+"&page="+strconv.Itoa(page)+"&q="+q, nil)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "DEBUG: %v\n", IssuesSearchURL+"?per_page="+perPage+"&page="+strconv.Itoa(page)+"&q="+q)
		req.SetBasicAuth(HTTPUser, HTTPToken)
		req.Header.Add("Accept", "application/vnd.github.symmetra-preview+json")
		res, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("search query failed: %s", res.Status)
		}
		if err = json.NewDecoder(res.Body).Decode(&sResult); err != nil {
			return nil, err
		}
		resGet += len(sResult.Items)
		if tResult.TotalCount < 0 {
			tResult.TotalCount = sResult.TotalCount
		}
		tResult.Items = append(tResult.Items, sResult.Items...)
		if resGet >= sResult.TotalCount || resGet >= maxResults {
			break
		}
		if rlReset, err = strconv.ParseInt(res.Header.Get("X-RateLimit-Reset"), 10, 64); err != nil {
			continue
		}
		if rlRemaining, err = strconv.ParseInt(res.Header.Get("X-RateLimit-Remaining"), 10, 64); err != nil {
			continue
		}
		if rlRemaining > 0 {
			continue
		}
		wait := rlReset - time.Now().Unix()
		time.Sleep(time.Duration(wait) * time.Second)
	}
	return &tResult, nil
}

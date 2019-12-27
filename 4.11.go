/*
Exercise 4.11:
Build a tool that lets users create, read, update, and delete GitHub issues from
the command line, invoking their preferred text editor when substantial text input
is required.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type reportbytime struct {
	title string
	qitem string
}

const (
	maxItems  = 100
	defEditor = "vi"
)

func main() {
	//s1, s2, err := edit("", "")
	//fmt.Fprintf(os.Stderr, "OOPS: %s %s %v\n", s1, s2, err)
	//os.Exit(0)
	me := &os.Args[0]
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}
	searchFs := flag.NewFlagSet("search", flag.ExitOnError)
	searchRepo := searchFs.String("r", "", "Repo to search: <user>/<repo>")

	readFs := flag.NewFlagSet("read", flag.ExitOnError)
	readRepo := readFs.String("r", "", "Repo to read: <user>/<repo>")
	readIssueNum := readFs.Int("n", -1, "Issuenumber to read")

	updateFs := flag.NewFlagSet("update", flag.ExitOnError)
	updateRepo := updateFs.String("r", "", "Repo to update: <user>/<repo>")
	updateIssueNum := updateFs.Int("n", 0, "Issuenumber to update")
	updateTitle := updateFs.String("t", "", "Issue-Title (optional)")
	updateBody := updateFs.String("b", "", "Issue-Content (optional)")
	updateAssignees := updateFs.String("a", "", "Assignees (comma separated list, optional)")
	updateLabels := updateFs.String("l", "", "Labels (comma separated list, optional)")
	updateMilestone := updateFs.Int("m", 0, "Milestone (optional)")

	closeFs := flag.NewFlagSet("close", flag.ExitOnError)
	closeRepo := closeFs.String("r", "", "Repo: <user>/<repo>")
	closeIssueNum := closeFs.Int("n", 0, "Issuenumber to close")
	closeTitle := closeFs.String("t", "", "Issue-Title (optional)")
	closeBody := closeFs.String("b", "", "Issue-Content (optional)")
	closeAssignees := closeFs.String("a", "", "Assignees (comma separated list, optional)")
	closeLabels := closeFs.String("l", "", "Labels (comma separated list, optional)")
	closeMilestone := closeFs.Int("m", 0, "Milestone (optional)")
	closeNia := closeFs.Bool("f", false, "Don't call editor")

	createFs := flag.NewFlagSet("create", flag.ExitOnError)
	createRepo := createFs.String("r", "", "Repo to use: <user>/<repo> (needed)")
	createTitle := createFs.String("t", "", "Issue-Title (optional)")
	createBody := createFs.String("b", "", "Issue-Content (optional)")
	createAssignees := createFs.String("a", "", "Assignees (comma separated list, optional)")
	createLabels := createFs.String("l", "", "Labels (comma separated list, optional)")
	createMilestone := createFs.Int("m", 0, "Milestone (optional)")

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "search":
		searchFs.Parse(os.Args[2:])
	case "read":
		readFs.Parse(os.Args[2:])
	case "create":
		createFs.Parse(os.Args[2:])
	case "update":
		updateFs.Parse(os.Args[2:])
	case "close":
		closeFs.Parse(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
	switch {
	case searchFs.Parsed() == true:
		ss := searchFs.Args()
		if *searchRepo != "" {
			ss = append(ss, "repo:"+*searchRepo)
		}
		if err := search(append(ss)); err != nil {
			fmt.Fprintf(os.Stderr, "%s: OOPS: An error occured: %v\n", *me, err)
			os.Exit(2)
		}
	case readFs.Parsed() == true:
		if *readRepo == "" {
			usage()
			os.Exit(1)
		}
		if *readIssueNum < 0 {
			usage()
			os.Exit(1)
		}
		if len(readFs.Args()) > 0 {
			usage()
			os.Exit(1)
		}
		if err := read(*readRepo, *readIssueNum); err != nil {
			fmt.Fprintf(os.Stderr, "%s: OOPS: An error occured: %v\n", *me, err)
			os.Exit(2)
		}
	case createFs.Parsed() == true:
		if *createRepo == "" {
			usage()
			os.Exit(1)
		}
		if *createMilestone < 0 {
			usage()
			os.Exit(1)
		}
		assignees := []string{}
		if *createAssignees != "" {
			assignees = strings.Split(*createAssignees, ",")
		}
		labels := []string{}
		if *createLabels != "" {
			labels = strings.Split(*createLabels, ",")
		}
		if len(createFs.Args()) > 0 {
			usage()
			os.Exit(1)
		}
		if *createTitle == "" || *createBody == "" {
			var err error
			*createTitle, *createBody, err = edit(*createTitle, *createBody)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: OOPS: An error occured: %v\n", *me, err)
			}
		}
		new_issue := github.NewIssue{
			Title:     *createTitle,
			Body:      *createBody,
			Assignees: assignees,
			Labels:    labels,
		}
		if err := create(*createRepo, new_issue); err != nil {
			fmt.Fprintf(os.Stderr, "%s: OOPS: An error occured: %v\n", *me, err)
			os.Exit(2)
		}
	case updateFs.Parsed() == true:
		if *updateRepo == "" {
			usage()
			os.Exit(1)
		}
		if *updateIssueNum <= 0 {
			usage()
			os.Exit(1)
		}
		if *updateMilestone < 0 {
			usage()
			os.Exit(1)
		}
		assignees := []string{}
		if *updateAssignees != "" {
			assignees = strings.Split(*updateAssignees, ",")
		}
		labels := []string{}
		if *updateLabels != "" {
			labels = strings.Split(*updateLabels, ",")
		}
		if *updateTitle == "" || *updateBody == "" {
			var err error
			*updateTitle, *updateBody, err = edit(*updateTitle, *updateBody)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: OOPS: An error occured: %v\n", *me, err)
			}
		}

		if len(updateFs.Args()) > 0 {
			usage()
			os.Exit(1)
		}

		new_issue := github.NewIssue{
			Title:     *updateTitle,
			Body:      *updateBody,
			Assignees: assignees,
			Labels:    labels,
		}
		if err := update(*updateRepo, *updateIssueNum, new_issue); err != nil {
			fmt.Fprintf(os.Stderr, "%s: OOPS: An error occured: %v\n", *me, err)
			os.Exit(2)
		}
	case closeFs.Parsed() == true:
		if *closeRepo == "" {
			usage()
			os.Exit(1)
		}
		if *closeIssueNum <= 0 {
			usage()
			os.Exit(1)
		}
		if *closeMilestone < 0 {
			usage()
			os.Exit(1)
		}
		assignees := []string{}
		if *closeAssignees != "" {
			assignees = strings.Split(*closeAssignees, ",")
		}
		labels := []string{}
		if *closeLabels != "" {
			labels = strings.Split(*closeLabels, ",")
		}
		if *closeTitle == "" || *closeBody == "" {
			var err error
			if *closeNia == false {
				*closeTitle, *closeBody, err = edit(*closeTitle, *closeBody)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: OOPS: An error occured: %v\n", *me, err)
			}
		}
		if len(closeFs.Args()) > 0 {
			usage()
			os.Exit(1)
		}

		new_issue := github.NewIssue{
			Title:     *closeTitle,
			Body:      *closeBody,
			Assignees: assignees,
			Labels:    labels,
			State:     "closed",
		}
		if err := update(*closeRepo, *closeIssueNum, new_issue); err != nil {
			fmt.Fprintf(os.Stderr, "%s: OOPS: An error occured: %v\n", *me, err)
			os.Exit(2)
		}

	default:
		usage()
		os.Exit(1)
	}
}

func edit(Title, Body string) (string, string, error) {
	editor := ""
	var err error
	var f *os.File
	for _, e := range []string{os.Getenv("EDITOR"), defEditor} {
		if editor, err = exec.LookPath(e); err == nil {
			break
		}
	}
	if err != nil {
		return "", "", err
	}
	for _, d := range []string{"", "/tmp"} {
		if f, err = ioutil.TempFile(d, ".ei"); err == nil {
			break
		}
	}
	defer os.Remove(f.Name())
	if err != nil {
		return "", "", err
	}
	Title = Title + "\n"
	if Body != "" {
		Body = Body + "\n"
	}
	comment :=
		"# Format of this file:\n" +
			"# First line is treated as issue's title\n" +
			"# Further lines are taken as its description\n" +
			"# Lines beginning with '#' will be deleted\n"
	f.Write([]byte(Title + Body + comment))
	cmd := exec.Command(editor, f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return "", "", err
	}
	Body = ""
	f.Seek(0, 0)
	{
		var buf string
		var err error
		ior := bufio.NewReader(f)
		buf, err = ior.ReadString('\n')
		if !(err == nil || err == io.EOF) {
			return "", "", err
		}
		Title = buf
		for buf, err = ior.ReadString('\n'); err == nil && err != io.EOF; buf, err = ior.ReadString('\n') {
			if strings.Index(buf, "#") != 0 {
				Body += buf
			}
		}
		if err != io.EOF {
			return "", "", err
		}
	}
	if strings.TrimSpace(Title) == "" {
		Title = ""
	}
	if strings.TrimSpace(Body) == "" {
		Body = ""
	}
	return Title, Body, nil
}

func create(repo string, new_issue github.NewIssue) error {
	if new_issue.Title == "" {
		return fmt.Errorf("MISSING ISSUE TITLE")
	}
	issue, err := github.UpdateOrCreateIssue(repo, new_issue, -1)
	if err != nil {
		return err
	}
	fmt.Printf("\nIssue Created:\n")
	fmt.Printf("\nISSUE: #%-d, REPO: %s, USER: %s, CREATED: %v\n", issue.Number, repo, issue.User.Login, issue.CreatedAt)
	fmt.Printf("\nSubject: %.65s\n", issue.Title)
	fmt.Printf("%s\n", strings.Repeat("-", len("Subject: "+issue.Title)))
	fmt.Printf("\n%s\n", issue.Body)
	return nil
}

func update(repo string, issue_nr int, new_issue github.NewIssue) error {
	if new_issue.Title == "" && new_issue.State != "closed" {
		return fmt.Errorf("MISSING ISSUE TITLE")
	}
	issue, err := github.UpdateOrCreateIssue(repo, new_issue, issue_nr)
	if err != nil {
		return err
	}
	fmt.Printf("\nIssue Updated:\n")
	fmt.Printf("\nISSUE: #%-d, REPO: %s, USER: %s, CREATED: %v STATE: %v\n", issue.Number, repo, issue.User.Login, issue.CreatedAt, issue.State)
	fmt.Printf("\nSubject: %.65s\n", issue.Title)
	fmt.Printf("%s\n", strings.Repeat("-", len("Subject: "+issue.Title)))
	fmt.Printf("\n%s\n", issue.Body)
	return nil
}
func read(repo string, issue int) error {
	item, err := github.ReadIssue(repo, issue)
	if err != nil {
		return err
	}
	fmt.Printf("\nISSUE: #%-d, REPO: %s, USER: %s, CREATED: %v\n", item.Number, repo, item.User.Login, item.CreatedAt)
	fmt.Printf("\nSubject: %.65s\n", item.Title)
	fmt.Printf("%s\n", strings.Repeat("-", len("Subject: "+item.Title)))
	fmt.Printf("\n%s\n", item.Body)
	return nil
}
func search(terms []string) error {
	itemsByRepo := make(map[string][]*github.Issue)
	if result, err := github.SearchIssues(append(terms[:], "sort:created"), maxItems); err != nil {
		return err
	} else {
		header := fmt.Sprintf("%d issues", result.TotalCount)
		if result.TotalCount > maxItems {
			header = fmt.Sprintf("%s (Reporting only the first %v)", header, maxItems)
		}
		header = fmt.Sprintf("%s:\n", header)
		fmt.Printf("%s", header)
		fmt.Printf("%s\n", strings.Repeat("=", len(header)))
		for _, item := range result.Items {
			repo := strings.Replace(item.RepositoryURL, github.ReposURL, "", 1)
			itemsByRepo[repo] = append(itemsByRepo[repo], item)
		}
		for k, v := range itemsByRepo {
			header := fmt.Sprintf("Repo: %s\n", k)
			fmt.Printf("%s", header)
			fmt.Printf("%s\n", strings.Repeat("-", len(header)))
			sort.Sort(ghi(v))
			for _, item := range v {
				fmt.Printf("#%-6[1]d %9.9[2]s %[3]v %.55[4]s\n", item.Number, item.User.Login, item.CreatedAt, item.Title)
			}
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
	return nil
}

/* For sorting Issues by number */
type ghi []*github.Issue

func (issues ghi) Len() int {
	return len(issues)
}
func (issues ghi) Less(i, j int) bool {
	return issues[i].Number < issues[j].Number
}
func (issues ghi) Swap(i, j int) {
	issues[i], issues[j] = issues[j], issues[i]
}

func usage() {
	fmt.Fprintf(os.Stderr,
		"\nUsage: %s {create|close|update|read|search} <ARGS>\n"+
			"\n%[1]s <subcommand> -h will tell you more\n\n",
		os.Args[0])
}

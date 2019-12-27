/*
Exercise 4.10: Modify issues to report the results in age categories, say less than a month old,
less than a year old, and more than a year old.
*/
package main

import (
	"fmt"
	"github"
	"log"
	"os"
	"time"
)

type reportbytime struct {
	title string
	qitem string
}

const (
	maxItems = 100
)

func main() {
	//t := time.Now().Add(31 * 24 * time.Hour)
	//monthAgo := time.Now().AddDate(0, -1, 0).Format("created:>=2006-01-02")
	//Mon Jan 2 15:04:05 MST 2006
	m := make([]reportbytime, 0, 3)
	monthAgo := time.Now().AddDate(0, -1, 0)
	yearAgo := time.Now().AddDate(-1, 0, 0)
	// create array for 3 categories by date of creation
	// - created last month,
	// - newer than one year ago but older than one month
	// - older than one year
	m = append(m, []reportbytime{{
		title: monthAgo.Format("Issues newer 2006-01-02T15:04:05"),
		qitem: monthAgo.Format("created:>=2006-01-02T15:04:05"),
	},
		{
			title: yearAgo.Format("Issues from 2006-01-02T15:04:05 to ") + monthAgo.Format("2006-01-02T15:04:05"),
			qitem: yearAgo.Format("created:2006-01-02T15:04:05..") + monthAgo.Format("2006-01-02T15:04:05"),
		},
		{
			title: yearAgo.Format("Created before 2006-01-02T15:04:05"),
			qitem: yearAgo.Format(" created:<2006-01-02T15:04:05"),
		},
	}...)
	//os.Exit(1)

	for _, v := range m {
		fmt.Printf("%s\n", v.title)
		fmt.Printf("---------------------------------------------------------------------\n")
		if result, err := github.SearchIssues(append(os.Args[1:], []string{v.qitem, "sort:created"}...), maxItems); err != nil {
			log.Fatal(err)
			os.Exit(1)
		} else {
			fmt.Printf("%d issues", result.TotalCount)
			if result.TotalCount > maxItems {
				fmt.Printf(" (Reporting only the first %v)", maxItems)
			}
			fmt.Printf(":\n")
			for _, item := range result.Items {
				fmt.Printf("#%-5[1]d %9.9[2]s %[4]v %.55[3]s\n", item.Number, item.User.Login, item.Title, item.CreatedAt)
			}
		}
		fmt.Printf("\n")
	}
}

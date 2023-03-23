package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	// Setup
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	userPtr := flag.String("user", "Jmainguy", "Github Username")
	orgsPtr := flag.String("orgs", "standouthost", "a list of orgs separated by commas")
	sincePtr := flag.String("since", "2022-01-01", "yyyy-mm-dd date to check for stats since")
	verbosePtr := flag.Bool("verbose", false, "print verbose information about each contribution or not")

	flag.Parse()

	const layout = "2006-01-02"
	orgs := strings.Split(*orgsPtr, ",")
	opened := 0
	merged := 0
	reviewsIDid := 0

	startDate, err := time.Parse(layout, *sincePtr)
	if err != nil {
		panic(err)
	}
	// End Setup

	prs := pullRequestsByUser(client, *userPtr, *sincePtr)
	if *verbosePtr {
		fmt.Println("Pull Requests:")
	}
	for _, pr := range prs {
		var print bool
		if contains(orgs, string(pr.owner)) {
			if pr.merged {
				print = true
				merged++
			}
			if pr.createdAt.After(startDate) {
				print = true
				opened++
			}
		}
		if *verbosePtr {
			if print {
				fmt.Println("  Author:", pr.author)
				fmt.Println("  CreatedAt:", pr.createdAt)
				fmt.Println("  Merged:", pr.merged)
				fmt.Println("  URL:", pr.url)
				fmt.Println("")
			}
		}
	}
	if *verbosePtr {
		fmt.Println("")
		fmt.Println("Pull Request Reviews")
	}

	prrs := pullRequestReviewsByUser(client, *userPtr, *sincePtr)
	for _, prr := range prrs {
		if contains(orgs, string(prr.owner)) {
			if prr.createdAt.After(startDate) {
				reviewsIDid++
				if *verbosePtr {
					fmt.Println("  CreatedAt:", prr.createdAt)
					fmt.Println("  URL:", prr.url)
					fmt.Println("")
				}

			}
		}
	}

	commits := commitsByUser(client, *userPtr, *sincePtr)

	fmt.Printf("Total PR's Opened %d\n", opened)
	fmt.Printf("Total PR's Merged %d\n", merged)
	fmt.Printf("Total Reviews %d\n", reviewsIDid)
	fmt.Printf("Commits: %d\n", commits)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
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
	orgsPtr := flag.String("orgs", "", "a list of orgs separated by commas (optional)")
	sincePtr := flag.String("since", "2022-01-01", "yyyy-mm-dd date to check for stats since")
	verbosePtr := flag.Bool("verbose", false, "print verbose information about each contribution or not")

	flag.Parse()

	const layout = "2006-01-02"
	// If --orgs is empty, don't filter by orgs.
	filterOrgs := *orgsPtr != ""
	var orgs []string
	if filterOrgs {
		orgs = strings.Split(*orgsPtr, ",")
	}
	opened := 0
	merged := 0
	reviewsIDid := 0

	startDate, err := time.Parse(layout, *sincePtr)
	if err != nil {
		panic(err)
	}
	// End Setup

	prs, prRepos := pullRequestsByUser(client, *userPtr, *sincePtr)
	if *verbosePtr {
		fmt.Println("Pull Requests:")
	}
	for _, pr := range prs {
		var print bool
		// If no org filter is set, accept all; otherwise check the org list.
		if !filterOrgs || contains(orgs, string(pr.owner)) {
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

	prrs, prReviewRepos := pullRequestReviewsByUser(client, *userPtr, *sincePtr)
	for _, prr := range prrs {
		// If no org filter is set, accept all; otherwise check the org list.
		if !filterOrgs || contains(orgs, string(prr.owner)) {
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

	commits, commitRepos := commitsByUser(client, *userPtr, *sincePtr)

	fmt.Printf("Total PR's Opened %d\n", opened)
	fmt.Printf("Total PR's Merged %d\n", merged)
	fmt.Printf("Total Reviews %d\n", reviewsIDid)
	fmt.Printf("Commits: %d\n", commits)

	// Build set of repos worked on from PRs, reviews, commits
	orgRepos := make(map[string]map[string]bool)
	addRepo := func(full string) {
		parts := strings.SplitN(full, "/", 2)
		if len(parts) != 2 {
			return
		}
		org := parts[0]
		repo := parts[1]
		if _, ok := orgRepos[org]; !ok {
			orgRepos[org] = make(map[string]bool)
		}
		orgRepos[org][repo] = true
	}
	for _, r := range prRepos {
		addRepo(r)
	}
	for _, r := range prReviewRepos {
		addRepo(r)
	}
	for _, r := range commitRepos {
		addRepo(r)
	}

	// Print sorted by org, then repo
	if len(orgRepos) > 0 {
		fmt.Println("")
		fmt.Println("Repositories worked on:")
		// sort orgs
		orgList := make([]string, 0, len(orgRepos))
		for o := range orgRepos {
			orgList = append(orgList, o)
		}
		sort.Strings(orgList)
		for _, org := range orgList {
			repoMap := orgRepos[org]
			repos := make([]string, 0, len(repoMap))
			for r := range repoMap {
				repos = append(repos, r)
			}
			sort.Strings(repos)
			fmt.Printf("%s:\n", org)
			for _, r := range repos {
				fmt.Printf("  - %s/%s\n", org, r)
			}
		}
	}
}

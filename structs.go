package main

import (
	"github.com/shurcooL/githubv4"
)

type pullRequest struct {
	author    string
	url       githubv4.URI
	merged    bool
	createdAt githubv4.DateTime
	owner     githubv4.String
}

type pullRequestReview struct {
	url       githubv4.URI
	createdAt githubv4.DateTime
	owner     githubv4.String
}

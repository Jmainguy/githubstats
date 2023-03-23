package main

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

func pullRequestsByUser(client *githubv4.Client, user, date string) (pullRequests []pullRequest) {
	var userQuery struct {
		User struct {
			ContributionsCollection struct {
				PullRequestContributions struct {
					Edges []struct {
						Node struct {
							PullRequest struct {
								Author struct {
									Login githubv4.String
								}
								URL        githubv4.URI
								Merged     githubv4.Boolean
								CreatedAt  githubv4.DateTime
								Repository struct {
									Owner struct {
										Login githubv4.String
									}
								}
							}
						}
					}
					PageInfo struct {
						HasNextPage githubv4.Boolean
						EndCursor   githubv4.String
					}
				} `graphql:"pullRequestContributions(first: 100, after: $cursor)"`
			} `graphql:"contributionsCollection(from: $since)"`
		} `graphql:"user(login: $login)"`
	}
	const layout = "2006-01-02"

	since, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}

	variables := map[string]interface{}{
		"login":  githubv4.String(user),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
		"since":  githubv4.DateTime{since},
	}
	for {
		err := client.Query(context.Background(), &userQuery, variables)
		if err != nil {
			panic(err)
		}

		for _, edge := range userQuery.User.ContributionsCollection.PullRequestContributions.Edges {
			var pr pullRequest
			pr.author = string(edge.Node.PullRequest.Author.Login)
			pr.merged = bool(edge.Node.PullRequest.Merged)
			pr.createdAt = edge.Node.PullRequest.CreatedAt
			pr.url = edge.Node.PullRequest.URL
			pr.owner = edge.Node.PullRequest.Repository.Owner.Login
			pullRequests = append(pullRequests, pr)
		}
		if !userQuery.User.ContributionsCollection.PullRequestContributions.PageInfo.HasNextPage {
			break
		} else {
			variables["cursor"] = githubv4.NewString(userQuery.User.ContributionsCollection.PullRequestContributions.PageInfo.EndCursor)
		}
	}
	return pullRequests
}

func pullRequestReviewsByUser(client *githubv4.Client, user, date string) (pullRequestReviews []pullRequestReview) {
	var userQuery struct {
		User struct {
			ContributionsCollection struct {
				PullRequestReviewContributions struct {
					Edges []struct {
						Node struct {
							PullRequestReview struct {
								CreatedAt githubv4.DateTime
								URL       githubv4.URI
							}
							Repository struct {
								Owner struct {
									Login githubv4.String
								}
							}
						}
					}
					PageInfo struct {
						HasNextPage githubv4.Boolean
						EndCursor   githubv4.String
					}
				} `graphql:"pullRequestReviewContributions(first: 100, after: $cursor)"`
			} `graphql:"contributionsCollection(from: $since)"`
		} `graphql:"user(login: $login)"`
	}
	const layout = "2006-01-02"

	since, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}

	variables := map[string]interface{}{
		"login":  githubv4.String(user),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
		"since":  githubv4.DateTime{since},
	}
	for {
		err := client.Query(context.Background(), &userQuery, variables)
		if err != nil {
			panic(err)
		}

		for _, edge := range userQuery.User.ContributionsCollection.PullRequestReviewContributions.Edges {
			var prr pullRequestReview
			prr.createdAt = edge.Node.PullRequestReview.CreatedAt
			prr.url = edge.Node.PullRequestReview.URL
			prr.owner = edge.Node.Repository.Owner.Login
			pullRequestReviews = append(pullRequestReviews, prr)
		}
		if !userQuery.User.ContributionsCollection.PullRequestReviewContributions.PageInfo.HasNextPage {
			break
		} else {
			variables["cursor"] = githubv4.NewString(userQuery.User.ContributionsCollection.PullRequestReviewContributions.PageInfo.EndCursor)
		}
	}
	return pullRequestReviews
}

func commitsByUser(client *githubv4.Client, user, date string) (commits int) {
	var userQuery struct {
		User struct {
			ContributionsCollection struct {
				TotalCommitContributions githubv4.Int
			} `graphql:"contributionsCollection(from: $since)"`
		} `graphql:"user(login: $login)"`
	}
	const layout = "2006-01-02"

	since, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}

	variables := map[string]interface{}{
		"login": githubv4.String(user),
		"since": githubv4.DateTime{since},
	}
	err = client.Query(context.Background(), &userQuery, variables)
	if err != nil {
		panic(err)
	}

	return int(userQuery.User.ContributionsCollection.TotalCommitContributions)
}

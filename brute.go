package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

func listPullRequests(client *githubv4.Client, repo repo) (pullRequests []pullRequest) {
	var pullRequestsQuery struct {
		Repository struct {
			PullRequests struct {
				Edges []struct {
					Node struct {
						Author struct {
							Login githubv4.String
						}
						State     githubv4.String
						Merged    githubv4.Boolean
						Permalink githubv4.URI
					}
				}
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
			} `graphql:"pullRequests(first: 100, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"name":   githubv4.String(repo.name),
		"owner":  githubv4.String(repo.owner),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	for {
		err := client.Query(context.Background(), &pullRequestsQuery, variables)
		if err != nil {
			panic(err)
		}

		for _, e := range pullRequestsQuery.Repository.PullRequests.Edges {
			var pr pullRequest
			pr.author = string(e.Node.Author.Login)
			pr.merged = bool(e.Node.Merged)
			pr.url = e.Node.Permalink
			pullRequests = append(pullRequests, pr)
		}
		if !pullRequestsQuery.Repository.PullRequests.PageInfo.HasNextPage {
			break
		} else {
			variables["cursor"] = githubv4.NewString(pullRequestsQuery.Repository.PullRequests.PageInfo.EndCursor)
		}
	}

	return pullRequests

}

func listRepos(client *githubv4.Client, organization string) (repos []repo) {
	var reposQuery struct {
		Organization struct {
			Name         githubv4.String
			URL          githubv4.String
			Repositories struct {
				Edges []struct {
					Node struct {
						Name githubv4.String
					}
				}
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
			} `graphql:"repositories(first: 100, after: $cursor)"`
		} `graphql:"organization(login: $login)"`
	}

	// RepoQuery specific variables
	variables := map[string]interface{}{
		"login":  githubv4.String(organization),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	for {
		err := client.Query(context.Background(), &reposQuery, variables)
		if err != nil {
			panic(err)
		}

		for _, edge := range reposQuery.Organization.Repositories.Edges {
			//fmt.Println("Repo Name:", edge.Node.Name)
			var repo repo
			repo.name = string(edge.Node.Name)
			repo.owner = organization
			repos = append(repos, repo)
			/*for _, e := range edge.Node.PullRequests.Edges {
				fmt.Println("  Author:", e.Node.Author.Login)
				fmt.Println("  Merged:", e.Node.Merged)
				fmt.Println("  URL:", e.Node.ResourcePath)
			}*/
		}
		if !reposQuery.Organization.Repositories.PageInfo.HasNextPage {
			break
		} else {
			variables["cursor"] = githubv4.NewString(reposQuery.Organization.Repositories.PageInfo.EndCursor)
		}
	}
	return repos
}

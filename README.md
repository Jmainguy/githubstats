# githubstats
[![Go Report Card](https://goreportcard.com/badge/github.com/Jmainguy/githubstats)](https://goreportcard.com/badge/github.com/Jmainguy/githubstats)
[![Release](https://img.shields.io/github/release/Jmainguy/githubstats.svg?style=flat-square)](https://github.com/Jmainguy/githubstats/releases/latest)
[![Coverage Status](https://coveralls.io/repos/github/Jmainguy/githubstats/badge.svg?branch=main)](https://coveralls.io/github/Jmainguy/githubstats?branch=main)

A golang application to return stats from github via their graphql endpoint

## Requirements
This tool uses the GitHub GraphQL API. You should set a personal access token in the GITHUB_TOKEN environment variable before running the tool.

export example:
```bash
export GITHUB_TOKEN=ghp_yourtokenhere
```

Token permissions (scopes)
- Public data only: a token is not strictly required, but providing a token increases rate limits. No special scopes required for public-only queries.
- Private repositories or private contributions: grant the `repo` scope on the token.
- (Optional) If you need to access organization-only resources beyond your public info, you can also add `read:org`.

Minimum recommendation:
- For typical use across public projects: a token with no additional scopes (or none) is fine.
- To include private repo stats: `repo` scope.

## Usage
```bash
Usage of githubstats:
  -orgs string
    	a list of orgs separated by commas (optional; if omitted, results from all orgs are included)
  -since string
    	yyyy-mm-dd date to check for stats since (default "2022-01-01")
  -user string
    	Github Username (default "Jmainguy")
  -verbose
    	print verbose information about each contribution or not
```

## What's changed
- --orgs is now optional. If you do not pass --orgs, the tool will include contributions from all organizations found in the results.
- At the end of the report the tool now prints a deduplicated list of repositories the user has worked on (from PRs, reviews, and commits), grouped and sorted by organization.
- The commits query now collects repositories the user committed to (up to the GraphQL limits).

## Example
```bash
export GITHUB_TOKEN=ghp_yourtokenhere
githubstats --user yourusername --since 2023-01-01
Total PR's Opened 7
Total PR's Merged 5
Total Reviews 3
Commits: 42

Repositories worked on:
acme:
  - acme/api
  - acme/ui
standouthost:
  - standouthost/infra
```

## PreBuilt Binaries
Grab Binaries from [The Releases Page](https://github.com/Jmainguy/githubstats/releases)

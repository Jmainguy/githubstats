# githubstats
[![Go Report Card](https://goreportcard.com/badge/github.com/Jmainguy/githubstats)](https://goreportcard.com/badge/github.com/Jmainguy/githubstats)
[![Release](https://img.shields.io/github/release/Jmainguy/githubstats.svg?style=flat-square)](https://github.com/Jmainguy/githubstats/releases/latest)
[![Coverage Status](https://coveralls.io/repos/github/Jmainguy/githubstats/badge.svg?branch=main)](https://coveralls.io/github/Jmainguy/githubstats?branch=main)

A golang application to return stats from github via their graphql endpoint

## Requirements
This tool requires a github personal access token to work, the token needs to repo permissions.

You need to export the token before running the tool, it has to be named GITHUB_TOKEN

```/bin/bash
export GITHUB_TOKEN=ghp_sdasdasdewsadsdad2asdsadasdas
```

## Usage
```/bin/bash
Usage of githubstats:
  -orgs string
    	a list of orgs separated by commas (default "standouthost")
  -since string
    	yyyy-mm-dd date to check for stats since (default "2022-01-01")
  -user string
    	Github Username (default "Jmainguy")
  -verbose
    	print verbose information about each contribution or not
```

## Example
```/bin/bash
[jmainguy@jmainguy githubstats]$ export GITHUB_TOKEN=ghp_sdasdasdewsadsdad2asdsadasdas
[jmainguy@jmainguy githubstats]$ githubstats 
Total PR's Opened 7
Total PR's Merged 5
Total Reviews 3
```

## PreBuilt Binaries
Grab Binaries from [The Releases Page](https://github.com/Jmainguy/githubstats/releases)

## Install

### Homebrew

```/bin/bash
brew install Jmainguy/tap/githubstats
```

## Build
```/bin/bash
export GO111MODULE=on
go build
```

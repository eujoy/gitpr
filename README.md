![version](https://img.shields.io/badge/version-v0.1.0-brightgreen)
![golang-version](https://img.shields.io/badge/Go-1.14-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![master-actions Actions Status](https://github.com/eujoy/gitpr/workflows/master-actions/badge.svg)](https://github.com/eujoy/gitpr/actions)

# Commands Usage

## Usage of the script in general

```text
[~/gitpr]$ go run cmd/gitpr/main.go -h
NAME:
   GitPullRequests - CLI tool to check status of pull requests in github, get pull requests of users and extract metrics on them.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.2

AUTHOR:
   Angelos Giannis

COMMANDS:
   find, f             Find the pull requests of multiple user repositories.
   pull-requests, p    Retrieves and prints all the pull requests of a user for a repository.
   user-repos, u       Retrieves and prints the repos of an authenticated user.
   widget, w           Display a widget based terminal which will include all the details required.
   commit-list, c      Retrieves and prints the list of commits between two provided tags or commits.
   create-release, cr  Retrieves all the commits between two tags and creates a list of them to be used a release description..
   pr-metrics, m       Retrieves and prints the number of pull requests for a repository that have been created during a specific time period as well as the lead time of those pull requests.
   help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Usage of `find` command

```text
[~/gitpr]$ go run cmd/gitpr/main.go find -h                                                                                                                                                                                                                [master]
NAME:
   main find - Find the pull requests of multiple user repositories.

USAGE:
   main find [command options] [arguments...]

OPTIONS:
   --auth_token value, -t value  Github authorization token. (default: "~")
   --help, -h                    show help (default: false)
```

## Usage of `pull-requests` command

```text
[~/gitpr]$ go run cmd/gitpr/main.go pull-requests -h                                                                                                                                                                                                       [master]
NAME:
   main pull-requests - Retrieves and prints all the pull requests of a user for a repository.

USAGE:
   main pull-requests [command options] [arguments...]

OPTIONS:
   --auth_token value, -t value  Github authorization token. (default: "~")
   --owner value, -o value       Owner of the repository to retrieve pull requests for.
   --repository value, -r value  Repository name to check.
   --base value, -b value        Base branch to check pull requests against. (default: "master")
   --state value, -a value       State of the pull request. (default: "open")
   --page_size value, -s value   Size of each page to load. (default: 10)
   --help, -h                    show help (default: false)
```

## Usage of `user-repos` command

```text
[~/gitpr]$ go run cmd/gitpr/main.go user-repos -h                                                                                                                                                                                                          [master]
NAME:
   main user-repos - Retrieves and prints the repos of an authenticated user.

USAGE:
   main user-repos [command options] [arguments...]

OPTIONS:
   --auth_token value, -t value  Github authorization token. (default: "~")
   --page_size value, -s value   Size of each page to load. (default: 10)
   --help, -h                    show help (default: false)
```

## Usage of `widget` command

```text
[~/gitpr]$ go run cmd/gitpr/main.go widget -h                                                                                                                                                                                                              [master]
NAME:
   main widget - Display a widget based terminal which will include all the details required.

USAGE:
   main widget [command options] [arguments...]

OPTIONS:
   --auth_token value, -t value  Github authorization token. (default: "~")
   --help, -h                    show help (default: false)
```

## Usage of `commit-list` command

```text
[~/gitpr]$ go run cmd/gitpr/main.go commit-list                                                                                                                                                           *[master]
NAME:
   main commit-list - Retrieves and prints the list of commits between two provided tags or commits.

USAGE:
   main commit-list [command options] [arguments...]

OPTIONS:
   --auth_token value, -t value  Github authorization token. (default: "~")
   --owner value, -o value       Owner of the repository to retrieve pull requests for.
   --repository value, -r value  Repository name to check.
   --start_tag value             The starting tag/commit to compare against.
   --end_tag value               The ending/latest tag/commit to compare against. (default: "HEAD")
   --help, -h                    show help (default: false)
```

## Usage of `create-release` command

```text
[~/gitpr]$ go run cmd/gitpr/main.go create-release -h                                                                                                                                                     *[master]
NAME:
   main create-release - Retrieves all the commits between two tags and creates a list of them to be used a release description..

USAGE:
   main create-release [command options] [arguments...]

OPTIONS:
   --auth_token value, -t value     Github authorization token. (default: "~")
   --owner value, -o value          Owner of the repository to retrieve pull requests for.
   --repository value, -r value     Repository name to check.
   --release_name value, -n value   Define the release name to be set. You can use a string pattern to set the place where the new release tag will be set. (default: "Release version : %v")
   --latest_tag value, -l value     The latest tag to compare against.
   --release_tag value, -v value    Release tag to be used (and checked against if exists). (default: "HEAD")
   --check_pattern value, -p value  Define the pattern to check the files modified against.
   --draft_release, -d              Defines if the release will be a draft or published. (default: false) (default: false)
   --force_create, -f               Forces the creation of the release without asking for confirmation. (default: false) (default: false)
   --help, -h                       show help (default: false)
```

## usage of `pr-metrics` command

```text
[~/gitpr]$ go run cmd/gitpr/main.go pr-metrics -h                                                                                                                                                     *[master]
NAME:
   main pr-metrics - Retrieves and prints the number of pull requests for a repository that have been created during a specific time period as well as the lead time of those pull requests.

USAGE:
   main pr-metrics [command options] [arguments...]

OPTIONS:
   --auth_token value, -t value  Github authorization token. (default: "~")
   --owner value, -o value       Owner of the repository to retrieve pull requests for.
   --repository value, -r value  Repository name to check.
   --base value, -b value        Base branch to check pull requests against. (default: "master")
   --state value, -a value       State of the pull request. (default: "open")
   --start_date value, -f value  Start date of the time range to check. [Expected format: 'yyyy-mm-dd']
   --end_date value, -e value    End date of the time range to check. [Expected format: 'yyyy-mm-dd']
   --print_json, --json          Define whether the output needs to be printed in json format. (default: false)
   --help, -h                    show help (default: false)
```

----

# Definition

## Examples

```shell
go run cmd/gitpr/main.go user-repos -t <your token>
go run cmd/gitpr/main.go pull-requests -t <your token> -o eujoy -r gitpr -a open
go run cmd/gitpr/main.go find -t <your token>
go run cmd/gitpr/main.go commit-list -o eujoy -r gitpr -s <from tag/commit> -e <until tag/commit>
go run cmd/gitpr/main.go create-release -o eujoy -r gitpr -l <previous version> -v <new version> -d -f
go run cmd/gitpr/main.go create-release -o eujoy -r erbuilder -l v0.5.0 -v v0.7.5 -d -p "app/service"
go run cmd/gitpr/main.go pr-metrics -o eujoy -r erbuilder --start_date "2021-01-01" --end_date "2021-03-05"
```

```shell
http://localhost:9999/userRepos?authToken=<your token>&pageSize=10&page=1
http://localhost:9999/pullRequests?authToken=<your token>&repoOwner=taxibeat&repository=core-business&prState=open&baseBranch=&pageSize=10&page=1
http://localhost:9999/defaults
```

## Useful Links

### Bitbucket API documentation

- https://developer.atlassian.com/bitbucket/api/2/reference/
- https://developer.atlassian.com/bitbucket/api/2/reference/meta/authentication
- https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories

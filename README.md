![version](https://img.shields.io/badge/version-v0.1.0-brightgreen)
![golang-version](https://img.shields.io/badge/Go-1.14-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![master-actions Actions Status](https://github.com/eujoy/gitpr/workflows/master-actions/badge.svg)](https://github.com/eujoy/gitpr/actions)

# Commands Usage

## Usage of the script in general

```text
[~/gitpr]$ go run cmd/gitpr/main.go -h                                                                                                                                                                                                                     [master]
NAME:
   GitPullRequests - CLI tool to check status of pull requests in github.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   Angelos Giannis

COMMANDS:
   find, f           Find the pull requests of multiple user repositories.
   pull-requests, p  Retrieves and prints all the pull requests of a user for a repository.
   user-repos, u     Retrieves and prints the repos of an authenticated user.
   widget, w         Display a widget based terminal which will include all the details required.
   help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
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
```

----

# Definition

## Examples

```shell
go run cmd/gitpr/main.go user-repos -t <your token>
go run cmd/gitpr/main.go pull-requests -t <your token> -o taxibeat -r rest -a open
go run cmd/gitpr/main.go find -t <your token>
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

![version](https://img.shields.io/badge/version-0.1.0-blue)
![golang-version](https://img.shields.io/badge/version-1.13-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![master-actions Actions Status](https://github.com/Angelos-Giannis/gitpr/workflows/master-actions/badge.svg)](https://github.com/Angelos-Giannis/gitpr/actions)

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

# Definition

## Examples

```shell
go run cmd/gitpr/main.go user-repos -t ba4b5f3439de62d77eab9d6d91779afda06b4d2b
go run cmd/gitpr/main.go pull-requests -t ba4b5f3439de62d77eab9d6d91779afda06b4d2b -o taxibeat -r rest -a open
go run cmd/gitpr/main.go find -t ba4b5f3439de62d77eab9d6d91779afda06b4d2b
```

```shell
http://localhost:9999/userRepos?authToken=ba4b5f3439de62d77eab9d6d91779afda06b4d2b&pageSize=10&page=1
http://localhost:9999/pullRequests?authToken=ba4b5f3439de62d77eab9d6d91779afda06b4d2b&repoOwner=taxibeat&repository=core-business&prState=open&baseBranch=&pageSize=10&page=1
http://localhost:9999/defaults
```

## Useful Links

### Bitbucket API documentation

- https://developer.atlassian.com/bitbucket/api/2/reference/
- https://developer.atlassian.com/bitbucket/api/2/reference/meta/authentication
- https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories

package domain

// Meta describes the metadata of a response.
type Meta struct {
	PageSize int `json:"page_size"`
	LastPage int `json:"last_page"`
}

// UserReposResponse describes the http response for retrieving user repositories.
type UserReposResponse struct {
	Repositories []Repository `json:"repositories"`
	Meta         Meta         `json:"meta"`
}

// RepoPullRequestsResponse describes the http response for retrieving the pull requests of a repository.
type RepoPullRequestsResponse struct {
	PullRequests []PullRequest `json:"pull_requests"`
	Meta         Meta          `json:"meta"`
}

// CompareTagsResponse describes the http response for retrieving the difference between two tags or commits.
type CompareTagsResponse struct {
	Commits []Commit `json:"commits"`
}
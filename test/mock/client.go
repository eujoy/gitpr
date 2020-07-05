package mock

import (
	"github.com/Angelos-Giannis/gitpr/internal/domain"
	"github.com/stretchr/testify/mock"
)

// Client mock.
type Client struct {
	mock.Mock
}

// GetUserRepos mock implementation.
func (c *Client) GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error) {
	args := c.MethodCalled("GetUserRepos", authToken, pageSize, pageNumber)

	return args.Get(0).(domain.UserReposResponse), args.Error(1)
}

// GetPullRequestsOfRepository mock implementation.
func (c *Client) GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error) {
	args := c.MethodCalled("GetPullRequestsOfRepository", authToken, repoOwner, repository, baseBranch, prState, pageSize, pageNumber)

	return args.Get(0).(domain.RepoPullRequestsResponse), args.Error(1)
}

// GetReviewStateOfPullRequest mock implementation.
func (c *Client) GetReviewStateOfPullRequest(authToken, repoOwner, repository string, pullRequestNumber int) ([]domain.PullRequestReview, error) {
	args := c.MethodCalled("GetReviewStateOfPullRequest", authToken, repoOwner, repository, pullRequestNumber)

	return args.Get(0).([]domain.PullRequestReview), args.Error(1)
}

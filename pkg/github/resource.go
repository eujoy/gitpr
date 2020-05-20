package github

import (
	"github.com/Angelos-Giannis/gitpr/internal/domain"
	"github.com/Angelos-Giannis/gitpr/pkg/github/http"
)

// Resource describes the github resource.
type Resource struct {
	githubClient *http.Client
}

// NewResource prepares and returns a github resource.
func NewResource(githubClient *http.Client) *Resource {
	return &Resource{
		githubClient: githubClient,
	}
}

// GetUserRepos retrieves all the user repositories from github.
func (r *Resource) GetUserRepos(authToken string, pageSize int, pageNumber int) ([]domain.Repository, error) {
	userRepos, err := r.githubClient.GetUserRepos(authToken, pageSize, pageNumber)
	return userRepos, err
}

// GetPullRequestsOfRepository retrieves the pull requests for a specified repo.
func (r *Resource) GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) ([]domain.PullRequest, error) {
	pullRequests, err := r.githubClient.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, pageNumber)
	return pullRequests, err
}

// GetReviewStateOfPullRequest retrieves the reviews of a pull request.
func (r *Resource) GetReviewStateOfPullRequest(authToken, repoOwner, repository string, pullRequestNumber int) ([]domain.PullRequestReview, error) {
	pullRequestReviews, err := r.githubClient.GetReviewStateOfPullRequest(authToken, repoOwner, repository, pullRequestNumber)
	return pullRequestReviews, err
}

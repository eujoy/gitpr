package github

import (
	"github.com/eujoy/gitpr/internal/domain"
)

type githubClient interface {
	GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error)
	GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error)
	GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error)
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
	GetReviewStateOfPullRequest(authToken, repoOwner, repository string, pullRequestNumber int) ([]domain.PullRequestReview, error)
	CreateRelease(authToken, repoOwner, repository, tagName string, draftRelease bool, name, body string) error
}

// Resource describes the github resource.
type Resource struct {
	githubClient githubClient
}

// NewResource prepares and returns a github resource.
func NewResource(githubClient githubClient) *Resource {
	return &Resource{
		githubClient: githubClient,
	}
}

// GetCommitDetails to get the details of a commit.
func (r *Resource) GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error) {
	commitDetails, err := r.githubClient.GetCommitDetails(authToken, repoOwner, repository, commitSha)
	return commitDetails, err
}

// GetDiffBetweenTags to get a list of commits.
func (r *Resource) GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error) {
	diffBetweenTags, err := r.githubClient.GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag)
	return diffBetweenTags, err
}

// GetUserRepos retrieves all the user repositories from github.
func (r *Resource) GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error) {
	userRepos, err := r.githubClient.GetUserRepos(authToken, pageSize, pageNumber)
	return userRepos, err
}

// GetPullRequestsOfRepository retrieves the pull requests for a specified repo.
func (r *Resource) GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error) {
	pullRequests, err := r.githubClient.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, pageNumber)
	return pullRequests, err
}

// GetReviewStateOfPullRequest retrieves the reviews of a pull request.
func (r *Resource) GetReviewStateOfPullRequest(authToken, repoOwner, repository string, pullRequestNumber int) ([]domain.PullRequestReview, error) {
	pullRequestReviews, err := r.githubClient.GetReviewStateOfPullRequest(authToken, repoOwner, repository, pullRequestNumber)
	return pullRequestReviews, err
}

// CreateRelease is responsible for creating a release against a desired repository.
func (r *Resource) CreateRelease(authToken, repoOwner, repository, tagName string, draftRelease bool, name, body string) error {
	err := r.githubClient.CreateRelease(authToken, repoOwner, repository, tagName, draftRelease, name, body)
	return err
}

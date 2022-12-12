package github

import (
	"github.com/eujoy/gitpr/internal/domain"
)

type githubClient interface {
	GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error)
	GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error)
	GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error)
	GetPullRequestsCommits(authToken, repoOwner, repository string, pullRequestNumber, pageSize, pageNumber int) ([]domain.Commit, error)
	GetPullRequestsDetails(authToken, repoOwner, repository string, pullRequestNumber int) (domain.PullRequest, error)
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
	GetReviewStateOfPullRequest(authToken, repoOwner, repository string, pullRequestNumber int) ([]domain.PullRequestReview, error)
	CreateRelease(authToken, repoOwner, repository, tagName string, draftRelease bool, name, body string) error
	GetReleaseList(authToken, repoOwner, repository string, pageSize, pageNumber int) ([]domain.Release, error)
	GetWorkflowExecutions(authToken, repoOwner, repository, startDateStr, endDateStr string, pageSize, pageNumber int) ([]domain.Workflow, error)
	GetWorkflowsOfRepository(authToken, repoOwner, repository string) ([]domain.Workflow, error)
	GetWorkflowTiming(authToken, repoOwner, repository string, runID int) (domain.WorkflowTiming, error)
	GetWorkflowUsage(authToken, repoOwner, repository string, workflowID int) (domain.WorkflowTiming, error)
}

// Resource describes the GitHub resource.
type Resource struct {
	githubClient githubClient
}

// NewResource prepares and returns a GitHub resource.
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

// GetUserRepos retrieves all the user repositories from GitHub.
func (r *Resource) GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error) {
	userRepos, err := r.githubClient.GetUserRepos(authToken, pageSize, pageNumber)
	return userRepos, err
}

// GetPullRequestsCommits retrieves the commits of a specific pull request.
func (r *Resource) GetPullRequestsCommits(authToken, repoOwner, repository string, pullRequestNumber, pageSize, pageNumber int) ([]domain.Commit, error) {
	pullRequestCommits, err := r.githubClient.GetPullRequestsCommits(authToken, repoOwner, repository, pullRequestNumber, pageSize, pageNumber)
	return pullRequestCommits, err
}

// GetPullRequestsDetails retrieves the details of a specific pull request.
func (r *Resource) GetPullRequestsDetails(authToken, repoOwner, repository string, pullRequestNumber int) (domain.PullRequest, error) {
	pullRequestDetails, err := r.githubClient.GetPullRequestsDetails(authToken, repoOwner, repository, pullRequestNumber)
	return pullRequestDetails, err
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

// GetReleaseList fetches the releases that have taken place in a repository.
func (r *Resource) GetReleaseList(authToken, repoOwner, repository string, pageSize, pageNumber int) ([]domain.Release, error) {
	releaseList, err := r.githubClient.GetReleaseList(authToken, repoOwner, repository, pageSize, pageNumber)
	return releaseList, err
}

// GetWorkflowExecutions retrieves the executions of the workflows of a repository.
func (r *Resource) GetWorkflowExecutions(authToken, repoOwner, repository, startDateStr, endDateStr string, pageSize, pageNumber int) ([]domain.Workflow, error) {
	workflows, err := r.githubClient.GetWorkflowExecutions(authToken, repoOwner, repository, startDateStr, endDateStr, pageSize, pageNumber)
	return workflows, err
}

// GetWorkflowsOfRepository retrieves and returns all the workflows of a repository.
func (r *Resource) GetWorkflowsOfRepository(authToken, repoOwner, repository string) ([]domain.Workflow, error) {
	workflows, err := r.githubClient.GetWorkflowsOfRepository(authToken, repoOwner, repository)
	return workflows, err
}

// GetWorkflowTiming retrieves the timing details of a workflow.
func (r *Resource) GetWorkflowTiming(authToken, repoOwner, repository string, runID int) (domain.WorkflowTiming, error) {
	workflowTiming, err := r.githubClient.GetWorkflowTiming(authToken, repoOwner, repository, runID)
	return workflowTiming, err
}

// GetWorkflowUsage retrieves the timing details of a workflow.
func (r *Resource) GetWorkflowUsage(authToken, repoOwner, repository string, workflowID int) (domain.WorkflowTiming, error) {
	workflowTiming, err := r.githubClient.GetWorkflowUsage(authToken, repoOwner, repository, workflowID)
	return workflowTiming, err
}

package pullrequests

import (
	"sort"

	"github.com/eujoy/gitpr/internal/domain"
)

type resource interface {
	GetPullRequestsCommits(authToken, repoOwner, repository string, pullRequestNumber, pageSize, pageNumber int) ([]domain.Commit, error)
	GetPullRequestsDetails(authToken, repoOwner, repository string, pullRequestNumber int) (domain.PullRequest, error)
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
	GetReviewStateOfPullRequest(authToken, repoOwner, repository string, pullRequestNumber int) ([]domain.PullRequestReview, error)
}

// Service describes the user repositories service.
type Service struct {
	resource resource
}

// NewService creates and returns a service instance.
func NewService(resource resource) *Service {
	return &Service{
		resource: resource,
	}
}

// GetPullRequestsCommits retrieves the commits of a specific pull request.
func (s *Service) GetPullRequestsCommits(authToken, repoOwner, repository string, pullRequestNumber, pageSize, pageNumber int) ([]domain.Commit, error) {
	pullRequestCommits, err := s.resource.GetPullRequestsCommits(authToken, repoOwner, repository, pullRequestNumber, pageSize, pageNumber)
	if err != nil {
		return []domain.Commit{}, err
	}

	return pullRequestCommits, nil
}

// GetPullRequestsDetails retrieves the details of a specific pull request.
func (s *Service) GetPullRequestsDetails(authToken, repoOwner, repository string, pullRequestNumber int) (domain.PullRequest, error) {
	pullRequestDetails, err := s.resource.GetPullRequestsDetails(authToken, repoOwner, repository, pullRequestNumber)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return pullRequestDetails, nil
}

// GetPullRequestsOfRepository retrieves the pull requests for a specified repo.
// @todo Improve performance of the flow.
func (s *Service) GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error) {
	pullRequests, err := s.resource.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, pageNumber)
	if err != nil {
		return domain.RepoPullRequestsResponse{}, err
	}

	for pr := range pullRequests.PullRequests {
		reviewStates, err := s.resource.GetReviewStateOfPullRequest(authToken, repoOwner, repository, pullRequests.PullRequests[pr].Number)
		if err != nil {
			return domain.RepoPullRequestsResponse{}, err
		}

		pullRequests.PullRequests[pr].ReviewStates = getLatestReviewStatus(pullRequests.PullRequests[pr].Reviewers, reviewStates)
	}

	return pullRequests, nil
}

// getLatestReviewStatus retrieve the latest pull request reviews state.
func getLatestReviewStatus(prReviewers []domain.User, reviews []domain.PullRequestReview) map[string]string {
	latestReviewState := map[string]string{}
	for i := range prReviewers {
		latestReviewState[prReviewers[i].Username] = "PENDING"
	}

	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].SubmittedAt.Before(reviews[j].SubmittedAt)
	})

	for i := range reviews {
		latestReviewState[reviews[i].User.Username] = reviews[i].State
	}

	return latestReviewState
}

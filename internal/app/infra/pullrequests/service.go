package pullrequests

import (
	"fmt"
	"sort"

	"github.com/Angelos-Giannis/gitpr/internal/domain"
)

type resource interface{
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) ([]domain.PullRequest, error)
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

// @todo Improve performance of the flow.
// GetPullRequestsOfRepository retrieves the pull requests for a specified repo.
func (s *Service) GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) ([]domain.PullRequest, error) {
	pullRequests, err := s.resource.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, pageNumber)
	if err != nil {
		fmt.Println(err)
		return []domain.PullRequest{}, err
	}

	for pr := range pullRequests {
		reviewStates, err := s.resource.GetReviewStateOfPullRequest(authToken, repoOwner, repository, pullRequests[pr].Number)
		if err != nil {
			return []domain.PullRequest{}, err
		}

		pullRequests[pr].ReviewStates = getLatestReviewStatus(pullRequests[pr].Reviewers, reviewStates)
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

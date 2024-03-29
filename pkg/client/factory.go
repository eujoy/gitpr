package client

import (
	"errors"
	"net/http"
	"time"

	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/eujoy/gitpr/pkg/client/github"
	githubHttp "github.com/eujoy/gitpr/pkg/client/github/http"
)

// Client describes the functions that muse be implemented by any client of the factory.
type Client interface {
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

// Factory describes the factory for allowing the usage of several external clients of git repos.
type Factory struct {
	client Client
}

// NewFactory creates the client and returns a factory.
func NewFactory(useClient string, cfg config.Config) (*Factory, error) {
	switch useClient {
	case "github":
		gcl := githubHttp.NewClient(&http.Client{Timeout: cfg.Clients.Github.Timeout * time.Second}, cfg)
		gr := github.NewResource(gcl)

		return &Factory{client: gr}, nil
	default:
		return nil, errors.New("failed to initialize client")
	}
}

// GetClient returns the generated factory.
func (f *Factory) GetClient() Client {
	return f.client
}

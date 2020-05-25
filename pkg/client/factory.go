package client

import (
	"errors"
	"net/http"
	"time"

	"github.com/Angelos-Giannis/gitpr/internal/config"
	"github.com/Angelos-Giannis/gitpr/internal/domain"
	"github.com/Angelos-Giannis/gitpr/pkg/client/github"
	githubhttp "github.com/Angelos-Giannis/gitpr/pkg/client/github/http"
)

// Client describes the functions that muse be implemented by any client of the factory.
type Client interface {
	GetUserRepos(authToken string, pageSize int, pageNumber int) ([]domain.Repository, error)
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) ([]domain.PullRequest, error)
	GetReviewStateOfPullRequest(authToken, repoOwner, repository string, pullRequestNumber int) ([]domain.PullRequestReview, error)
}

// Factory describes the factory for allowing the usage of several external clients of git repos.
type Factory struct {
	client Client
}

// NewFactory creates the client and returns a factory.
func NewFactory(useClient string, cfg config.Config) (*Factory, error) {
	switch useClient {
	case "github":
		gcl := githubhttp.NewClient(&http.Client{Timeout: cfg.Clients.Github.Timeout * time.Second}, cfg)
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

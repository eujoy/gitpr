package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Angelos-Giannis/gitpr/internal/config"
	"github.com/Angelos-Giannis/gitpr/internal/domain"
)

// Client describes a github client structure.
type Client struct {
	httpClient    *http.Client
	configuration config.Config
}

// NewClient builds and returns a github client
func NewClient(httpClient *http.Client, configuration config.Config) *Client {
	return &Client{
		httpClient:    httpClient,
		configuration: configuration,
	}
}

// GetUserRepos retrieves all the usekkr repositories from github.
func (c *Client) GetUserRepos(authToken string, pageSize int, pageNumber int) ([]domain.Repository, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Github.APIURL, c.configuration.Github.Endpoints.GetUserRepos)
	URL = strings.Replace(URL, "{pageSize}", strconv.Itoa(pageSize), -1)
	URL = strings.Replace(URL, "{pageNumber}", strconv.Itoa(pageNumber), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []domain.Repository{}, err
	}

	req.Header.Add("Accept", c.configuration.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var userRepos []domain.Repository
	err = c.getResponse(req, &userRepos)

	return userRepos, err
}

// GetPullRequestsOfRepository retrieves the pull requests for a specified repo.
func (c *Client) GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) ([]domain.PullRequest, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Github.APIURL, c.configuration.Github.Endpoints.GetUserPullRequestsForRepo)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{prState}", prState, -1)
	if baseBranch != "" {
		URL = strings.Replace(URL, "{baseBranch}", "base=" + baseBranch, -1)
	}
	URL = strings.Replace(URL, "{pageSize}", strconv.Itoa(pageSize), -1)
	URL = strings.Replace(URL, "{pageNumber}", strconv.Itoa(pageNumber), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []domain.PullRequest{}, err
	}

	req.Header.Add("Accept", c.configuration.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var pullRequests []domain.PullRequest
	err = c.getResponse(req, &pullRequests)

	return pullRequests, err
}

// GetReviewStateOfPullRequest retrieves the reviews of a pull request.
func (c *Client) GetReviewStateOfPullRequest(authToken, repoOwner, repository string, pullRequestNumber int) ([]domain.PullRequestReview, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Github.APIURL, c.configuration.Github.Endpoints.GetReviewStatusOfPullRequest)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{pullRequestNumber}", strconv.Itoa(pullRequestNumber), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []domain.PullRequestReview{}, err
	}

	req.Header.Add("Accept", c.configuration.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var pullRequestReviews []domain.PullRequestReview
	err = c.getResponse(req, &pullRequestReviews)

	return pullRequestReviews, err
}

func (c *Client) getResponse(req *http.Request, data interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// fmt.Println("==========================")
	// fmt.Println(string(body))
	// fmt.Println("==========================")

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	return nil
}

package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
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

// GetUserRepos retrieves all the user repositories from github.
func (c *Client) GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.APIURL, c.configuration.Clients.Github.Endpoints.GetUserRepos)
	URL = strings.Replace(URL, "{pageSize}", strconv.Itoa(pageSize), -1)
	URL = strings.Replace(URL, "{pageNumber}", strconv.Itoa(pageNumber), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return domain.UserReposResponse{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var userReposResponse domain.UserReposResponse
	err = c.getResponse(req, &userReposResponse.Repositories, &userReposResponse.Meta)
	userReposResponse.Meta.PageSize = pageSize

	return userReposResponse, err
}

// GetPullRequestsOfRepository retrieves the pull requests for a specified repo.
func (c *Client) GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.APIURL, c.configuration.Clients.Github.Endpoints.GetUserPullRequestsForRepo)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{prState}", prState, -1)
	if baseBranch != "" {
		URL = strings.Replace(URL, "{baseBranch}", "base="+baseBranch, -1)
	}
	URL = strings.Replace(URL, "{pageSize}", strconv.Itoa(pageSize), -1)
	URL = strings.Replace(URL, "{pageNumber}", strconv.Itoa(pageNumber), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return domain.RepoPullRequestsResponse{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var pullRequestResponse domain.RepoPullRequestsResponse
	err = c.getResponse(req, &pullRequestResponse.PullRequests, &pullRequestResponse.Meta)
	pullRequestResponse.Meta.PageSize = pageSize

	return pullRequestResponse, err
}

// GetReviewStateOfPullRequest retrieves the reviews of a pull request.
func (c *Client) GetReviewStateOfPullRequest(authToken, repoOwner, repository string, pullRequestNumber int) ([]domain.PullRequestReview, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.APIURL, c.configuration.Clients.Github.Endpoints.GetReviewStatusOfPullRequest)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{pullRequestNumber}", strconv.Itoa(pullRequestNumber), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []domain.PullRequestReview{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var pullRequestReviews []domain.PullRequestReview
	err = c.getResponse(req, &pullRequestReviews, nil)

	return pullRequestReviews, err
}

// getResponse makes the actual request and converts the response to the respective required format.
// Also, it parses the meta data in case it is required.
func (c *Client) getResponse(req *http.Request, data interface{}, meta *domain.Meta) error {
	// fmt.Println("==========================")
	// fmt.Println(req)
	// fmt.Println("==========================")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if meta != nil {
		err = parseMetaData(resp, meta)
		if err != nil {
			return err
		}
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

// parseMetaData prepares the metadata for the response.
func parseMetaData(response *http.Response, meta *domain.Meta) error {
	lastPage := 1
	regex := regexp.MustCompile(`&page=([0-9]*)`)
	for _, pageNum := range regex.FindAllString(response.Header.Get("Link"), -1) {
		parts := strings.Split(pageNum, "=")

		convertedInt, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		if convertedInt > lastPage {
			lastPage = convertedInt
		}
	}

	meta.LastPage = lastPage

	return nil
}

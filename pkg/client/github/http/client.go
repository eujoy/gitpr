package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

// GetCommitDetails to get the details of a commit.
func (c *Client) GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.GetCommitDetails)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{commitSha}", commitSha, -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return domain.Commit{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var commitInfo domain.Commit
	err = c.getResponse(req, &commitInfo, nil)

	return commitInfo, err
}

// GetDiffBetweenTags to get a list of commits.
func (c *Client) GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.GetDiffBetweenTags)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{existingTag}", existingTag, -1)
	if latestTag == "" {
		latestTag = "HEAD"
	}
	URL = strings.Replace(URL, "{newTag}", latestTag, -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return domain.CompareTagsResponse{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var compareTagsResponse domain.CompareTagsResponse
	err = c.getResponse(req, &compareTagsResponse, nil)

	return compareTagsResponse, err
}

// GetUserRepos retrieves all the user repositories from github.
func (c *Client) GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.GetUserRepos)
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

// GetPullRequestsCommits retrieves the commits of a specific pull request.
func (c *Client) GetPullRequestsCommits(authToken, repoOwner, repository string, pullRequestNumber, pageSize, pageNumber int) ([]domain.Commit, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.GetPullRequestCommits)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{pullRequestNumber}", strconv.Itoa(pullRequestNumber), -1)
	URL = strings.Replace(URL, "{pageSize}", strconv.Itoa(pageSize), -1)
	URL = strings.Replace(URL, "{pageNumber}", strconv.Itoa(pageNumber), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []domain.Commit{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var pullRequestCommits []domain.Commit
	err = c.getResponse(req, &pullRequestCommits, nil)

	return pullRequestCommits, err
}

// GetPullRequestsDetails retrieves the details of a specific pull request.
func (c *Client) GetPullRequestsDetails(authToken, repoOwner, repository string, pullRequestNumber int) (domain.PullRequest, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.GetPullRequestDetails)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{pullRequestNumber}", strconv.Itoa(pullRequestNumber), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return domain.PullRequest{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var pullRequestDetails domain.PullRequest
	err = c.getResponse(req, &pullRequestDetails, nil)

	return pullRequestDetails, err
}

// GetPullRequestsOfRepository retrieves the pull requests for a specified repo.
func (c *Client) GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.GetUserPullRequestsForRepo)
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
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.GetReviewStatusOfPullRequest)
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

// CreateRelease makes a post request to github api to create a new release with description.
func (c *Client) CreateRelease(authToken, repoOwner, repository, tagName string, draftRelease bool, name, body string) error {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.PostCreateRelease)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)

	data := url.Values{}
	data.Set("tag_name", tagName)
	data.Add("draft_release", strconv.FormatBool(draftRelease))
	data.Add("body", body)
	if name != "" {
		data.Add("name", name)
	}

	values := map[string]interface{}{
		"tag_name": tagName,
		"draft":    draftRelease,
		"body":     body,
		"name":     name,
	}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))
	req.Header.Set("Content-Type", "application/json")

	err = c.getResponse(req, nil, nil)

	return err
}

// GetReleaseList fetches the releases that have taken place in a repository.
func (c *Client) GetReleaseList(authToken, repoOwner, repository string, pageSize, pageNumber int) ([]domain.Release, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.GetReleaseList)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{pageSize}", strconv.Itoa(pageSize), -1)
	URL = strings.Replace(URL, "{pageNumber}", strconv.Itoa(pageNumber), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []domain.Release{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var releaseList []domain.Release
	err = c.getResponse(req, &releaseList, nil)

	return releaseList, err
}

// GetWorkflowExecutions retrieves the executions of the workflows of a repository.
func (c *Client) GetWorkflowExecutions(authToken, repoOwner, repository, startDateStr, endDateStr string, pageSize, pageNumber int) ([]domain.Workflow, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.WorkflowRuns)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{pageSize}", strconv.Itoa(pageSize), -1)
	URL = strings.Replace(URL, "{pageNumber}", strconv.Itoa(pageNumber), -1)
	URL = strings.Replace(URL, "{createdFrom}", startDateStr, -1)
	URL = strings.Replace(URL, "{createdTo}", endDateStr, -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []domain.Workflow{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var workflowResp domain.WorkflowResponse
	err = c.getResponse(req, &workflowResp, nil)
	if err != nil {
		return []domain.Workflow{}, err
	}

	return workflowResp.WorkflowRuns, nil
}

// GetWorkflowsOfRepository retrieves and returns all the workflows of a repository.
func (c *Client) GetWorkflowsOfRepository(authToken, repoOwner, repository string) ([]domain.Workflow, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.WorkflowsOfRepository)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []domain.Workflow{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var workflowResp domain.WorkflowResponse
	err = c.getResponse(req, &workflowResp, nil)
	if err != nil {
		return []domain.Workflow{}, err
	}

	return workflowResp.WorkflowDetails, nil
}

// GetWorkflowTiming retrieves the timing details of a workflow.
func (c *Client) GetWorkflowTiming(authToken, repoOwner, repository string, runID int) (domain.WorkflowTiming, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.WorkflowTiming)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{run_id}", strconv.Itoa(runID), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return domain.WorkflowTiming{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var workflowTimingResp domain.WorkflowTiming
	err = c.getResponse(req, &workflowTimingResp, nil)
	if err != nil {
		return domain.WorkflowTiming{}, err
	}

	return workflowTimingResp, nil
}

// GetWorkflowUsage retrieves the timing details of a workflow.
func (c *Client) GetWorkflowUsage(authToken, repoOwner, repository string, workflowID int) (domain.WorkflowTiming, error) {
	URL := fmt.Sprintf("%s%s", c.configuration.Clients.Github.ApiUrl, c.configuration.Clients.Github.Endpoints.WorkflowUsage)
	URL = strings.Replace(URL, "{repoOwner}", repoOwner, -1)
	URL = strings.Replace(URL, "{repository}", repository, -1)
	URL = strings.Replace(URL, "{workflowID}", strconv.Itoa(workflowID), -1)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return domain.WorkflowTiming{}, err
	}

	req.Header.Add("Accept", c.configuration.Clients.Github.Headers.Accept)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))

	var workflowTimingResp domain.WorkflowTiming
	err = c.getResponse(req, &workflowTimingResp, nil)
	if err != nil {
		return domain.WorkflowTiming{}, err
	}

	return workflowTimingResp, nil
}

// getResponse makes the actual request and converts the response to the respective required format.
// Also, it parses the metadata in case it is required.
func (c *Client) getResponse(req *http.Request, data interface{}, meta *domain.Meta) error {
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

	if data != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
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

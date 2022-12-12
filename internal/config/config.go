package config

import (
    "io/ioutil"
    "os"
    "time"

    "gopkg.in/yaml.v2"
)

type application struct {
    Author  string `yaml:"author"`
    Name    string `yaml:"name"`
    Usage   string `yaml:"usage"`
    Version string `yaml:"version"`
}

type clients struct {
    Github github `yaml:"github"`
}

type endpoints struct {
    GetCommitDetails             string `yaml:"get_commit_details"`
    GetDiffBetweenTags           string `yaml:"get_diff_between_tags"`
    GetPullRequestCommits        string `yaml:"get_pull_request_commits"`
    GetPullRequestDetails        string `yaml:"get_pull_request_details"`
    GetReleaseList               string `yaml:"get_release_list"`
    GetReviewStatusOfPullRequest string `yaml:"get_review_status_of_pull_request"`
    GetUserRepos                 string `yaml:"get_user_repos"`
    GetUserPullRequestsForRepo   string `yaml:"get_user_pull_requests_for_repo"`
    PostCreateRelease            string `yaml:"post_create_release"`
    WorkflowRuns                 string `yaml:"get_workflow_details"`
    WorkflowsOfRepository        string `yaml:"get_workflows_of_repository"`
    WorkflowTiming               string `yaml:"get_workflow_execution_timing"`
    WorkflowUsage                string `yaml:"get_workflow_usage"`
}

type headers struct {
    Accept string `yaml:"accept"`
}

type token struct {
    DefaultEnvVar string `yaml:"default_env_var"`
    DefaultValue  string `yaml:"default_value"`
}

type billing struct {
    Linux   float32 `yaml:"linux"`
    MacOS   float32 `yaml:"macOS"`
    Windows float32 `yaml:"windows"`
}

type github struct {
    ApiUrl    string        `yaml:"api_url"`
    Headers   headers       `yaml:"headers"`
    Endpoints endpoints     `yaml:"endpoints"`
    Timeout   time.Duration `yaml:"timeout"`
    Token     token         `yaml:"token"`
    Billing   billing       `yaml:"billing"`
}

type pagination struct {
    Next     string `yaml:"next"`
    Previous string `yaml:"previous"`
    Exit     string `yaml:"exit"`
}

type service struct {
    Mode string `yaml:"mode"`
    Port string `yaml:"port"`
}

type settings struct {
    AllowedPullRequestStates []string `yaml:"allowed_pull_request_states"`
    AvailableClients         []string `yaml:"available_clients"`
    BaseBranch               string   `yaml:"base_branch"`
    DefaultClient            string   `yaml:"default_client"`
    PageSize                 int      `yaml:"page_size"`
    PullRequestState         string   `yaml:"pull_request_state"`
}

type spinner struct {
    HideCursor bool          `yaml:"hide_cursor"`
    Type       int           `yaml:"type"`
    Time       time.Duration `yaml:"time"`
}

// Config describes the configuration of the service.
type Config struct {
    Application application `yaml:"application"`
    Clients     clients     `yaml:"clients"`
    Pagination  pagination  `yaml:"pagination"`
    Service     service     `yaml:"service"`
    Settings    settings    `yaml:"settings"`
    Spinner     spinner     `yaml:"spinner"`
}

// New creates and returns a configuration object for the service.
func New(configFile string) (Config, error) {
    var config Config

    yamlBytes, err := ioutil.ReadFile(configFile)
    if err != nil {
        return Config{}, err
    }

    err = yaml.Unmarshal(yamlBytes, &config)
    if err != nil {
        return Config{}, err
    }

    if config.Clients.Github.Token.DefaultValue == "" {
        config.Clients.Github.Token.DefaultValue = os.Getenv(config.Clients.Github.Token.DefaultEnvVar)
    }

    return config, err
}

package domain

import "time"

// Repository describes the required details to keep for a repo.
type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
	SSHURL   string `json:"ssh_url"`
	Private  bool   `json:"private"`
	Language string `json:"language"`
	Stars    int    `json:"stargazers_count"`
}

// User describes a user account.
type User struct {
	ID int `json:"id"`
	Username string `json:"login"`
}

// Label describes a label of a pull request.
type Label struct {
	Name string `json:"name"`
}

// PullRequestReview describes the state of the pull request reviews.
type PullRequestReview struct {
	ID          int       `json:"id"`
	State       string    `json:"state"`
	User        User      `json:"user"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// PullRequest describes the details of a pull request.
type PullRequest struct {
	ID           int               `json:"id"`
	HTMLURL      string            `json:"html_url"`
	Number       int               `json:"number"`
	Title        string            `json:"title"`
	Reviewers    []User            `json:"requested_reviewers"`
	Labels       []Label           `json:"labels"`
	State        string            `json:"state"`
	ReviewStates map[string]string `json:"reviews"`
	Mergeable    bool              `json:"mergeable"`
}

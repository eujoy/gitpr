package domain

import "time"

// Repository describes the required details to keep for a repo.
type Repository struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HtmlUrl     string `json:"html_url"`
	SshUrl      string `json:"ssh_url"`
	Private     bool   `json:"private"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
}

// User describes a user account.
type User struct {
	ID       int    `json:"id"`
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
	ID             int               `json:"id"`
	HtmlUrl        string            `json:"html_url"`
	Number         int               `json:"number"`
	Title          string            `json:"title"`
	Creator        User              `json:"user"`
	Reviewers      []User            `json:"requested_reviewers"`
	Labels         []Label           `json:"labels"`
	State          string            `json:"state"`
	ReviewStates   map[string]string `json:"reviews"`
	Mergeable      bool              `json:"mergeable"`
	MergeCommitSha string            `json:"merge_commit_sha"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	ClosedAt       time.Time         `json:"closed_at"`
	MergedAt       time.Time         `json:"merged_at"`

	Comments       int `json:"comments"`
	ReviewComments int `json:"review_comments"`
	Commits        int `json:"commits"`
	Additions      int `json:"additions"`
	Deletions      int `json:"deletions"`
	ChangedFiles   int `json:"changed_files"`
}

// Committer describes the
type Committer struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// CommitDetails describes the deepest level commit details that we need.
type CommitDetails struct {
	Message      string    `json:"Message"`
	CommentCount int       `json:"comment_count"`
	Committer    Committer `json:"committer"`
}

// CommitFile describes the file details modified in a commit.
type CommitFile struct {
	Filename string `json:"filename"`
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	Changes int `json:"changes"`
	Status string `json:"status"`
}

// Commit describes the information of a commit.
type Commit struct {
	Sha     string        `json:"sha"`
	Url     string        `json:"url"`
	Details CommitDetails `json:"commit"`
	Author  User          `json:"author"`
	Files   []CommitFile   `json:"files"`
}

// PullRequestMetricDetails describes the pull request lead time details to be kept.
type PullRequestMetricDetails struct {
	Number      int           `json:"number"`
	Title       string        `json:"title"`
	LeadTime    time.Duration `json:"lead_time"`
	TimeToMerge time.Duration `json:"time_to_merge"`
	CreatedAt   time.Time     `json:"created_at"`

	Comments       int `json:"comments"`
	ReviewComments int `json:"review_comments"`
	Commits        int `json:"commits"`
	Additions      int `json:"additions"`
	Deletions      int `json:"deletions"`
	ChangedFiles   int `json:"changed_files"`
}
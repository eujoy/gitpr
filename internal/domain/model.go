package domain

import (
    "math"
    "time"
)

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
    Filename  string `json:"filename"`
    Additions int    `json:"additions"`
    Deletions int    `json:"deletions"`
    Changes   int    `json:"changes"`
    Status    string `json:"status"`
}

// Commit describes the information of a commit.
type Commit struct {
    Sha     string        `json:"sha"`
    Url     string        `json:"url"`
    Details CommitDetails `json:"commit"`
    Author  User          `json:"author"`
    Files   []CommitFile  `json:"files"`
}

// PullRequestMetricDetails describes the pull request lead time details to be kept.
type PullRequestMetricDetails struct {
    Number         int           `json:"number"`
    Title          string        `json:"title"`
    LeadTime       time.Duration `json:"lead_time"`
    TimeToMerge    time.Duration `json:"time_to_merge"`
    StrLeadTime    string        `json:"str_lead_time"`
    StrTimeToMerge string        `json:"str_time_to_merge"`
    CreatedAt      time.Time     `json:"created_at"`

    Comments       int `json:"comments"`
    ReviewComments int `json:"review_comments"`
    Commits        int `json:"commits"`
    Additions      int `json:"additions"`
    Deletions      int `json:"deletions"`
    ChangedFiles   int `json:"changed_files"`
}

// TotalAggregation describes the total aggregation metrics of data.
type TotalAggregation struct {
    Comments       int `json:"comments"`
    ReviewComments int `json:"review_comments"`
    Commits        int `json:"commits"`
    Additions      int `json:"additions"`
    Deletions      int `json:"deletions"`
    ChangedFiles   int `json:"changed_files"`

    LeadTime       time.Duration `json:"lead_time"`
    TimeToMerge    time.Duration `json:"time_to_merge"`
    StrLeadTime    string        `json:"str_lead_time"`
    StrTimeToMerge string        `json:"str_time_to_merge"`
}

// AverageAggregation describes the average aggregation metrics of data.
type AverageAggregation struct {
    Comments       float64 `json:"comments"`
    ReviewComments float64 `json:"review_comments"`
    Commits        float64 `json:"commits"`
    Additions      float64 `json:"additions"`
    Deletions      float64 `json:"deletions"`
    ChangedFiles   float64 `json:"changed_files"`

    LeadTime       time.Duration `json:"lead_time"`
    TimeToMerge    time.Duration `json:"time_to_merge"`
    StrLeadTime    string        `json:"str_lead_time"`
    StrTimeToMerge string        `json:"str_time_to_merge"`
}

// PullRequestMetrics describes the full details required for the metrics.
type PullRequestMetrics struct {
    PRDetails []PullRequestMetricDetails `json:"pr_metrics"`
    Total     TotalAggregation           `json:"total"`
    Average   AverageAggregation         `jsom:"average"`
}

// PullRequestFlowRatio describes the flow ratio information for the pull requests.
type PullRequestFlowRatio struct {
    Created int    `json:"created"`
    Merged  int    `json:"merged"`
    Ratio   string `json:"ratio"`
}

// Release describes the release details.
type Release struct {
    ID          int       `json:"id"`
    Url         string    `json:"url"`
    HtmlUrl     string    `json:"html_url"`
    TagName     string    `json:"tag_name"`
    Name        string    `json:"name"`
    Body        string    `json:"body"`
    Draft       bool      `json:"draft"`
    PreRelease  bool      `json:"prerelease"`
    CreatedAt   time.Time `json:"created_at"`
    PublishedAt time.Time `json:"published_at"`
}

// ReleaseReport describes the fields to generate reports for releases.
type ReleaseReport struct {
    NumberOfDraftReleases     int `json:"number_of_draft_releases"`
    NumberOfPreReleases       int `json:"number_of_pre_releases"`
    NumberOfReleasesCreated   int `json:"number_of_releases_created"`
    NumberOfReleasesPublished int `json:"number_of_releases_published"`

    CreatedToPublishedRatio float64 `json:"created_to_published_ratio"`
}

// CalculateRatioFields of the release report.
func (rr *ReleaseReport) CalculateRatioFields() {
    if rr.NumberOfReleasesPublished > 0 {
        rr.CreatedToPublishedRatio = math.Round((float64(rr.NumberOfReleasesCreated)/float64(rr.NumberOfReleasesPublished))*100) / 100
    } else {
        rr.CreatedToPublishedRatio = float64(rr.NumberOfReleasesCreated)
    }
}

// WorkflowResponse describes the response to get the workflows of a repository.
type WorkflowResponse struct {
    TotalCount      int        `json:"total_count"`
    WorkflowRuns    []Workflow `json:"workflow_runs"`
    WorkflowDetails []Workflow `json:"workflows"`
}

// Workflow describes the information of a commit.
type Workflow struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    State string `json:"state"`
}

// WorkflowTiming describes the timing details for a workflow execution.
type WorkflowTiming struct {
    Billable      Billable `json:"billable"`
    RunDurationMs int      `json:"run_duration_ms"`
}

// Billable describes the environment specific details for the job executions.
type Billable struct {
    Ubuntu  JobDetails `json:"UBUNTU"`
    MacOs   JobDetails `json:"MACOS"`
    Windows JobDetails `json:"WINDOWS"`
}

// JobDetails describes the job details of an environment
type JobDetails struct {
    TotalMs int64    `json:"total_ms"`
    Jobs    int      `json:"jobs"`
    JobRuns []JobRun `json:"job_runs"`
}

// JobRun describes the details of a job execution.
type JobRun struct {
    JobID      int   `json:"job_id"`
    DurationMs int64 `json:"duration_ms"`
}

// WorkflowBilling describes the billing information of a workflow.
type WorkflowBilling struct {
    Name          string
    WorkflowCosts []WorkflowCosts
}

// WorkflowCosts describes the cost of a workflow.
type WorkflowCosts struct {
    EnvType     string
    ExecMinutes int64
    Cost        float32
}

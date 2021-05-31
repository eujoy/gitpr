package command

import (
	"time"

	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/eujoy/gitpr/internal/infra/command/commitlist"
	"github.com/eujoy/gitpr/internal/infra/command/createrelease"
	"github.com/eujoy/gitpr/internal/infra/command/find"
	"github.com/eujoy/gitpr/internal/infra/command/prmetrics"
	"github.com/eujoy/gitpr/internal/infra/command/publishmetrics"
	"github.com/eujoy/gitpr/internal/infra/command/pullrequests"
	"github.com/eujoy/gitpr/internal/infra/command/releasereport"
	"github.com/eujoy/gitpr/internal/infra/command/userrepos"
	"github.com/eujoy/gitpr/internal/infra/command/widget"
	"github.com/urfave/cli/v2"
)

type userReposService interface {
	GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error)
}

type pullRequestsService interface {
	GetPullRequestsCommits(authToken, repoOwner, repository string, pullRequestNumber, pageSize, pageNumber int) ([]domain.Commit, error)
	GetPullRequestsDetails(authToken, repoOwner, repository string, pullRequestNumber int) (domain.PullRequest, error)
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
}

type repositoryService interface {
	GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error)
	CreateRelease(authToken, repoOwner, repository, tagName string, draftRelease bool, name, body string) error
	GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error)
	GetReleaseList(authToken, repoOwner, repository string, pageSize, pageNumber int) ([]domain.Release, error)
	PrintCommitList(commitList []domain.Commit, useTmpl string) (string, error)
}

type tablePrinter interface {
	PrintRepos(repos []domain.Repository)
	PrintPullRequest(pullRequests []domain.PullRequest)
	PrintPullRequestFlowRatio(flowRatioData map[string]*domain.PullRequestFlowRatio)
	PrintPullRequestMetrics(pullRequests domain.PullRequestMetrics)
	PrintReleaseReport(releaseReport domain.ReleaseReport, captionText string)
}

type utilities interface {
	ClearTerminalScreen()
	GetPageOptions(respLength int, pageSize int, currentPage int) []string
	GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool)
	ConvertDurationToString(dur time.Duration) string
}

// Builder describes the builder of the cli commands.
type Builder struct {
	commands            []*cli.Command
	cfg                 config.Config
	userReposService    userReposService
	pullRequestsService pullRequestsService
	repositoryService   repositoryService
	tablePrinter        tablePrinter
	utils               utilities
}

// NewBuilder creates and returns a new command builder.
func NewBuilder(cfg config.Config, userReposService userReposService, pullRequestsService pullRequestsService, repositoryService repositoryService, tablePrinter tablePrinter, utils utilities) *Builder {
	return &Builder{
		commands:            []*cli.Command{},
		cfg:                 cfg,
		userReposService:    userReposService,
		pullRequestsService: pullRequestsService,
		repositoryService:   repositoryService,
		tablePrinter:        tablePrinter,
		utils:               utils,
	}
}

// GetCommands returns the list of allowed commands.
func (b *Builder) GetCommands() []*cli.Command {
	return b.commands
}

// UserRepos retrieves the repositories that the authenticated used has access to.
func (b *Builder) UserRepos() *Builder {
	userReposCmd := userrepos.NewCmd(b.cfg, b.userReposService, b.tablePrinter, b.utils)
	b.commands = append(b.commands, userReposCmd)

	return b
}

// CreatedPullRequests retrieves the number pull requests in a repo that have been created during a specific time period.
func (b *Builder) CreatedPullRequests() *Builder {
	pullRequestsCmd := prmetrics.NewCmd(b.cfg, b.pullRequestsService, b.repositoryService, b.tablePrinter, b.utils)
	b.commands = append(b.commands, pullRequestsCmd)

	return b
}

// PullRequests retrieves the pull requests that the authenticated user has in a specific repo.
func (b *Builder) PullRequests() *Builder {
	pullRequestsCmd := pullrequests.NewCmd(b.cfg, b.pullRequestsService, b.tablePrinter, b.utils)
	b.commands = append(b.commands, pullRequestsCmd)

	return b
}

// Find retrieves the repositories a user has access to and then allows the user to select multiple repos to retrieve
// the pull requests that are open against the selected repositories.
func (b *Builder) Find() *Builder {
	findCmd := find.NewCmd(b.cfg, b.userReposService, b.pullRequestsService, b.tablePrinter, b.utils)
	b.commands = append(b.commands, findCmd)

	return b
}

// Widget is used to display all the details in widgets in terminal.
func (b *Builder) Widget() *Builder {
	widgetCmd := widget.NewCmd(b.cfg, b.userReposService, b.pullRequestsService)
	b.commands = append(b.commands, widgetCmd)

	return b
}

// CommitList is used to retrieve and print a list of all the commits between 2 tags or commits.
func (b *Builder) CommitList() *Builder {
	commitListCmd := commitlist.NewCmd(b.cfg, b.repositoryService)
	b.commands = append(b.commands, commitListCmd)

	return b
}

// CreateRelease is used to create a new release tag using the provided tag value and also define the description of the new release.
func (b *Builder) CreateRelease() *Builder {
	createReleaseCmd := createrelease.NewCmd(b.cfg, b.repositoryService)
	b.commands = append(b.commands, createReleaseCmd)

	return b
}

// ReleaseReport is used to fetch the releases for a desired period and based on the provided pattern to prepare reports.
func (b *Builder) ReleaseReport() *Builder {
	releaseReportCmd := releasereport.NewCmd(b.cfg, b.repositoryService, b.tablePrinter)
	b.commands = append(b.commands, releaseReportCmd)

	return b
}

// PublishPullRequestMetrics retrieves the metrics for pull requests and publishes them to google spreadsheets.
func (b *Builder) PublishPullRequestMetrics() *Builder {
	publishMetricsCmd := publishmetrics.NewCmd(b.cfg, b.pullRequestsService, b.repositoryService, b.utils)
	b.commands = append(b.commands, publishMetricsCmd)

	return b
}

package command

import (
	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/eujoy/gitpr/internal/infra/command/find"
	"github.com/eujoy/gitpr/internal/infra/command/pullrequests"
	"github.com/eujoy/gitpr/internal/infra/command/userrepos"
	"github.com/eujoy/gitpr/internal/infra/command/widget"
	"github.com/urfave/cli"
)

type userReposService interface {
	GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error)
}

type pullRequestsService interface {
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
}

type tablePrinter interface {
	PrintRepos(repos []domain.Repository)
	PrintPullRequest(pullRequests []domain.PullRequest)
}

type utilities interface {
	ClearTerminalScreen()
	GetPageOptions(respLength int, pageSize int, currentPage int) []string
	GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool)
}

// Builder describes the builder of the cli commands.
type Builder struct {
	commands            []cli.Command
	cfg                 config.Config
	userReposService    userReposService
	pullRequestsService pullRequestsService
	tablePrinter        tablePrinter
	utils               utilities
}

// NewBuilder creates and returns a new command builder.
func NewBuilder(cfg config.Config, userReposService userReposService, pullRequestsService pullRequestsService, tablePrinter tablePrinter, utils utilities) *Builder {
	return &Builder{
		commands:            []cli.Command{},
		cfg:                 cfg,
		userReposService:    userReposService,
		pullRequestsService: pullRequestsService,
		tablePrinter:        tablePrinter,
		utils:               utils,
	}
}

// GetCommands returns the list of allowed commands.
func (b *Builder) GetCommands() []cli.Command {
	return b.commands
}

// UserRepos retrieves the repositories that the authenticated used has access to.
func (b *Builder) UserRepos() *Builder {
	userReposCmd := userrepos.NewCmd(b.cfg, b.userReposService, b.tablePrinter, b.utils)
	b.commands = append(b.commands, userReposCmd)

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

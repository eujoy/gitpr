package command

import (
	"github.com/Angelos-Giannis/gitpr/internal/app/infra/command/pullrequests"
	"github.com/Angelos-Giannis/gitpr/internal/app/infra/command/userrepos"
	"github.com/Angelos-Giannis/gitpr/internal/domain"
	"github.com/urfave/cli"
)

type userReposService interface {
	GetUserRepos(authToken string, pageSize int, pageNumber int) ([]domain.Repository, error)
}

type pullRequestsService interface {
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) ([]domain.PullRequest, error)
}

type tablePrinter interface{
	PrintRepos(repos []domain.Repository)
	PrintPullRequest(pullRequests []domain.PullRequest)
}

type utilities interface{
	ClearTerminalScreen()
	GetPageOptions(respLength int, pageSize int, currentPage int) []string
	GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool)
}

// Builder describes the builder of the cli commands.
type Builder struct {
	commands            []cli.Command
	userReposService    userReposService
	pullRequestsService pullRequestsService
	tablePrinter        tablePrinter
	utils               utilities
}

// NewBuilder creates and returns a new command builder.
func NewBuilder(userReposService userReposService, pullRequestsService pullRequestsService, tablePrinter tablePrinter, utils utilities) *Builder {
	return &Builder{
		commands:            []cli.Command{},
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
	userReposCmd := userrepos.NewCmd(b.userReposService, b.tablePrinter, b.utils)
	b.commands = append(b.commands, userReposCmd)

	return b
}

// PullRequests retrieves the pull requests that the authenticated user has in a specific repo.
func (b *Builder) PullRequests() *Builder {
	pullRequestsCmd := pullrequests.NewCmd(b.pullRequestsService, b.tablePrinter, b.utils)
	b.commands = append(b.commands, pullRequestsCmd)

	return b
}

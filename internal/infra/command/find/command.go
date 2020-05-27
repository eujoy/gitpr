package find

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Angelos-Giannis/gitpr/internal/config"
	"github.com/Angelos-Giannis/gitpr/internal/domain"
	"github.com/briandowns/spinner"
	"github.com/urfave/cli"
)

type userReposService interface{
	GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error)
}

type pullRequestsService interface{
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
}

type tablePrinter interface{
	PrintPullRequest(pullRequests []domain.PullRequest)
}

type utilities interface{
	ClearTerminalScreen()
}

func NewCmd(cfg config.Config, userReposService userReposService, pullRequestsService pullRequestsService, tablePrinter tablePrinter, utilities utilities) cli.Command {
	var authToken string
	// var pageSize  int

	findCmd := cli.Command{
		Name:    "find",
		Aliases: []string{"f"},
		Usage:   "Find the pull requests of multiple user repositories.",
		Flags:   []cli.Flag{
			cli.StringFlag{
				Name:        "auth_token, t",
				Usage:       "Github authorization token.",
				Value:       "",
				Destination: &authToken,
				Required:    true,
			},
			//
			// @todo Add a flag to allow providing a comma separated list of pull request creator(s)  - by survey.Input or by cli.Flag.
			// @todo Add a flag to allow providing a comma separated list of pull request reviewer(s) - by survey.Input or by cli.Flag.
			//
		},
		Action:  func(c *cli.Context) {
			spinLoader := spinner.New(spinner.CharSets[cfg.Spinner.Type], cfg.Spinner.Time * time.Millisecond, spinner.WithHiddenCursor(cfg.Spinner.HideCursor))

			utilities.ClearTerminalScreen()
			spinLoader.Start()

			userRepositories := getAllUserRepoNames(userReposService, authToken, cfg.Settings.PageSize)

			spinLoader.Stop()
			utilities.ClearTerminalScreen()

			userReposPrompt := &survey.MultiSelect{
				Message:       "Select the repos to retrieve the pull request of:",
				Options:       userRepositories,
				PageSize:      20,
			}

			var selectedRepos []string
			err := survey.AskOne(userReposPrompt, &selectedRepos)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			baseBranch := promptBranchInput()

			prState := promptPrStateOptions(cfg.Settings.AllowedPullRequestStates, cfg.Settings.PullRequestState)

			utilities.ClearTerminalScreen()

			pullRequests := getPullRequestsOfRepos(pullRequestsService, selectedRepos, authToken, baseBranch, prState, cfg.Settings.PageSize, spinLoader)

			utilities.ClearTerminalScreen()

			tablePrinter.PrintPullRequest(pullRequests)
		},
	}

	return findCmd
}

// promptBranchInput asks for branch and returns the provided option.
func promptBranchInput() string {
	var baseBranch string
	baseBranchPrompt := &survey.Input{
		Message: "Define branch to retrieve pull requests created against of (leave empty for any):",
	}

	err := survey.AskOne(baseBranchPrompt, &baseBranch)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return baseBranch
}

// promptPrStateOptions asks user to select the state of the pull requests from a available list of states.
func promptPrStateOptions(availableOptions []string, defaultPrState string) string {
	var prState string
	prStatePrompt := &survey.Select{
		Message: fmt.Sprintf("Select a pull request state (default '%v'):", defaultPrState),
		Options: availableOptions,
		Default: defaultPrState,
	}

	err := survey.AskOne(prStatePrompt, &prState)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return prState
}

// getAllUserRepoNames retrieve all the repos that the user has access to and returns a list of their names.
func getAllUserRepoNames(userReposService userReposService, authToken string, pageSize int) []string {
	var userRepositories []string
	currentPage := 1

	for {
		userRepos, err := userReposService.GetUserRepos(authToken, pageSize, currentPage)
		if err != nil {
			fmt.Println(err)
			return []string{}
		}

		for _, r := range userRepos.Repositories {
			userRepositories = append(userRepositories, r.FullName)
		}

		if len(userRepos.Repositories) < pageSize {
			break
		}

		currentPage++
	}

	return userRepositories
}

// getPullRequestsOfRepos retrieves all the pull requests of the provided repos.
func getPullRequestsOfRepos(pullRequestsService pullRequestsService, userRepos []string, authToken, baseBranch, prState string, pageSize int, spinLoader *spinner.Spinner) []domain.PullRequest {
	var pullRequests []domain.PullRequest
	for _, r := range userRepos {
		details := strings.Split(r, "/")
		currentPage := 1

		fmt.Printf("Retrieving pull requests of : %v/%v...\n", details[0], details[1])

		spinLoader.Start()

		for {
			prs, err := pullRequestsService.GetPullRequestsOfRepository(authToken, details[0], details[1], baseBranch, prState, pageSize, currentPage)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			pullRequests = append(pullRequests, prs.PullRequests...)

			if len(prs.PullRequests) < pageSize {
				break
			}

			currentPage++
		}

		spinLoader.Stop()
	}

	return pullRequests
}
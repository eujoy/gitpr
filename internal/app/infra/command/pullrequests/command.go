package pullrequests

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Angelos-Giannis/gitpr/internal/domain"
	"github.com/briandowns/spinner"
	"github.com/urfave/cli"
)

var (
	defaultPageSize = 10
	defaultPrState = "open"

	availablePrStates = []string{"all", "open", "closed"}
)

type service interface{
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) ([]domain.PullRequest, error)
}

type tablePrinter interface{
	PrintPullRequest(pullRequests []domain.PullRequest)
}

type utilities interface{
	ClearTerminalScreen()
	GetPageOptions(respLength int, pageSize int, currentPage int) []string
	GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool)
}

// NewCmd creates a new command to retrieve pull requests for a repo.
func NewCmd(service service, tablePrinter tablePrinter, utilities utilities) cli.Command {
	var authToken, repoOwner, repository, baseBranch, prState string
	var pageSize  int

	pullRequestsCmd := cli.Command{
		Name: "pull-requests",
		Aliases: []string{"p"},
		Usage: "Retrieves and prints all the pull requests of a user for a repository.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "auth_token, t",
				Usage:       "Github authorization token.",
				Value:       "",
				Destination: &authToken,
				Required:    true,
			},
			cli.StringFlag{
				Name:        "owner, o",
				Usage:       "Owner of the repository to retrieve pull requests for.",
				Value:       "",
				Destination: &repoOwner,
				Required:    false,
			},
			cli.StringFlag{
				Name:        "repository, r",
				Usage:       "Repository name to check.",
				Value:       "",
				Destination: &repository,
				Required:    false,
			},
			cli.StringFlag{
				Name:        "base, b",
				Usage:       "Base branch to check pull requests against.",
				Value:       "",
				Destination: &baseBranch,
				Required:    false,
			},
			cli.StringFlag{
				Name:        "state, a",
				Usage:       "State of the pull request.",
				Value:       defaultPrState,
				Destination: &prState,
				Required:    false,
			},
			cli.IntFlag{
				Name:        "page_size, s",
				Usage:       "Size of each page to load.",
				Value:       defaultPageSize,
				Destination: &pageSize,
				Required:    false,
			},
		},
		Action: func(c *cli.Context) {
			spinLoader := spinner.New(spinner.CharSets[4], 200*time.Millisecond, spinner.WithHiddenCursor(true))

			currentPage := 1
			shallContinue := true

			prState = validatePrStateAndGetDefault(prState)

			for {
				utilities.ClearTerminalScreen()
				spinLoader.Start()

				pullRequests, err := service.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, currentPage)

				utilities.ClearTerminalScreen()
				spinLoader.Stop()

				tablePrinter.PrintPullRequest(pullRequests)

				var whatToDo string
				prompt := &survey.Select{
					Message: "Choose an option:",
					Options: utilities.GetPageOptions(len(pullRequests), pageSize, currentPage),
				}

				err = survey.AskOne(prompt, &whatToDo)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				currentPage, shallContinue = utilities.GetNextPageNumberOrExit(whatToDo, currentPage)
				if shallContinue {
					continue
				} else {
					fmt.Println("Finished!!")
					return
				}
			}
		},
	}

	return pullRequestsCmd
}

// validatePrStateAndGetDefault checks if the requested state of pull requests is valid and returns
// it in case it is, otherwise it returns the default pull request state.
func validatePrStateAndGetDefault(prState string) string {
	for _, prs := range availablePrStates {
		if prState == prs {
			return prState
		}
	}

	return defaultPrState
}

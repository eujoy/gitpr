package pullrequests

import (
    "fmt"
    "time"

    "github.com/AlecAivazis/survey/v2"
    "github.com/briandowns/spinner"
    "github.com/eujoy/gitpr/internal/config"
    "github.com/eujoy/gitpr/internal/domain"
    "github.com/eujoy/gitpr/internal/infra/flag"
    "github.com/urfave/cli/v2"
)

type service interface {
    GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
}

type tablePrinter interface {
    PrintPullRequest(pullRequests []domain.PullRequest)
}

type utilities interface {
    ClearTerminalScreen()
    GetPageOptions(respLength int, pageSize int, currentPage int) []string
    GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool)
}

// NewCmd creates a new command to retrieve pull requests for a repo.
func NewCmd(cfg config.Config, service service, tablePrinter tablePrinter, utilities utilities) *cli.Command {
    var authToken, repoOwner, repository, baseBranch, prState string
    var pageSize int

    flagBuilder := flag.New(cfg)

    pullRequestsCmd := cli.Command{
        Name:    "pull-requests",
        Aliases: []string{"p"},
        Usage:   "Retrieves and prints all the pull requests of a user for a repository.",
        Flags: flagBuilder.
            AppendAuthFlag(&authToken).
            AppendOwnerFlag(&repoOwner).
            AppendRepositoryFlag(&repository).
            AppendBaseFlag(&baseBranch).
            AppendStateFlag(&prState).
            AppendPageSizeFlag(&pageSize, cfg.Settings.PageSize).
            GetFlags(),
        Action: func(c *cli.Context) error {
            var shallContinue bool
            spinLoader := spinner.New(spinner.CharSets[cfg.Spinner.Type], cfg.Spinner.Time*time.Millisecond, spinner.WithHiddenCursor(cfg.Spinner.HideCursor))

            currentPage := 1

            prState = validatePrStateAndGetDefault(cfg, prState)

            for {
                utilities.ClearTerminalScreen()
                spinLoader.Start()

                prResp, err := service.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, currentPage)
                if err != nil {
                    fmt.Println(err)
                    return err
                }

                spinLoader.Stop()
                utilities.ClearTerminalScreen()

                tablePrinter.PrintPullRequest(prResp.PullRequests)

                var whatToDo string
                prompt := &survey.Select{
                    Message: "Choose an option:",
                    Options: utilities.GetPageOptions(len(prResp.PullRequests), pageSize, currentPage),
                }

                err = survey.AskOne(prompt, &whatToDo)
                if err != nil {
                    fmt.Println(err)
                    return err
                }

                currentPage, shallContinue = utilities.GetNextPageNumberOrExit(whatToDo, currentPage)
                if shallContinue {
                    continue
                } else {
                    fmt.Println("Finished!!")
                    return nil
                }
            }
        },
    }

    return &pullRequestsCmd
}

// validatePrStateAndGetDefault checks if the requested state of pull requests is valid and returns
// it in case it is, otherwise it returns the default pull request state.
func validatePrStateAndGetDefault(cfg config.Config, prState string) string {
    for _, prs := range cfg.Settings.AllowedPullRequestStates {
        if prState == prs {
            return prState
        }
    }

    return cfg.Settings.PullRequestState
}

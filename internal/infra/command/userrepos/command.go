package userrepos

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/eujoy/gitpr/internal/infra/flag"
	"github.com/urfave/cli/v2"
)

type service interface {
	GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error)
}

type tablePrinter interface {
	PrintRepos(repos []domain.Repository)
}

type utilities interface {
	ClearTerminalScreen()
	GetPageOptions(respLength int, pageSize int, currentPage int) []string
	GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool)
}

// NewCmd creates a new command to retrieve the repos of a user.
func NewCmd(cfg config.Config, service service, tablePrinter tablePrinter, utilities utilities) *cli.Command {
	var authToken string
	var pageSize int

	flagBuilder := flag.New(cfg)

	userReposCmd := cli.Command{
		Name:    "user-repos",
		Aliases: []string{"u"},
		Usage:   "Retrieves and prints the repos of an authenticated user.",
		Flags: flagBuilder.
			AppendAuthFlag(&authToken).
			AppendPageSizeFlag(&pageSize).
			GetFlags(),
		Action: func(c *cli.Context) error {
			var shallContinue bool
			spinLoader := spinner.New(spinner.CharSets[cfg.Spinner.Type], cfg.Spinner.Time*time.Millisecond, spinner.WithHiddenCursor(cfg.Spinner.HideCursor))

			currentPage := 1

			for {
				utilities.ClearTerminalScreen()
				spinLoader.Start()

				userRepos, err := service.GetUserRepos(authToken, pageSize, currentPage)
				if err != nil {
					fmt.Println(err)
					return err
				}

				spinLoader.Stop()
				utilities.ClearTerminalScreen()

				tablePrinter.PrintRepos(userRepos.Repositories)

				var whatToDo string
				prompt := &survey.Select{
					Message: "Choose an option:",
					Options: utilities.GetPageOptions(len(userRepos.Repositories), pageSize, currentPage),
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
					break
				}
			}

			return nil
		},
	}

	return &userReposCmd
}

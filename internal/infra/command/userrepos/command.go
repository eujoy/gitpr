package userrepos

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Angelos-Giannis/gitpr/internal/config"
	"github.com/Angelos-Giannis/gitpr/internal/domain"
	"github.com/briandowns/spinner"
	"github.com/urfave/cli"
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
func NewCmd(cfg config.Config, service service, tablePrinter tablePrinter, utilities utilities) cli.Command {
	var authToken string
	var pageSize int

	userReposCmd := cli.Command{
		Name:    "user-repos",
		Aliases: []string{"u"},
		Usage:   "Retrieves and prints the repos of an authenticated user.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "auth_token, t",
				Usage:       "Github authorization token.",
				Value:       "",
				Destination: &authToken,
				Required:    true,
			},
			cli.IntFlag{
				Name:        "page_size, s",
				Usage:       "Size of each page to load.",
				Value:       cfg.Settings.PageSize,
				Destination: &pageSize,
				Required:    false,
			},
		},
		Action: func(c *cli.Context) {
			var shallContinue bool
			spinLoader := spinner.New(spinner.CharSets[cfg.Spinner.Type], cfg.Spinner.Time*time.Millisecond, spinner.WithHiddenCursor(cfg.Spinner.HideCursor))

			currentPage := 1

			for {
				utilities.ClearTerminalScreen()
				spinLoader.Start()

				userRepos, err := service.GetUserRepos(authToken, pageSize, currentPage)
				if err != nil {
					fmt.Println(err)
					return
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
					return
				}
			}
		},
	}

	return userReposCmd
}

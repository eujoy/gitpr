package widget

import (
	"fmt"

	"github.com/Angelos-Giannis/gitpr/internal/config"
	"github.com/Angelos-Giannis/gitpr/internal/domain"
	"github.com/rivo/tview"
	"github.com/urfave/cli"
)

type userReposService interface {
	GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error)
}

type pullRequestService interface {
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
}

// NewCmd creates a new command to display the details retrieved as widgets in terminal.
func NewCmd(cfg config.Config, userReposService userReposService, pullRequestService pullRequestService) cli.Command {
	var authToken string

	widgetCmd := cli.Command{
		Name:    "widget",
		Aliases: []string{"w"},
		Usage:   "Display a widget based terminal which will include all the details required.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "auth_token, t",
				Usage:       "Github authorization token.",
				Value:       "",
				Destination: &authToken,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) {
			app := tview.NewApplication()

			userRepositories := getAllUserRepoNames(userReposService, authToken, cfg.Settings.PageSize)

			newPrimitive := func(text string) tview.Primitive {
				return tview.NewTextView().
					SetTextAlign(tview.AlignCenter).
					SetText(text)
			}

			pullRequestsList := tview.NewList()
			userReposList := tview.NewList()
			for _, r := range userRepositories {
				userReposList.AddItem(r.Name, "", '-', func() {
					pullRequestsList.Clear()

					prList := getAllPullRequestsForRepo(pullRequestService, authToken, "Angelos-Giannis", r.Name, cfg.Settings.BaseBranch, cfg.Settings.PullRequestState, cfg.Settings.PageSize)
					if len(prList) == 0 {
						pullRequestsList.AddItem("No pull requests found!", "", '-', nil)
					} else {
						for _, pr := range prList {
							pullRequestsList.AddItem(pr.Title, "", '-', nil)
						}
					}
					// pullRequestsList.AddItem("SOme pull request...", "", '-', nil)
				})
			}

			userReposList.AddItem("Quit", "Press to exit", 'q', func() {
				app.Stop()
			})

			grid := tview.NewGrid().
				SetRows(2).
				SetColumns(40, 0).
				SetBorders(true).
				AddItem(newPrimitive("Header"), 0, 0, 1, 2, 0, 0, false)

			grid.AddItem(userReposList, 1, 0, 1, 1, 0, 0, true).
				AddItem(pullRequestsList, 1, 1, 1, 1, 0, 0, false)

			if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
				panic(err)
			}
		},
	}

	return widgetCmd
}

// getAllUserRepoNames retrieve all the repos that the user has access to and returns a list of their names.
func getAllUserRepoNames(userReposService userReposService, authToken string, pageSize int) []domain.Repository {
	var userRepositories []domain.Repository
	currentPage := 1

	for {
		userRepos, err := userReposService.GetUserRepos(authToken, pageSize, currentPage)
		if err != nil {
			fmt.Println(err)
			return []domain.Repository{}
		}

		for _, r := range userRepos.Repositories {
			userRepositories = append(userRepositories, r)
		}

		if len(userRepos.Repositories) < pageSize {
			break
		}

		currentPage++
	}

	return userRepositories
}

// getAllPullRequestsForRepo retrieves all the repositories for a respective service.
func getAllPullRequestsForRepo(pullRequestService pullRequestService, authToken, repoOwner, repository, baseBranch, prState string, pageSize int) []domain.PullRequest {
	var pullRequestsOfRepository []domain.PullRequest
	currentPage := 1

	for {
		pullRequests, err := pullRequestService.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, currentPage)
		if err != nil {
			// fmt.Println(err)
			return []domain.PullRequest{}
		}

		for _, pr := range pullRequests.PullRequests {
			pullRequestsOfRepository = append(pullRequestsOfRepository, pr)
		}

		if len(pullRequests.PullRequests) < pageSize {
			break
		}

		currentPage++
	}

	return pullRequestsOfRepository
}

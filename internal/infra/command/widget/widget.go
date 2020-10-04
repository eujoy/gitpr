package widget

import (
	"fmt"
	"strings"

	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/gdamore/tcell"
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
	authToken := cfg.Clients.Github.Token.DefaultValue

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
				Required:    false,
			},
		},
		Action: func(c *cli.Context) {
			app := tview.NewApplication()

			selectedPrState := cfg.Settings.PullRequestState
			baseBranch := "master"

			// newPrimitive := func(text string) tview.Primitive {
			// 	return tview.NewTextView().
			// 		SetTextAlign(tview.AlignCenter).
			// 		SetText(text)
			// }

			defaultPrStateIndex := 0
			for idx, s := range cfg.Settings.AllowedPullRequestStates {
				if s == cfg.Settings.PullRequestState {
					defaultPrStateIndex = idx
					break
				}
			}

			headerForm := tview.NewForm().
				SetHorizontal(true).
				SetButtonBackgroundColor(tcell.ColorLightGoldenrodYellow).
				SetButtonTextColor(tcell.ColorBlack).
				SetFieldBackgroundColor(tcell.ColorLightGreen).
				SetFieldTextColor(tcell.ColorBlack).
				SetLabelColor(tcell.ColorLightGoldenrodYellow).
				AddDropDown("PR Status", cfg.Settings.AllowedPullRequestStates, defaultPrStateIndex, func(option string, optionIndex int) {
					selectedPrState = option
					return
				}).
				AddInputField("Base Branch", baseBranch, 40, nil, func(text string) {
					baseBranch = text
					return
				}).
				AddPasswordField("Auth Token", authToken, 50, '*', func(text string) {
					authToken = text
					return
				})
			headerForm.SetBorder(true).SetTitle("Change the filters (navigate using tab button)").SetTitleAlign(tview.AlignCenter)

			grid := tview.NewGrid().
				SetRows(5).
				SetColumns(40, 0).
				SetBorders(true).
				AddItem(headerForm, 0, 0, 1, 2, 20, 0, false)

			userRepositories := getAllUserRepoNames(userReposService, authToken, cfg.Settings.PageSize)

			pullRequestsList := tview.NewList().SetSelectedBackgroundColor(tcell.ColorWhiteSmoke)
			userReposList := tview.NewList().SetSelectedBackgroundColor(tcell.ColorLightYellow)
			for _, r := range userRepositories {
				userRepo := r
				access := "Public"
				if userRepo.Private {
					access = "Private"
				}

				userRepoSecondaryText := fmt.Sprintf("%v (%v Stars)", access, userRepo.Stars)
				userReposList.AddItem(userRepo.FullName, userRepoSecondaryText, '-', func() {
					pullRequestsList.Clear()

					details := strings.Split(userRepo.FullName, "/")

					prList := getAllPullRequestsForRepo(pullRequestService, authToken, details[0], details[1], baseBranch, selectedPrState, cfg.Settings.PageSize)
					if len(prList) == 0 {
						pullRequestsList.AddItem("No pull requests found!", "", '-', nil)
					} else {
						for _, pr := range prList {
							approved, pending, requestedChanges, total := 0, 0, 0, 0
							for key := range pr.ReviewStates {
								total++
								if pr.ReviewStates[key] == "APPROVED" {
									approved++
								}
								if pr.ReviewStates[key] == "PENDING" {
									pending++
								}
								if pr.ReviewStates[key] != "APPROVED" && pr.ReviewStates[key] != "PENDING" {
									requestedChanges++
								}
							}

							primaryText := fmt.Sprintf("%v - by %v", pr.Title, pr.Creator.Username)
							secondaryText := fmt.Sprintf("%v | Status: %v | [APPROVED: %v - PENDING: %v - REQUEST CHANGES: %v - TOTAL: %v]",
								pr.HTMLURL,
								pr.State,
								approved,
								pending,
								requestedChanges,
								total,
							)

							pullRequestsList.AddItem(
								fmt.Sprintf(primaryText),
								secondaryText,
								'-',
								nil,
							)
						}
					}
				})
			}

			userReposList.AddItem("Quit", "Press to exit", 'q', func() {
				app.Stop()
			})

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

	// fmt.Println(authToken, repoOwner, repository, baseBranch, prState, pageSize, currentPage)

	for {
		pullRequests, err := pullRequestService.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, currentPage)
		if err != nil {
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

package releasereport

import (
	"fmt"
	"math"
	"time"

	"github.com/briandowns/spinner"
	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/urfave/cli/v2"
)

type service interface {
	GetReleaseList(authToken, repoOwner, repository string, pageSize, pageNumber int) ([]domain.Release, error)
}

type tablePrinter interface {
	PrintReleaseReport(releaseReport domain.ReleaseReport)
}

// NewCmd creates a new command to retrieve pull requests for a repo.
func NewCmd(cfg config.Config, service service, tablePrinter tablePrinter) *cli.Command {
	var authToken, repoOwner, repository string
	var startDateStr, endDateStr string

	releaseReportCmd := cli.Command{
		Name:    "release-report",
		Aliases: []string{"r"},
		Usage:   "Retrieves the releases that were published and/or created within a time range for a repository and prints a report based on them.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "auth_token, t",
				Usage:       "Github authorization token.",
				Value:       cfg.Clients.Github.Token.DefaultValue,
				Destination: &authToken,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "owner, o",
				Usage:       "Owner of the repository to retrieve pull requests for.",
				Value:       "",
				Destination: &repoOwner,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "repository, r",
				Usage:       "Repository name to check.",
				Value:       "",
				Destination: &repository,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "start_date",
				Aliases:     []string{"f"},
				Usage:       "Start date of the time range to check. [Expected format: 'yyyy-mm-dd']",
				Destination: &startDateStr,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "end_date",
				Aliases:     []string{"e"},
				Usage:       "End date of the time range to check. [Expected format: 'yyyy-mm-dd']",
				Destination: &endDateStr,
				Required:    false,
			},

		},
		Action: func(c *cli.Context) error {
			startDate, startDateParseErr := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v 00:00:00", startDateStr))
			if startDateParseErr != nil {
				fmt.Printf("Failed to parse date %q with error : %v\n", startDateStr, startDateParseErr)
			}

			endDate, endDateParseErr := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v 23:59:59", endDateStr))
			if startDateParseErr != nil {
				fmt.Printf("Failed to parse date %q with error : %v\n", endDateStr, endDateParseErr)
			}

			spinLoader := spinner.New(spinner.CharSets[cfg.Spinner.Type], cfg.Spinner.Time*time.Millisecond, spinner.WithHiddenCursor(cfg.Spinner.HideCursor))

			var releaseList []domain.Release
			currentPage := 1
			for {
				spinLoader.Start()
				currentReleaseListPage, err := service.GetReleaseList(authToken, repoOwner, repository, 10, currentPage)
				if err != nil {
					spinLoader.Stop()
					fmt.Println(err)
					return err
				}

				needToBreak := false
				for _, rel := range currentReleaseListPage {
					if (rel.CreatedAt.After(startDate) && rel.CreatedAt.Before(endDate)) || (rel.PublishedAt.After(startDate) && rel.PublishedAt.Before(endDate)) {
						releaseList = append(releaseList, rel)
						continue
					}

					if rel.CreatedAt.Before(startDate) && rel.PublishedAt.Before(startDate) {
						needToBreak = true
					}
				}

				if needToBreak {
					spinLoader.Stop()
					break
				}

				currentPage++
			}

			releaseReport := domain.ReleaseReport{
				NumberOfDraftReleases:     0,
				NumberOfReleasesCreated:   0,
				NumberOfReleasesPublished: 0,
				CreatedToPublishedRatio:   0.0,
			}

			for _, rel := range releaseList {
				if rel.Draft {
					releaseReport.NumberOfDraftReleases++
				}

				if rel.CreatedAt.After(startDate) && rel.CreatedAt.Before(endDate) {
					releaseReport.NumberOfReleasesCreated++
				}

				if rel.PublishedAt.After(startDate) && rel.PublishedAt.Before(endDate) {
					releaseReport.NumberOfReleasesPublished++
				}
			}

			if releaseReport.NumberOfReleasesPublished > 0 {
				releaseReport.CreatedToPublishedRatio = math.Round((float64(releaseReport.NumberOfReleasesCreated)/float64(releaseReport.NumberOfReleasesPublished))*100) / 100
			} else {
				releaseReport.CreatedToPublishedRatio = float64(releaseReport.NumberOfReleasesCreated)
			}

			tablePrinter.PrintReleaseReport(releaseReport)

			return nil
		},
	}

	return &releaseReportCmd
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

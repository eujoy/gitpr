package releasereport

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/urfave/cli/v2"
)

const (
	captionTextPattern                = "Releases report for pattern : %v"
	captionTextPatternServiceInitials = "Releases report for service with initials : %v"

	defaultVersionPattern             = "^(v[\\d]+.[\\d]+.[\\d]+)$"
	versionPatternWithServiceInitials = "^(v[\\d]+.[\\d]+.[\\d]+-(\\w){numOfInitialLetters,numOfInitialLetters})$"
)

type service interface {
	GetReleaseList(authToken, repoOwner, repository string, pageSize, pageNumber int) ([]domain.Release, error)
}

type tablePrinter interface {
	PrintReleaseReport(releaseReport domain.ReleaseReport, captionText string)
}

// NewCmd creates a new command to generate report for release on a repo.
func NewCmd(cfg config.Config, service service, tablePrinter tablePrinter) *cli.Command {
	var authToken, repoOwner, repository string
	var startDateStr, endDateStr string

	var enableDefaultVersionPattern bool
	var enableDefaultVersionPatternWithServiceInitials bool
	var numOfInitialLetters int
	var useVersionPatternWithServiceInitials string

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
			&cli.BoolFlag{
				Name:        "default_version_pattern",
				Aliases:     []string{"dvp"},
				Usage:       fmt.Sprintf("Enables the default release version pattern to be used. (default pattern: %v)", defaultVersionPattern),
				Destination: &enableDefaultVersionPattern,
				Required:    false,
			},
			&cli.IntFlag{
				Name:        "version_pattern_with_service_initials",
				Aliases:     []string{"vpwsi"},
				Usage:       fmt.Sprintf("Enables the release version pattern that uses the provided number of letters for service initials to be used. [pattern format: %v]", versionPatternWithServiceInitials),
				Value:       0,
				Destination: &numOfInitialLetters,
				Required:    false,
			},
		},
		Action: func(c *cli.Context) error {
			if numOfInitialLetters > 0 {
				enableDefaultVersionPatternWithServiceInitials = true
				useVersionPatternWithServiceInitials = strings.Replace(versionPatternWithServiceInitials, "numOfInitialLetters", strconv.Itoa(numOfInitialLetters), -1)
			}

			startDate, startDateParseErr := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v 00:00:00", startDateStr))
			if startDateParseErr != nil {
				fmt.Printf("Failed to parse date %q with error : %v\n", startDateStr, startDateParseErr)
				return startDateParseErr
			}

			endDate, endDateParseErr := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v 23:59:59", endDateStr))
			if endDateParseErr != nil {
				fmt.Printf("Failed to parse date %q with error : %v\n", endDateStr, endDateParseErr)
				return endDateParseErr
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

				if len(currentReleaseListPage) == 0 {
					spinLoader.Stop()
					break
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

			validDefaultReleaseVersion := regexp.MustCompile(defaultVersionPattern)
			validDefaultReleaseVersionWithServiceInitials := regexp.MustCompile(useVersionPatternWithServiceInitials)

			if enableDefaultVersionPatternWithServiceInitials {
				releaseReportMap := make(map[string]*domain.ReleaseReport)
				for _, rel := range releaseList {
					if validDefaultReleaseVersionWithServiceInitials.MatchString(rel.TagName) {
						tagNameSlice := strings.Split(rel.TagName, "-")
						serviceInitials := tagNameSlice[1]

						if _, ok := releaseReportMap[serviceInitials]; !ok {
							releaseReportMap[tagNameSlice[1]] = &domain.ReleaseReport{
								NumberOfDraftReleases:     0,
								NumberOfReleasesCreated:   0,
								NumberOfReleasesPublished: 0,
								CreatedToPublishedRatio:   0.0,
							}
						}

						if rel.Draft {
							releaseReportMap[serviceInitials].NumberOfDraftReleases++
						}

						if rel.PreRelease {
							releaseReportMap[serviceInitials].NumberOfDraftReleases++
						}

						if rel.CreatedAt.After(startDate) && rel.CreatedAt.Before(endDate) {
							releaseReportMap[serviceInitials].NumberOfReleasesCreated++
						}

						if rel.PublishedAt.After(startDate) && rel.PublishedAt.Before(endDate) {
							releaseReportMap[serviceInitials].NumberOfReleasesPublished++
						}
					}
				}

				for serviceInitials, releaseReport := range releaseReportMap {
					releaseReport.CalculateRatioFields()

					captionText := fmt.Sprintf(captionTextPatternServiceInitials, serviceInitials)
					tablePrinter.PrintReleaseReport(*releaseReport, captionText)
				}
			} else {
				releaseReport := &domain.ReleaseReport{
					NumberOfDraftReleases:     0,
					NumberOfReleasesCreated:   0,
					NumberOfReleasesPublished: 0,
					CreatedToPublishedRatio:   0.0,
				}

				for _, rel := range releaseList {
					if enableDefaultVersionPattern {
						if !validDefaultReleaseVersion.MatchString(rel.TagName) {
							continue
						}
					}

					if rel.Draft {
						releaseReport.NumberOfDraftReleases++
					}

					if rel.PreRelease {
						releaseReport.NumberOfDraftReleases++
					}

					if rel.CreatedAt.After(startDate) && rel.CreatedAt.Before(endDate) {
						releaseReport.NumberOfReleasesCreated++
					}

					if rel.PublishedAt.After(startDate) && rel.PublishedAt.Before(endDate) {
						releaseReport.NumberOfReleasesPublished++
					}
				}

				releaseReport.CalculateRatioFields()

				captionText := ""
				if enableDefaultVersionPattern {
					captionText = fmt.Sprintf(captionTextPattern, defaultVersionPattern)
				}
				tablePrinter.PrintReleaseReport(*releaseReport, captionText)
			}

			return nil
		},
	}

	return &releaseReportCmd
}

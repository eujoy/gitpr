package prmetrics

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/urfave/cli/v2"
)

type pullRequestService interface {
	GetPullRequestsCommits(authToken, repoOwner, repository string, pullRequestNumber, pageSize, pageNumber int) ([]domain.Commit, error)
	GetPullRequestsDetails(authToken, repoOwner, repository string, pullRequestNumber int) (domain.PullRequest, error)
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
}

type repositoryService interface {
	GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error)
}

type tablePrinter interface {
	PrintPullRequestLeadTime(pullRequests []domain.PullRequestMetricDetails)
}

type utilities interface {
	ClearTerminalScreen()
	GetPageOptions(respLength int, pageSize int, currentPage int) []string
	GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool)
}

// NewCmd creates a new command to retrieve pull requests for a repo.
func NewCmd(cfg config.Config, pullRequestService pullRequestService, repositoryService repositoryService, tablePrinter tablePrinter, utilities utilities) *cli.Command {
	var authToken, repoOwner, repository, baseBranch, prState string
	var startDateStr, endDateStr string
	var pageSize int
	var printJson bool

	pullRequestsCmd := cli.Command{
		Name:    "pr-metrics",
		Aliases: []string{"m"},
		Usage:   "Retrieves and prints the number of pull requests for a repository that have been created during a specific time period as well as the lead time of those pull requests.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "auth_token",
				Aliases:     []string{"t"},
				Usage:       "Github authorization token.",
				Value:       cfg.Clients.Github.Token.DefaultValue,
				Destination: &authToken,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "owner",
				Aliases:     []string{"o"},
				Usage:       "Owner of the repository to retrieve pull requests for.",
				Value:       "",
				Destination: &repoOwner,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "repository",
				Aliases:     []string{"r"},
				Usage:       "Repository name to check.",
				Value:       "",
				Destination: &repository,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "base",
				Aliases:     []string{"b"},
				Usage:       "Base branch to check pull requests against.",
				Value:       cfg.Settings.BaseBranch,
				Destination: &baseBranch,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "state",
				Aliases:     []string{"a"},
				Usage:       "State of the pull request.",
				Value:       cfg.Settings.PullRequestState,
				Destination: &prState,
				Required:    false,
			},
			&cli.IntFlag{
				Name:        "page_size",
				Aliases:     []string{"s"},
				Usage:       "Size of each page to load.",
				Value:       cfg.Settings.PageSize,
				Destination: &pageSize,
				Required:    false,
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
				Name:        "print_json",
				Aliases:     []string{"json"},
				Usage:       "Define whether the output needs to be printed in json format.",
				Value:       false,
				Destination: &printJson,
				HasBeenSet:  false,
			},
		},
		Action: func(c *cli.Context) error {
			var prLeadTimes []domain.PullRequestMetricDetails
			shallContinue := true
			spinLoader := spinner.New(spinner.CharSets[cfg.Spinner.Type], cfg.Spinner.Time*time.Millisecond, spinner.WithHiddenCursor(cfg.Spinner.HideCursor))

			startDate, startDateParseErr := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v 00:00:00", startDateStr))
			if startDateParseErr != nil {
				fmt.Printf("Failed to parse date %q with error : %v\n", startDateStr, startDateParseErr)
			}

			endDate, endDateParseErr := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v 23:59:59", endDateStr))
			if startDateParseErr != nil {
				fmt.Printf("Failed to parse date %q with error : %v\n", endDateStr, endDateParseErr)
			}

			prState = validatePrStateAndGetDefault(cfg, prState)
			currentPage := 1
			for {
				spinLoader.Start()

				prResp, err := pullRequestService.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, currentPage)
				if err != nil {
					spinLoader.Stop()
					fmt.Println(err)
					return err
				}

				for _, pr := range prResp.PullRequests {
					if pr.CreatedAt.After(startDate) && pr.CreatedAt.Before(endDate) {
						actualLeadTime := time.Duration(0)
						if !pr.MergedAt.IsZero() {
							actualLeadTime = pr.MergedAt.Sub(pr.CreatedAt)
						}

						pullRequestDetails, err := pullRequestService.GetPullRequestsDetails(authToken, repoOwner, repository, pr.Number)
						if err != nil {
							spinLoader.Stop()
							fmt.Println(err)
							return err
						}

						firstCommitsList, err := pullRequestService.GetPullRequestsCommits(authToken, repoOwner, repository, pr.Number, 1, 1)
						if err != nil {
							spinLoader.Stop()
							fmt.Printf("Failed to get detials of first commit with error : %v\n", err)
							return err
						}

						actualTimeToMerge := firstCommitsList[0].Details.Committer.Date.Sub(time.Now())

						if pullRequestDetails.MergeCommitSha != "" {
							lastCommit, err := repositoryService.GetCommitDetails(authToken, repoOwner, repository, pullRequestDetails.MergeCommitSha)
							if err != nil {
								spinLoader.Stop()
								fmt.Printf("Failed to get detials of last commit with error : %v\n", err)
								return err
							}

							actualTimeToMerge = lastCommit.Details.Committer.Date.Sub(firstCommitsList[0].Details.Committer.Date)
						}

						leadTm := domain.PullRequestMetricDetails{
							Number:         pr.Number,
							Title:          pr.Title,
							LeadTime:       actualLeadTime,
							TimeToMerge:    actualTimeToMerge,
							CreatedAt:      pullRequestDetails.CreatedAt,
							Comments:       pullRequestDetails.Comments,
							ReviewComments: pullRequestDetails.ReviewComments,
							Commits:        pullRequestDetails.Commits,
							Additions:      pullRequestDetails.Additions,
							Deletions:      pullRequestDetails.Deletions,
							ChangedFiles:   pullRequestDetails.ChangedFiles,
						}

						prLeadTimes = append(prLeadTimes, leadTm)
					}

					if pr.CreatedAt.Before(startDate) {
						shallContinue = false
					}
				}

				if !shallContinue {
					break
				}
				currentPage++
			}

			spinLoader.Stop()

			if printJson {
				type jsonOutput struct {
					NumOfPullRequests int                               `json:"num_of_pull_requests"`
					PrMetrics         []domain.PullRequestMetricDetails `json:"pr_metrics"`
				}

				jOut := jsonOutput{
					NumOfPullRequests: len(prLeadTimes),
					PrMetrics:         prLeadTimes,
				}

				jsonBytes, err := json.Marshal(jOut)
				if err != nil {
					fmt.Printf("Failed to generate json with error : %v", err)
					return err
				}

				fmt.Printf("%s\n", string(jsonBytes))
			} else {
				fmt.Printf("Number of pull requests : %v\n", len(prLeadTimes))
				fmt.Println()
				tablePrinter.PrintPullRequestLeadTime(prLeadTimes)
			}

			return nil
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

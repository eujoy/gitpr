package prmetrics

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/briandowns/spinner"
    "github.com/eujoy/gitpr/internal/config"
    "github.com/eujoy/gitpr/internal/domain"
    "github.com/eujoy/gitpr/internal/infra/flag"

    "github.com/urfave/cli/v2"
)

const (
    defaultPageSize = 20
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
    PrintPullRequestFlowRatio(flowRatioData map[string]*domain.PullRequestFlowRatio)
    PrintPullRequestMetrics(pullRequests domain.PullRequestMetrics)
}

type utilities interface {
    ClearTerminalScreen()
    GetPageOptions(respLength int, pageSize int, currentPage int) []string
    GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool)
    ConvertDurationToString(dur time.Duration) string
}

// NewCmd creates a new command to retrieve pull requests for a repo.
func NewCmd(cfg config.Config, pullRequestService pullRequestService, repositoryService repositoryService, tablePrinter tablePrinter, utilities utilities) *cli.Command {
    var authToken, repoOwner, repository, baseBranch, prState string
    var startDateStr, endDateStr string
    var printJson bool

    flagBuilder := flag.New(cfg)

    pullRequestsCmd := cli.Command{
        Name:    "pr-metrics",
        Aliases: []string{"m"},
        Usage:   "Retrieves and prints the number of pull requests for a repository that have been created during a specific time period as well as the lead time of those pull requests.",
        Flags: flagBuilder.
            AppendAuthFlag(&authToken).
            AppendOwnerFlag(&repoOwner).
            AppendRepositoryFlag(&repository).
            AppendBaseFlag(&baseBranch).
            AppendStateFlag(&prState).
            AppendStartDateFlag(&startDateStr, false).
            AppendEndDateFlag(&endDateStr, false).
            AppendPrintJsonFlag(&printJson).
            GetFlags(),
        Action: func(c *cli.Context) error {
            prFlowRatio := make(map[string]*domain.PullRequestFlowRatio)

            var prMetricsDetails []domain.PullRequestMetricDetails
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

            totalAggregation := domain.TotalAggregation{
                LeadTime:    time.Duration(0),
                TimeToMerge: time.Duration(0),
            }

            prState = validatePrStateAndGetDefault(cfg, prState)
            currentPage := 1
            for {
                spinLoader.Start()

                prResp, err := pullRequestService.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, defaultPageSize, currentPage)
                if err != nil {
                    spinLoader.Stop()
                    fmt.Println(err)
                    return err
                }

                if len(prResp.PullRequests) == 0 {
                    spinLoader.Stop()
                    fmt.Println("Retrieved empty list.")
                    break
                }

                for _, pr := range prResp.PullRequests {
                    createdAtStr := pr.CreatedAt.Format("2006-01-02")
                    mergedAtStr := ""

                    if pr.MergeCommitSha != "" {
                        if pr.MergedAt.After(startDate) && pr.MergedAt.Before(endDate) {
                            mergedAtStr = pr.MergedAt.Format("2006-01-02")
                            if _, ok := prFlowRatio[mergedAtStr]; !ok {
                                prFlowRatio[mergedAtStr] = &domain.PullRequestFlowRatio{
                                    Created: 0,
                                    Merged:  0,
                                }
                            }

                            prFlowRatio[mergedAtStr].Merged++
                        }
                    }

                    if pr.CreatedAt.After(startDate) && pr.CreatedAt.Before(endDate) {
                        if _, ok := prFlowRatio[createdAtStr]; !ok {
                            prFlowRatio[createdAtStr] = &domain.PullRequestFlowRatio{
                                Created: 0,
                                Merged:  0,
                            }
                        }

                        prFlowRatio[createdAtStr].Created++

                        actualLeadTime := time.Duration(0)
                        if !pr.MergedAt.IsZero() {
                            actualLeadTime = pr.MergedAt.Sub(pr.CreatedAt)
                            totalAggregation.LeadTime += actualLeadTime
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
                            fmt.Printf("Failed to get details of first commit with error : %v\n", err)
                            return err
                        }

                        actualTimeToMerge := time.Until(firstCommitsList[0].Details.Committer.Date)

                        if pullRequestDetails.MergeCommitSha != "" {
                            lastCommit, err := repositoryService.GetCommitDetails(authToken, repoOwner, repository, pullRequestDetails.MergeCommitSha)
                            if err != nil {
                                spinLoader.Stop()
                                fmt.Printf("Failed to get details of last commit with error : %v\n", err)
                                return err
                            }

                            actualTimeToMerge = lastCommit.Details.Committer.Date.Sub(firstCommitsList[0].Details.Committer.Date)
                            totalAggregation.TimeToMerge += actualTimeToMerge
                        }

                        prMetric := domain.PullRequestMetricDetails{
                            Number:         pr.Number,
                            Title:          pr.Title,
                            LeadTime:       actualLeadTime,
                            TimeToMerge:    actualTimeToMerge,
                            StrLeadTime:    utilities.ConvertDurationToString(actualLeadTime),
                            StrTimeToMerge: utilities.ConvertDurationToString(actualTimeToMerge),
                            CreatedAt:      pullRequestDetails.CreatedAt,
                            Comments:       pullRequestDetails.Comments,
                            ReviewComments: pullRequestDetails.ReviewComments,
                            Commits:        pullRequestDetails.Commits,
                            Additions:      pullRequestDetails.Additions,
                            Deletions:      pullRequestDetails.Deletions,
                            ChangedFiles:   pullRequestDetails.ChangedFiles,
                        }
                        updateTotals(&totalAggregation, prMetric)

                        prMetricsDetails = append(prMetricsDetails, prMetric)
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

            totalAggregation.StrLeadTime = utilities.ConvertDurationToString(totalAggregation.LeadTime)
            totalAggregation.StrTimeToMerge = utilities.ConvertDurationToString(totalAggregation.TimeToMerge)
            prMetrics := domain.PullRequestMetrics{
                PRDetails: prMetricsDetails,
                Total:     totalAggregation,
                Average:   calculateAvgAggregation(utilities, len(prMetricsDetails), totalAggregation),
            }

            totalCreated := 0
            totalMerged := 0
            for _, fd := range prFlowRatio {
                totalCreated += fd.Created
                totalMerged += fd.Merged

                ratio := float64(fd.Created) / float64(fd.Merged)
                fd.Ratio = fmt.Sprintf("%.2f", ratio)
            }

            prFlowRatio["Summary"] = &domain.PullRequestFlowRatio{
                Created: totalCreated,
                Merged:  totalMerged,
                Ratio:   fmt.Sprintf("%.2f", float64(totalCreated)/float64(totalMerged)),
            }

            spinLoader.Stop()

            if printJson {
                type jsonOutput struct {
                    NumOfPullRequests int                                     `json:"num_of_pull_requests"`
                    PrMetrics         domain.PullRequestMetrics               `json:"data"`
                    PrFlowRatio       map[string]*domain.PullRequestFlowRatio `json:"flow_ratio"`
                }

                jOut := jsonOutput{
                    NumOfPullRequests: len(prMetricsDetails),
                    PrMetrics:         prMetrics,
                    PrFlowRatio:       prFlowRatio,
                }

                jsonBytes, err := json.Marshal(jOut)
                if err != nil {
                    fmt.Printf("Failed to generate json with error : %v", err)
                    return err
                }

                fmt.Printf("%s\n", string(jsonBytes))
            } else {
                fmt.Printf("Number of pull requests : %v\n", len(prMetricsDetails))
                fmt.Println()
                tablePrinter.PrintPullRequestMetrics(prMetrics)
                fmt.Println()
                tablePrinter.PrintPullRequestFlowRatio(prFlowRatio)
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

func updateTotals(totalData *domain.TotalAggregation, metricDetails domain.PullRequestMetricDetails) {
    totalData.Comments += metricDetails.Comments
    totalData.ReviewComments += metricDetails.ReviewComments
    totalData.Commits += metricDetails.Commits
    totalData.Additions += metricDetails.Additions
    totalData.Deletions += metricDetails.Deletions
    totalData.ChangedFiles += metricDetails.ChangedFiles
}

func calculateAvgAggregation(utilities utilities, prCount int, totalData domain.TotalAggregation) domain.AverageAggregation {
    avgLeadTime := time.Duration(totalData.LeadTime.Seconds()/float64(prCount)) * time.Second
    avgTimeToMerge := time.Duration(totalData.TimeToMerge.Seconds()/float64(prCount)) * time.Second

    return domain.AverageAggregation{
        Comments:       float64(totalData.Comments) / float64(prCount),
        ReviewComments: float64(totalData.ReviewComments) / float64(prCount),
        Commits:        float64(totalData.Commits) / float64(prCount),
        Additions:      float64(totalData.Additions) / float64(prCount),
        Deletions:      float64(totalData.Deletions) / float64(prCount),
        ChangedFiles:   float64(totalData.ChangedFiles) / float64(prCount),
        LeadTime:       avgLeadTime,
        TimeToMerge:    avgTimeToMerge,
        StrLeadTime:    utilities.ConvertDurationToString(avgLeadTime),
        StrTimeToMerge: utilities.ConvertDurationToString(avgTimeToMerge),
    }
}

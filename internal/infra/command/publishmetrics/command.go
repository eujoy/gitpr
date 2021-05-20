package publishmetrics

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/eujoy/gitpr/internal/infra/flag"
	"github.com/urfave/cli/v2"
)

const (
	defaultPageSize                     = 20
	pullRequestSheetNameDefaultTemplate = "OverallData-{repositoryName}"
	releaseSheetNameDefaultTemplate     = "Release-{repositoryName}"

	defaultVersionPattern             = "^(v[\\d]+.[\\d]+.[\\d]+)$"
	versionPatternWithServiceInitials = "^(v[\\d]+.[\\d]+.[\\d]+-(\\w){numOfInitialLetters,numOfInitialLetters})$"
)

type pullRequestService interface {
	GetPullRequestsCommits(authToken, repoOwner, repository string, pullRequestNumber, pageSize, pageNumber int) ([]domain.Commit, error)
	GetPullRequestsDetails(authToken, repoOwner, repository string, pullRequestNumber int) (domain.PullRequest, error)
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
}

type repositoryService interface {
	GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error)
	GetReleaseList(authToken, repoOwner, repository string, pageSize, pageNumber int) ([]domain.Release, error)
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

type googleSheetsService interface{
	CreateAndCleanupOverallSheet(spreadsheetID string, sheetName string) error
	CreateAndCleanupReleaseOverallSheet(spreadsheetID string, sheetName string) error
	WritePullRequestReportData(spreadsheetID string, sheetName string, cellRange string, sprint *domain.SprintSummary, prMetrics *domain.PullRequestMetrics, prFlowRatio *domain.PullRequestFlowRatio) error
	WriteReleaseReportData(spreadsheetID string, sheetName string, cellRange string, sprint *domain.SprintSummary, releaseTagType string, releaseReport *domain.ReleaseReport) error
}

// NewCmd creates a new command to retrieve pull requests for a repo.
func NewCmd(cfg config.Config, pullRequestService pullRequestService, repositoryService repositoryService, tablePrinter tablePrinter, utilities utilities, googleSheetsService googleSheetsService) *cli.Command {
	var authToken, repoOwner, repository, baseBranch, prState string
	var spreadsheetID, prSheetName, sprintSummary, relSheetName string

	var enableDefaultVersionPattern bool
	var enableDefaultVersionPatternWithServiceInitials bool
	var numOfInitialLetters int
	var useVersionPatternWithServiceInitials string

	flagBuilder := flag.New(cfg)

	pullRequestsCmd := cli.Command{
		Name:    "publish-metrics",
		Aliases: []string{"pm"},
		Usage:   "Retrieves the metric details for a list of sprints, prepares the report information for each one of them and publishes the report data the provided google spreadsheet.",
		Flags: flagBuilder.
			AppendAuthFlag(&authToken).
			AppendOwnerFlag(&repoOwner).
			AppendRepositoryFlag(&repository).
			AppendBaseFlag(&baseBranch).
			AppendStateFlag(&prState).
			AppendSpreadsheetID(&spreadsheetID).
			AppendPullRequestSheetName(&prSheetName).
			AppendReleaseSheetName(&relSheetName).
			AppendSprintSummary(&sprintSummary).
			AppendDefaultVersionPatternFlag(&enableDefaultVersionPattern, defaultVersionPattern).
			AppendVersionPatternWithServiceInitialsFlag(&numOfInitialLetters, versionPatternWithServiceInitials).
			GetFlags(),
		Action: func(c *cli.Context) error {
			fmt.Println("Starting the process...")

			if numOfInitialLetters > 0 {
				enableDefaultVersionPatternWithServiceInitials = true
				useVersionPatternWithServiceInitials = strings.Replace(versionPatternWithServiceInitials, "numOfInitialLetters", strconv.Itoa(numOfInitialLetters), -1)
			}

			shallContinue := true

			if prSheetName == "" {
				prSheetName = strings.Replace(pullRequestSheetNameDefaultTemplate, "{repositoryName}", repository, -1)
			}

			if relSheetName == "" {
				relSheetName = strings.Replace(releaseSheetNameDefaultTemplate, "{repositoryName}", repository, -1)
			}

			err := googleSheetsService.CreateAndCleanupOverallSheet(spreadsheetID, prSheetName)
			if err != nil {
				fmt.Println(err)
				return err
			}

			err = googleSheetsService.CreateAndCleanupReleaseOverallSheet(spreadsheetID, prSheetName)
			if err != nil {
				fmt.Println(err)
				return err
			}

			var sprintSummaryList []domain.SprintSummary
			err = json.Unmarshal([]byte(sprintSummary), &sprintSummaryList)
			if err != nil {
				fmt.Println(err)
				return err
			}

			var cleanedSummaryList []domain.SprintSummary
			var startAt, endAt time.Time
			for _, sprint := range sprintSummaryList {
				fmt.Printf("Cleanup of sprint with number : %v\n", sprint.Number)

				if !sprint.StartDate.IsZero() && !sprint.EndDate.IsZero() {
					cleanedSummaryList = append(cleanedSummaryList, sprint)

					if startAt.IsZero() || startAt.After(sprint.StartDate.Time) {
						startAt = sprint.StartDate.Time
					}

					if endAt.IsZero() || endAt.Before(sprint.EndDate.Time) {
						endAt = sprint.EndDate.Time
					}
				}
			}

			pullRequestListPerDay := make(map[string][]domain.PullRequest)
			prState = validatePrStateAndGetDefault(cfg, prState)
			currentPage := 1
			for {
				fmt.Printf("Fetch pull requests for page : %v\n", currentPage)

				prResp, err := pullRequestService.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, defaultPageSize, currentPage)
				if err != nil {
					fmt.Println(err)
					return err
				}

				if len(prResp.PullRequests) == 0 {
					fmt.Println("Retrieved empty list.")
					break
				}

				fmt.Printf("Moving pull requests from response to the map - number of pull requests : %v\n", len(prResp.PullRequests))
				for _, pr := range prResp.PullRequests {
					createdAtStr := pr.CreatedAt.Format("2006-01-02")

					pullRequestListPerDay[createdAtStr] = append(pullRequestListPerDay[createdAtStr], pr)

					if pr.CreatedAt.Before(startAt) {
						shallContinue = false
					}
				}

				if !shallContinue {
					break
				}
				currentPage++
			}

			for _, sprint := range cleanedSummaryList {
				var prMetricsDetails []domain.PullRequestMetricDetails

				fmt.Printf("Prepare report data for sprint with number : %v\n", sprint.Number)

				var prsInSprint []domain.PullRequest
				currentDate := sprint.StartDate.Time
				for {
					fmt.Printf("Combine all pull requests for sprint - current date is : %v\n", currentDate.Format("2006-01-02"))

					prsInSprint = append(prsInSprint, pullRequestListPerDay[currentDate.Format("2006-01-02")]...)
					currentDate = currentDate.Add(24 * time.Hour)
					if currentDate.After(sprint.EndDate.Time) {
						break
					}
				}

				prFlowRatio := make(map[string]*domain.PullRequestFlowRatio)
				totalAggregation := domain.TotalAggregation{
					LeadTime:       time.Duration(0),
					TimeToMerge:    time.Duration(0),
				}

				for _, pr := range prsInSprint {
					fmt.Printf("Prepare report for sprint with id : %v\n", pr.ID)

					createdAtStr := pr.CreatedAt.Format("2006-01-02")
					mergedAtStr := ""

					if pr.MergeCommitSha != "" {
						if pr.MergedAt.After(sprint.StartDate.Time) && pr.MergedAt.Before(sprint.EndDate.Time) {
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

					if pr.CreatedAt.After(sprint.StartDate.Time) && pr.CreatedAt.Before(sprint.EndDate.Time) {
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

						fmt.Printf("Fetch pull request details for sprint with id : %v\n", pr.ID)

						pullRequestDetails, err := pullRequestService.GetPullRequestsDetails(authToken, repoOwner, repository, pr.Number)
						if err != nil {
							fmt.Println(err)
							return err
						}

						fmt.Printf("Fetch pull request first commit for sprint with id : %v\n", pr.ID)

						firstCommitsList, err := pullRequestService.GetPullRequestsCommits(authToken, repoOwner, repository, pr.Number, 1, 1)
						if err != nil {
							fmt.Printf("Failed to get details of first commit with error : %v\n", err)
							return err
						}

						actualTimeToMerge := time.Until(firstCommitsList[0].Details.Committer.Date)

						if pullRequestDetails.MergeCommitSha != "" {
							fmt.Printf("Fetch pull request last commit for sprint with id : %v\n", pr.ID)

							lastCommit, err := repositoryService.GetCommitDetails(authToken, repoOwner, repository, pullRequestDetails.MergeCommitSha)
							if err != nil {
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
				}

				totalAggregation.StrLeadTime = utilities.ConvertDurationToString(totalAggregation.LeadTime)
				totalAggregation.StrTimeToMerge = utilities.ConvertDurationToString(totalAggregation.TimeToMerge)
				prMetrics := &domain.PullRequestMetrics{
					PRDetails: prMetricsDetails,
					Total:     totalAggregation,
					Average:   calculateAvgAggregation(utilities, len(prMetricsDetails), totalAggregation),
				}

				totalCreated := 0
				totalMerged := 0
				for _, fd := range prFlowRatio {
					totalCreated += fd.Created
					totalMerged += fd.Merged

					ratio := float64(fd.Created)/float64(fd.Merged)
					fd.Ratio = fmt.Sprintf("%.2f", ratio)
				}

				prFlowRatio["Summary"] = &domain.PullRequestFlowRatio{
					Created: totalCreated,
					Merged:  totalMerged,
					Ratio:   fmt.Sprintf("%.2f", float64(totalCreated)/float64(totalMerged)),
				}

				err = googleSheetsService.WritePullRequestReportData(spreadsheetID, prSheetName, fmt.Sprintf("A%d", sprint.Number+1), &sprint, prMetrics, prFlowRatio["Summary"])
				if err != nil {
					fmt.Printf("Failed to write report for pull requests for sprint with error : %v\n", err)
					return err
				}
			}

			fmt.Println("Fetch information about the releases.")

			releaseList := make(map[string][]domain.Release)
			currentPage = 1
			for {
				currentReleaseListPage, err := repositoryService.GetReleaseList(authToken, repoOwner, repository, 10, currentPage)
				if err != nil {
					fmt.Println(err)
					return err
				}

				if len(currentReleaseListPage) == 0 {
					break
				}

				needToBreak := false
				for _, rel := range currentReleaseListPage {
					if rel.CreatedAt.After(startAt) && rel.CreatedAt.Before(endAt) {
						createdAtStr := rel.CreatedAt.Format("2006-01-02")
						releaseList[createdAtStr] = append(releaseList[createdAtStr], rel)
						continue
					} else {
						if rel.PublishedAt.After(startAt) && rel.PublishedAt.Before(endAt) {
							publishedAtStr := rel.PublishedAt.Format("2006-01-02")
							releaseList[publishedAtStr] = append(releaseList[publishedAtStr], rel)
							continue
						}
					}

					if rel.CreatedAt.Before(startAt) && rel.PublishedAt.Before(endAt) {
						needToBreak = true
					}
				}

				if needToBreak {
					break
				}

				currentPage++
			}

			for _, sprint := range cleanedSummaryList {
				var relList []domain.Release

				fmt.Printf("Prepare release report data for sprint with number : %v\n", sprint.Number)

				currentDate := sprint.StartDate.Time
				for {
					fmt.Printf("Combine all releases for sprint - current date is : %v\n", currentDate.Format("2006-01-02"))

					relList = append(relList, releaseList[currentDate.Format("2006-01-02")]...)
					currentDate = currentDate.Add(24 * time.Hour)
					if currentDate.After(sprint.EndDate.Time) {
						break
					}
				}

				validDefaultReleaseVersion := regexp.MustCompile(defaultVersionPattern)
				validDefaultReleaseVersionWithServiceInitials := regexp.MustCompile(useVersionPatternWithServiceInitials)

				if enableDefaultVersionPatternWithServiceInitials {
					releaseReportMap := make(map[string]*domain.ReleaseReport)
					for _, rel := range relList {
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

							if rel.CreatedAt.After(sprint.StartDate.Time) && rel.CreatedAt.Before(sprint.EndDate.Time) {
								releaseReportMap[serviceInitials].NumberOfReleasesCreated++
							}

							if rel.PublishedAt.After(sprint.StartDate.Time) && rel.PublishedAt.Before(sprint.EndDate.Time) {
								releaseReportMap[serviceInitials].NumberOfReleasesPublished++
							}
						}
					}

					extraLine := 0
					for serviceInitials, releaseReport := range releaseReportMap {
						releaseReport.CalculateRatioFields()

						err = googleSheetsService.WriteReleaseReportData(spreadsheetID, relSheetName, "A1", &sprint, serviceInitials, releaseReport)
						if err != nil {
							fmt.Printf("Failed to write release report data to spreadsheet with error : %v", err)
							return err
						}

						extraLine++
					}
				} else {
					releaseReport := &domain.ReleaseReport{
						NumberOfDraftReleases:     0,
						NumberOfReleasesCreated:   0,
						NumberOfReleasesPublished: 0,
						CreatedToPublishedRatio:   0.0,
					}

					for _, rel := range relList {
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

						if rel.CreatedAt.After(sprint.StartDate.Time) && rel.CreatedAt.Before(sprint.EndDate.Time) {
							releaseReport.NumberOfReleasesCreated++
						}

						if rel.PublishedAt.After(sprint.StartDate.Time) && rel.PublishedAt.Before(sprint.EndDate.Time) {
							releaseReport.NumberOfReleasesPublished++
						}
					}

					releaseReport.CalculateRatioFields()

					err = googleSheetsService.WriteReleaseReportData(spreadsheetID, relSheetName, "A1", &sprint, "", releaseReport)
					if err != nil {
						fmt.Printf("Failed to write release report data to spreadsheet with error : %v", err)
						return err
					}
				}
			}

			fmt.Println("Finished process successfully!!")

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
	totalData.Comments       += metricDetails.Comments
	totalData.ReviewComments += metricDetails.ReviewComments
	totalData.Commits        += metricDetails.Commits
	totalData.Additions      += metricDetails.Additions
	totalData.Deletions      += metricDetails.Deletions
	totalData.ChangedFiles   += metricDetails.ChangedFiles
}

func calculateAvgAggregation(utilities utilities, prCount int, totalData domain.TotalAggregation) domain.AverageAggregation {
	avgLeadTime := time.Duration(totalData.LeadTime.Seconds()/float64(prCount)) * time.Second
	avgTimeToMerge := time.Duration(totalData.TimeToMerge.Seconds()/float64(prCount)) * time.Second

	return domain.AverageAggregation{
		Comments:       float64(totalData.Comments)/float64(prCount),
		ReviewComments: float64(totalData.ReviewComments)/float64(prCount),
		Commits:        float64(totalData.Commits)/float64(prCount),
		Additions:      float64(totalData.Additions)/float64(prCount),
		Deletions:      float64(totalData.Deletions)/float64(prCount),
		ChangedFiles:   float64(totalData.ChangedFiles)/float64(prCount),
		LeadTime:       avgLeadTime,
		TimeToMerge:    avgTimeToMerge,
		StrLeadTime:    utilities.ConvertDurationToString(avgLeadTime),
		StrTimeToMerge: utilities.ConvertDurationToString(avgTimeToMerge),
	}
}

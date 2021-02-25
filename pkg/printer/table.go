package printer

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/eujoy/gitpr/internal/domain"
	"github.com/jedib0t/go-pretty/v6/table"
)

type total struct {
	comments       int
	reviewComments int
	commits        int
	additions      int
	deletions      int
	changedFiles   int

	leadTime time.Duration
	timeToMerge time.Duration
}

// TablePrinter wraps the printout for models as table.
type TablePrinter struct{}

// NewTablePrinter creates and returns a new table printer struct.
func NewTablePrinter() *TablePrinter {
	return &TablePrinter{}
}

// PrintRepos prints repositories as table.
func (t *TablePrinter) PrintRepos(repos []domain.Repository) {
	outputTable := table.NewWriter()
	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"ID", "Name", "Full Name", "Url", "ssh Url", "Private", "Language", "Stars"})

	for _, r := range repos {
		outputTable.AppendRow(table.Row{r.ID, r.Name, r.FullName, r.HtmlUrl, r.SshUrl, r.Private, r.Language, r.Stars})
	}

	outputTable.AppendSeparator()

	// To find styles : https://github.com/jedib0t/go-pretty/blob/master/table/style.go
	outputTable.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	// outputTable.SetStyle(table.StyleBold)

	// customStyle := table.Style{
	// 	Name:    "CustomStyle",
	// 	Box:     table.StyleBoxBold,
	// 	Color:   table.ColorOptionsCyanWhiteOnBlack,
	// 	Format:  table.FormatOptionsDefault,
	// 	Options: table.OptionsDefault,
	// 	Title:   table.TitleOptionsDefault,
	// }
	// outputTable.SetStyle(customStyle)
	outputTable.Render()
}

// PrintPullRequest prints pull requests as table.
func (t *TablePrinter) PrintPullRequest(pullRequests []domain.PullRequest) {
	outputTable := table.NewWriter()
	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"#", "Url", "Title", "Labels", "State", "Approved", "Req. Changes", "Req. Reviews"})

	for _, p := range pullRequests {
		approved, requestedChanges, total := 0, 0, 0
		for key := range p.ReviewStates {
			total++
			if p.ReviewStates[key] == "APPROVED" {
				approved++
			}
			if p.ReviewStates[key] != "APPROVED" && p.ReviewStates[key] != "PENDING" {
				requestedChanges++
			}
		}

		outputTable.AppendRow(table.Row{p.Number, p.HtmlUrl, p.Title, p.Labels, p.State, approved, requestedChanges, total})
	}

	outputTable.AppendSeparator()
	outputTable.SetStyle(table.StyleBold)
	outputTable.Render()
}

// PrintPullRequestLeadTime prints pull requests lead time as table.
func (t *TablePrinter) PrintPullRequestLeadTime(pullRequests []domain.PullRequestMetricDetails) {
	totalData := total{}

	outputTable := table.NewWriter()
	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"#", "Title", "Comments", "Review Comments", "Commits", "Additions", "Deletions", "Changed Files", "Lead Time", "Time to Merge", "Created At"})

	for _, p := range pullRequests {
		leadTime := ""
		if p.LeadTime != time.Duration(0) {
			totalData.leadTime += p.LeadTime
			leadTime = convertDurationToString(p.LeadTime)
		}

		timeToMerge := ""
		if p.TimeToMerge != time.Duration(0) {
			totalData.timeToMerge += p.TimeToMerge
			timeToMerge = convertDurationToString(p.TimeToMerge)
		}

		outputTable.AppendRow(table.Row{p.Number, p.Title, p.Comments, p.ReviewComments, p.Commits, p.Additions, p.Deletions, p.ChangedFiles, leadTime, timeToMerge, p.CreatedAt})

		updateTotals(&totalData, p)
	}

	totalRow, averageRow := getTotalAndAverageRows(len(pullRequests), totalData)

	outputTable.AppendSeparator()
	outputTable.AppendRow(totalRow)
	outputTable.AppendSeparator()
	outputTable.AppendRow(averageRow)
	outputTable.SetStyle(table.StyleBold)
	outputTable.Render()
}

func updateTotals(totalData *total, metricDetails domain.PullRequestMetricDetails) {
	totalData.comments       += metricDetails.Comments
	totalData.reviewComments += metricDetails.ReviewComments
	totalData.commits        += metricDetails.Commits
	totalData.additions      += metricDetails.Additions
	totalData.deletions      += metricDetails.Deletions
	totalData.changedFiles   += metricDetails.ChangedFiles
}

func convertDurationToString(dur time.Duration) string {
	days := int64(dur.Hours() / 24)
	hours := int64(math.Mod(dur.Hours(), 24))
	minutes := int64(math.Mod(dur.Minutes(), 60))
	seconds := int64(math.Mod(dur.Seconds(), 60))

	formattedDuration := fmt.Sprintf("%d days & %02d:%02d:%02d", days, hours, minutes, seconds)

	return formattedDuration
}

func getTotalAndAverageRows(totalPullRequests int, totalData total) (table.Row, table.Row) {
	totalPullRequestsFloat := float64(totalPullRequests)

	totalTableRow := table.Row{
		"",
		"Total",
		totalData.comments,
		totalData.reviewComments,
		totalData.commits,
		totalData.additions,
		totalData.deletions,
		totalData.changedFiles,
		convertDurationToString(totalData.leadTime),
		convertDurationToString(totalData.timeToMerge),
		"",
	}

	averageTableRow := table.Row{
		"",
		"Average",
		fmt.Sprintf("%.2f", float64(totalData.comments)/totalPullRequestsFloat),
		fmt.Sprintf("%.2f", float64(totalData.reviewComments)/totalPullRequestsFloat),
		fmt.Sprintf("%.2f", float64(totalData.commits)/totalPullRequestsFloat),
		fmt.Sprintf("%.2f", float64(totalData.additions)/totalPullRequestsFloat),
		fmt.Sprintf("%.2f", float64(totalData.deletions)/totalPullRequestsFloat),
		fmt.Sprintf("%.2f", float64(totalData.changedFiles)/totalPullRequestsFloat),
		"",
		"",
		"",
	}

	return totalTableRow, averageTableRow
}
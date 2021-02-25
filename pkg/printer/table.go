package printer

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/eujoy/gitpr/internal/domain"
	"github.com/jedib0t/go-pretty/v6/table"
)

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
	outputTable := table.NewWriter()
	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"#", "Title", "Comments", "Review Comments", "Commits", "Additions", "Deletions", "Changed Files", "Lead Time", "Time to Merge", "Created At"})

	for _, p := range pullRequests {
		leadTime := ""
		if p.LeadTime != time.Duration(0) {
			days := int64(p.LeadTime.Hours() / 24)
			hours := int64(math.Mod(p.LeadTime.Hours(), 24))
			minutes := int64(math.Mod(p.LeadTime.Minutes(), 60))
			seconds := int64(math.Mod(p.LeadTime.Seconds(), 60))

			leadTime = fmt.Sprintf("%d days & %d:%d:%d", days, hours, minutes, seconds)
		}

		timeToMerge := ""
		if p.TimeToMerge != time.Duration(0) {
			days := int64(p.TimeToMerge.Hours() / 24)
			hours := int64(math.Mod(p.TimeToMerge.Hours(), 24))
			minutes := int64(math.Mod(p.TimeToMerge.Minutes(), 60))
			seconds := int64(math.Mod(p.TimeToMerge.Seconds(), 60))

			timeToMerge = fmt.Sprintf("%d days & %d:%d:%d", days, hours, minutes, seconds)
		}

		outputTable.AppendRow(table.Row{p.Number, p.Title, p.Comments, p.ReviewComments, p.Commits, p.Additions, p.Deletions, p.ChangedFiles, leadTime, timeToMerge, p.CreatedAt})
	}

	outputTable.AppendSeparator()
	outputTable.SetStyle(table.StyleBold)
	outputTable.Render()
}

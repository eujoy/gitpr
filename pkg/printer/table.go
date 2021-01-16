package printer

import (
	"os"

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

package printer

import (
    "fmt"
    "os"
    "sort"

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

// PrintPullRequestFlowRatio prints pull request flow ratio details as table.
func (t *TablePrinter) PrintPullRequestFlowRatio(flowRatioData map[string]*domain.PullRequestFlowRatio) {
    outputTable := table.NewWriter()
    outputTable.SetOutputMirror(os.Stdout)
    outputTable.AppendHeader(table.Row{"Date", "Created", "Merged", "Ratio"})

    dates := make([]string, 0, len(flowRatioData))
    for k := range flowRatioData {
        dates = append(dates, k)
    }
    sort.Strings(dates)

    for _, d := range dates {
        if d == "Summary" {
            continue
        }
        fd := flowRatioData[d]
        outputTable.AppendRow(table.Row{d, fd.Created, fd.Merged, fd.Ratio})
    }

    outputTable.AppendSeparator()
    outputTable.AppendRow(table.Row{"Summary", flowRatioData["Summary"].Created, flowRatioData["Summary"].Merged, flowRatioData["Summary"].Ratio})
    outputTable.AppendSeparator()
    outputTable.SetStyle(table.StyleBold)
    outputTable.Render()
}

// PrintPullRequestMetrics prints pull request related metrics as table.
func (t *TablePrinter) PrintPullRequestMetrics(pullRequests domain.PullRequestMetrics) {
    outputTable := table.NewWriter()
    outputTable.SetOutputMirror(os.Stdout)
    outputTable.AppendHeader(table.Row{"#", "Title", "Comments", "Review Comments", "Commits", "Additions", "Deletions", "Changed Files", "Lead Time", "Time to Merge", "Created At"})

    for _, p := range pullRequests.PRDetails {
        outputTable.AppendRow(table.Row{p.Number, p.Title, p.Comments, p.ReviewComments, p.Commits, p.Additions, p.Deletions, p.ChangedFiles, p.StrLeadTime, p.StrTimeToMerge, p.CreatedAt})
    }

    totalRow, averageRow := t.getTotalAndAverageRows(pullRequests.Total, pullRequests.Average)

    outputTable.AppendSeparator()
    outputTable.AppendRow(totalRow)
    outputTable.AppendSeparator()
    outputTable.AppendRow(averageRow)
    outputTable.SetStyle(table.StyleBold)
    outputTable.Render()
}

// PrintReleaseReport prints release report details as table.
func (t *TablePrinter) PrintReleaseReport(releaseReport domain.ReleaseReport, captionText string) {
    outputTable := table.NewWriter()
    outputTable.SetOutputMirror(os.Stdout)

    outputTable.AppendHeader(table.Row{"Label", "Value"})
    outputTable.AppendRow(table.Row{"Draft", releaseReport.NumberOfDraftReleases})
    outputTable.AppendRow(table.Row{"Pre-Releases", releaseReport.NumberOfPreReleases})
    outputTable.AppendRow(table.Row{"Created", releaseReport.NumberOfReleasesCreated})
    outputTable.AppendRow(table.Row{"Published", releaseReport.NumberOfReleasesPublished})
    outputTable.AppendSeparator()
    outputTable.AppendRow(table.Row{"Ratio", releaseReport.CreatedToPublishedRatio})

    if captionText != "" {
        outputTable.SetCaption(fmt.Sprintf("%v", captionText))
    }

    outputTable.SetStyle(table.StyleBold)
    outputTable.Render()
}

// PrintWorkflowCosts prints cost details of workflows.
func (t *TablePrinter) PrintWorkflowCosts(workflowBilling []domain.WorkflowBilling) {
    outputTable := table.NewWriter()
    outputTable.SetOutputMirror(os.Stdout)

    totalMinutes := int64(0)
    totalCost := float32(0)
    outputTable.AppendHeader(table.Row{"Workflow", "Env", "Billable Minutes", "Cost"})
    for _, wfBill := range workflowBilling {
        nonZeroEnvs := 0
        wfTotalMinutes := int64(0)
        wfTotalCost := float32(0)

        outputTable.AppendRow(table.Row{wfBill.Name, "", "", ""})
        for _, wfCost := range wfBill.WorkflowCosts {
            if wfCost.Cost > 0 {
                outputTable.AppendRow(table.Row{"", wfCost.EnvType, wfCost.ExecMinutes, fmt.Sprintf("%.2f", wfCost.Cost)})

                totalMinutes += wfCost.ExecMinutes
                totalCost += wfCost.Cost

                wfTotalMinutes += wfCost.ExecMinutes
                wfTotalCost += wfCost.Cost

                nonZeroEnvs++
            }
        }

        if nonZeroEnvs > 1 {
            outputTable.AppendSeparator()
            outputTable.AppendRow(table.Row{"", "", wfTotalMinutes, fmt.Sprintf("%.2f", wfTotalCost)})
        }

        outputTable.AppendSeparator()
    }

    outputTable.AppendFooter(table.Row{"Total", "", totalMinutes, fmt.Sprintf("%.2f", totalCost)})

    outputTable.SetStyle(table.StyleBold)
    outputTable.Render()
}

func (t *TablePrinter) getTotalAndAverageRows(totalData domain.TotalAggregation, averageData domain.AverageAggregation) (table.Row, table.Row) {
    totalTableRow := table.Row{
        "",
        "Total",
        totalData.Comments,
        totalData.ReviewComments,
        totalData.Commits,
        totalData.Additions,
        totalData.Deletions,
        totalData.ChangedFiles,
        totalData.StrLeadTime,
        totalData.StrTimeToMerge,
        "",
    }

    averageTableRow := table.Row{
        "",
        "Average",
        fmt.Sprintf("%.2f", averageData.Comments),
        fmt.Sprintf("%.2f", averageData.ReviewComments),
        fmt.Sprintf("%.2f", averageData.Comments),
        fmt.Sprintf("%.2f", averageData.Additions),
        fmt.Sprintf("%.2f", averageData.Deletions),
        fmt.Sprintf("%.2f", averageData.ChangedFiles),
        averageData.StrLeadTime,
        averageData.StrTimeToMerge,
        "",
    }

    return totalTableRow, averageTableRow
}

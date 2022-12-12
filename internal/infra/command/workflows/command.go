package workflows

import (
    "fmt"
    "sort"
    "time"

    "github.com/briandowns/spinner"
    "github.com/eujoy/gitpr/internal/config"
    "github.com/eujoy/gitpr/internal/domain"
    "github.com/eujoy/gitpr/internal/infra/flag"
    "github.com/urfave/cli/v2"
)

type workflowEnvCost struct {
    Ubuntu  int64
    MacOs   int64
    Windows int64
}

type service interface {
    GetWorkflowExecutions(authToken, repoOwner, repository, startDateStr, endDateStr string, pageSize, pageNumber int) ([]domain.Workflow, error)
    GetWorkflowsOfRepository(authToken, repoOwner, repository string) ([]domain.Workflow, error)
    GetWorkflowTiming(authToken, repoOwner, repository string, runID int) (domain.WorkflowTiming, error)
    GetWorkflowUsage(authToken, repoOwner, repository string, workflowID int) (domain.WorkflowTiming, error)
}

type tablePrinter interface {
    PrintWorkflowCosts(workflowBilling []domain.WorkflowBilling)
}

type utilities interface {
    ClearTerminalScreen()
    GetPageOptions(respLength int, pageSize int, currentPage int) []string
    GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool)
}

// NewCmd creates a new command to retrieve the repos of a user.
func NewCmd(cfg config.Config, service service, tablePrinter tablePrinter, utilities utilities) *cli.Command {
    var authToken, repoOwner, repository string
    var startDateStr, endDateStr string
    var pageSize int

    flagBuilder := flag.New(cfg)
    runsSubCommandFlagBuilder := flag.New(cfg)

    userReposCmd := cli.Command{
        Name:    "workflows",
        Aliases: []string{"wf_exec"},
        Usage:   "Retrieves and prints the workflow executions of a repository.",
        Flags: flagBuilder.
            AppendAuthFlag(&authToken).
            AppendOwnerFlag(&repoOwner).
            AppendRepositoryFlag(&repository).
            GetFlags(),
        Subcommands: []*cli.Command{
            {
                Name:    "runs",
                Aliases: []string{"r"},
                Usage:   "Retrieves and prints the runs of a repository for a given amount of time.",
                Flags: runsSubCommandFlagBuilder.
                    AppendStartDateFlag(&startDateStr, true).
                    AppendEndDateFlag(&endDateStr, true).
                    AppendPageSizeFlag(&pageSize, 100).
                    GetFlags(),
                Action: func(c *cli.Context) error {
                    spinLoader := spinner.New(spinner.CharSets[cfg.Spinner.Type], cfg.Spinner.Time*time.Millisecond, spinner.WithHiddenCursor(cfg.Spinner.HideCursor))

                    workflowTotalTiming := make(map[string]*workflowEnvCost)
                    distinctWorkflows := make(map[string]struct{})

                    currentPage := 1

                    utilities.ClearTerminalScreen()
                    spinLoader.Start()

                    totalWFExecutions := 0

                    for {
                        workflows, err := service.GetWorkflowExecutions(authToken, repoOwner, repository, startDateStr, endDateStr, pageSize, currentPage)
                        if err != nil {
                            fmt.Println(err)
                            return err
                        }

                        totalWFExecutions += len(workflows)

                        if len(workflows) == 0 {
                            break
                        }

                        for _, wf := range workflows {
                            wfTiming, err := service.GetWorkflowTiming(authToken, repoOwner, repository, wf.ID)
                            if err != nil {
                                fmt.Println(err)
                                return err
                            }

                            distinctWorkflows[wf.Name] = struct{}{}

                            if _, exists := workflowTotalTiming[wf.Name]; !exists {
                                workflowTotalTiming[wf.Name] = &workflowEnvCost{Ubuntu: 0, MacOs: 0, Windows: 0}
                            }

                            workflowTotalTiming[wf.Name].Ubuntu += wfTiming.Billable.Ubuntu.TotalMs
                            workflowTotalTiming[wf.Name].MacOs += wfTiming.Billable.MacOs.TotalMs
                            workflowTotalTiming[wf.Name].Windows += wfTiming.Billable.Windows.TotalMs
                        }

                        currentPage++
                    }

                    var distinctWorkflowsList []string
                    for n := range distinctWorkflows {
                        distinctWorkflowsList = append(distinctWorkflowsList, n)
                    }

                    sort.Slice(distinctWorkflowsList, func(i int, j int) bool {
                        return distinctWorkflowsList[i] < distinctWorkflowsList[j]
                    })

                    var wfBillingList []domain.WorkflowBilling
                    for _, wfName := range distinctWorkflowsList {
                        ubuntuMinutes := workflowTotalTiming[wfName].Ubuntu / 60000
                        macOsMinutes := workflowTotalTiming[wfName].MacOs / 60000
                        windowsMinutes := workflowTotalTiming[wfName].Windows / 60000

                        wfCosts := []domain.WorkflowCosts{
                            {
                                EnvType:     "Ubuntu",
                                ExecMinutes: ubuntuMinutes,
                                Cost:        float32(ubuntuMinutes) * cfg.Clients.Github.Billing.Linux,
                            },
                            {
                                EnvType:     "MacOs",
                                ExecMinutes: macOsMinutes,
                                Cost:        float32(macOsMinutes) * cfg.Clients.Github.Billing.MacOS,
                            },
                            {
                                EnvType:     "Windows",
                                ExecMinutes: windowsMinutes,
                                Cost:        float32(windowsMinutes) * cfg.Clients.Github.Billing.Windows,
                            },
                        }

                        wfBillingList = append(wfBillingList, domain.WorkflowBilling{
                            Name:          wfName,
                            WorkflowCosts: wfCosts,
                        })
                    }

                    spinLoader.Stop()
                    utilities.ClearTerminalScreen()

                    tablePrinter.PrintWorkflowCosts(wfBillingList)

                    return nil
                },
            },
            {
                Name:    "billable",
                Aliases: []string{"b"},
                Usage:   "Retrieves and prints summary of billing for all workflows of a repository duting the current billing cycle.",
                Action: func(c *cli.Context) error {
                    spinLoader := spinner.New(spinner.CharSets[cfg.Spinner.Type], cfg.Spinner.Time*time.Millisecond, spinner.WithHiddenCursor(cfg.Spinner.HideCursor))

                    utilities.ClearTerminalScreen()
                    spinLoader.Start()

                    workflows, err := service.GetWorkflowsOfRepository(authToken, repoOwner, repository)
                    if err != nil {
                        fmt.Println(err)
                        return err
                    }

                    var wfBillingList []domain.WorkflowBilling
                    for _, wf := range workflows {
                        usage, err := service.GetWorkflowUsage(authToken, repoOwner, repository, wf.ID)
                        if err != nil {
                            fmt.Println(err)
                            return err
                        }

                        ubuntuMinutes := usage.Billable.Ubuntu.TotalMs / 60000
                        macOsMinutes := usage.Billable.MacOs.TotalMs / 60000
                        windowsMinutes := usage.Billable.Windows.TotalMs / 60000

                        wfCosts := []domain.WorkflowCosts{
                            {
                                EnvType:     "Ubuntu",
                                ExecMinutes: ubuntuMinutes,
                                Cost:        float32(ubuntuMinutes) * cfg.Clients.Github.Billing.Linux,
                            },
                            {
                                EnvType:     "MacOs",
                                ExecMinutes: macOsMinutes,
                                Cost:        float32(macOsMinutes) * cfg.Clients.Github.Billing.MacOS,
                            },
                            {
                                EnvType:     "Windows",
                                ExecMinutes: windowsMinutes,
                                Cost:        float32(windowsMinutes) * cfg.Clients.Github.Billing.Windows,
                            },
                        }

                        wfBillingList = append(wfBillingList, domain.WorkflowBilling{
                            Name:          wf.Name,
                            WorkflowCosts: wfCosts,
                        })
                    }

                    spinLoader.Stop()
                    utilities.ClearTerminalScreen()

                    tablePrinter.PrintWorkflowCosts(wfBillingList)

                    return nil
                },
            },
        },
    }

    return &userReposCmd
}

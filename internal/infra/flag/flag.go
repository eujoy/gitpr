package flag

import (
    "fmt"

    "github.com/eujoy/gitpr/internal/config"
    "github.com/urfave/cli/v2"
)

// builder describes the flag builder definition.
type builder struct{
    cfg            config.Config
    flagDefinition []cli.Flag
}

// New creates and returns a new flag builder.
func New(cfg config.Config) *builder {
    return &builder{
        cfg: cfg,
        flagDefinition: []cli.Flag{},
    }
}

// GetFlags returns all the flags that have been assigned in the list.
func (b *builder) GetFlags() []cli.Flag {
    return b.flagDefinition
}

// AppendAuthFlag appends the 'auth_token' flag in the flag list.
func (b *builder) AppendAuthFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "auth_token",
            Aliases:     []string{"t"},
            Usage:       "Github authorization token.",
            Value:       b.cfg.Clients.Github.Token.DefaultValue,
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendOwnerFlag appends the 'owner' flag in the flag list.
func (b *builder) AppendOwnerFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "owner",
            Aliases:     []string{"o"},
            Usage:       "Owner of the repository to use.",
            Value:       "",
            Destination: destination,
            Required:    true,
        },
    )

    return b
}

// AppendRepositoryFlag appends the 'repository' flag in the flag list.
func (b *builder) AppendRepositoryFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "repository",
            Aliases:     []string{"r"},
            Usage:       "Repository name to use.",
            Value:       "",
            Destination: destination,
            Required:    true,
        },
    )

    return b
}

// AppendBaseFlag appends the 'base' flag in the flag list.
func (b *builder) AppendBaseFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "base",
            Aliases:     []string{"b"},
            Usage:       "Base branch to check pull requests against.",
            Value:       b.cfg.Settings.BaseBranch,
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendStateFlag appends the 'state' flag in the flag list.
func (b *builder) AppendStateFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "state",
            Aliases:     []string{"a"},
            Usage:       "State of the pull request.",
            Value:       b.cfg.Settings.PullRequestState,
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendStartDateFlag appends the 'start_date' flag in the flag list.
func (b *builder) AppendStartDateFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "start_date",
            Aliases:     []string{"f"},
            Usage:       "Start date of the time range to check. [Expected format: 'yyyy-mm-dd']",
            Value:       "",
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendEndDateFlag appends the 'end_date' flag in the flag list.
func (b *builder) AppendEndDateFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "end_date",
            Aliases:     []string{"e"},
            Usage:       "End date of the time range to check. [Expected format: 'yyyy-mm-dd']",
            Value:       "",
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendStartTagFlag appends the 'start_tag' flag in the flag list.
func (b *builder) AppendStartTagFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "start_tag, s",
            Usage:       "The starting tag/commit to compare against.",
            Value:       "",
            Destination: destination,
            Required:    true,
        },
    )

    return b
}

// AppendEndTagFlag appends the 'start_tag' flag in the flag list.
func (b *builder) AppendEndTagFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "end_tag, e",
            Usage:       "The ending/latest tag/commit to compare against.",
            Value:       "HEAD",
            Destination: destination,
            Required:    true,
        },
    )

    return b
}

// AppendReleaseNameFlag appends the 'release_name' flag in the flag list.
func (b *builder) AppendReleaseNameFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "release_name",
            Aliases:     []string{"n"},
            Usage:       "Define the release name to be set. You can use a string pattern to set the place where the new release tag will be set.",
            Value:       "Release version : %v",
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendLatestTagFlag appends the 'latest_tag' flag in the flag list.
func (b *builder) AppendLatestTagFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "latest_tag",
            Aliases:     []string{"l"},
            Usage:       "The latest tag to compare against.",
            Value:       "",
            Destination: destination,
            Required:    true,
        },
    )

    return b
}

// AppendReleaseTagFlag appends the 'release_tag' flag in the flag list.
func (b *builder) AppendReleaseTagFlag(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "release_tag",
            Aliases:     []string{"v"},
            Usage:       "Release tag to be used (and checked against if exists).",
            Value:       "HEAD",
            Destination: destination,
            Required:    true,
        },
    )

    return b
}

// AppendCheckPatternFlag appends the 'check_pattern' flag in the flag list.
func (b *builder) AppendCheckPatternFlag(destination *cli.StringSlice) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringSliceFlag{
            Name:        "check_pattern",
            Aliases:     []string{"p"},
            Usage:       "Define the pattern to check the files modified against.",
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendSpreadsheetID appends the 'spreadsheet_id' flag in the flag list.
func (b *builder) AppendSpreadsheetID(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "spreadsheet_id",
            Aliases:     []string{"sid"},
            Usage:       "Define the id of the spreadsheet to publish the data to.",
            Destination: destination,
            Required:    true,
        },
    )

    return b
}

// AppendSheetName appends the 'sheet_name' flag in the flag list.
func (b *builder) AppendSheetName(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "sheet_name",
            Aliases:     []string{"sheet"},
            Usage:       "Define the name of the sheet to store the report data to. By default, it will be 'OverallData-{repositoryName}'",
            Value:       "",
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendSprintSummary appends the 'sprint_summary' flag in the flag list.
func (b *builder) AppendSprintSummary(destination *string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.StringFlag{
            Name:        "sprint_summary",
            Aliases:     []string{"sprints"},
            Usage:       "Define the summary level details of the sprints to take into account. Expecting json formatted array of sprint summary details. Example : [{\"number\":1,\"name\":\"My Sprint 1\",\"start_date\":\"2020-08-03\",\"end_date\":\"2020-08-17\"},{\"number\":2,\"name\":\"My Sprint 2\",\"start_date\":\"2020-09-08\",\"end_date\":\"2020-09-22\"}]",
            Destination: destination,
            Required:    true,
        },
    )

    return b
}

// AppendPageSizeFlag appends the 'page_size' flag in the flag list.
func (b *builder) AppendPageSizeFlag(destination *int) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.IntFlag{
            Name:        "page_size",
            Aliases:     []string{"s"},
            Usage:       "Size of each page to load.",
            Value:       b.cfg.Settings.PageSize,
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendVersionPatternWithServiceInitials appends the 'version_pattern_with_service_initials' flag in the flag list.
func (b *builder) AppendVersionPatternWithServiceInitialsFlag(destination *int, versionPatternWithServiceInitials string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.IntFlag{
            Name:        "version_pattern_with_service_initials",
            Aliases:     []string{"vpwsi"},
            Usage:       fmt.Sprintf("Enables the release version pattern that uses the provided number of letters for service initials to be used. [pattern format: %v]", versionPatternWithServiceInitials),
            Value:       0,
            Destination: destination,
            Required:    false,
        },
    )

    return b
}

// AppendPrintJsonFlag appends the 'print_json' flag in the flag list.
func (b *builder) AppendPrintJsonFlag(destination *bool) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.BoolFlag{
            Name:        "print_json",
            Aliases:     []string{"json"},
            Usage:       "Define whether the output needs to be printed in json format.",
            Value:       false,
            Destination: destination,
            HasBeenSet:  false,
            Required:    false,
        },
    )

    return b
}

// AppendDraftReleaseFlag appends the 'draft_release' flag in the flag list.
func (b *builder) AppendDraftReleaseFlag(destination *bool) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.BoolFlag{
            Name:        "draft_release",
            Aliases:     []string{"d"},
            Usage:       "Defines if the release will be a draft or published. (default: false)",
            Destination: destination,
            HasBeenSet:  false,
            Required:    false,
        },
    )

    return b
}

// AppendForceCreateFlag appends the 'force_create' flag in the flag list.
func (b *builder) AppendForceCreateFlag(destination *bool) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.BoolFlag{
            Name:        "force_create",
            Aliases:     []string{"f"},
            Usage:       "Forces the creation of the release without asking for confirmation. (default: false)",
            Destination: destination,
            HasBeenSet:  false,
            Required:    false,
        },
    )

    return b
}

// AppendDefaultVersionPatternFlag appends the 'default_version_pattern' flag in the flag list.
func (b *builder) AppendDefaultVersionPatternFlag(destination *bool, defaultVersionPattern string) *builder {
    b.flagDefinition = append(
        b.flagDefinition,
        &cli.BoolFlag{
            Name:        "default_version_pattern",
            Aliases:     []string{"dvp"},
            Usage:       fmt.Sprintf("Enables the default release version pattern to be used. (default pattern: %v)", defaultVersionPattern),
            Destination: destination,
            HasBeenSet:  false,
            Required:    false,
        },
    )

    return b
}

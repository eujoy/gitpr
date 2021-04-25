package flag

import (
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
            Usage:       "Owner of the repository to retrieve pull requests for.",
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
            Usage:       "Repository name to check.",
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

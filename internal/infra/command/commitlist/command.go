package commitlist

import (
	"fmt"
	"os"

	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
	"github.com/eujoy/gitpr/internal/infra/flag"
	"github.com/urfave/cli/v2"
)

type service interface {
	GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error)
	PrintCommitList(commitList []domain.Commit, useTmpl string) (string, error)
}

// NewCmd creates a new command to retrieve the commits between 2 provided tags or commits.
func NewCmd(cfg config.Config, service service) *cli.Command {
	var authToken, repoOwner, repository, startTag, endTag string

	flagBuilder := flag.New(cfg)

	commitListCmd := cli.Command{
		Name:    "commit-list",
		Aliases: []string{"c"},
		Usage:   "Retrieves and prints the list of commits between two provided tags or commits.",
		Flags: flagBuilder.
			AppendAuthFlag(&authToken).
			AppendOwnerFlag(&repoOwner).
			AppendRepositoryFlag(&repository).
			AppendStartTagFlag(&startTag).
			AppendEndTagFlag(&endTag).
			GetFlags(),
		Action: func(c *cli.Context) error {
			commitList, err := service.GetDiffBetweenTags(authToken, repoOwner, repository, startTag, endTag)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			commitListPrintout, err := service.PrintCommitList(commitList.Commits, domain.CommitListTerminalTemplate)
			if err != nil {
				os.Exit(1)
			}

			fmt.Println(commitListPrintout)

			return nil
		},
	}

	return &commitListCmd
}

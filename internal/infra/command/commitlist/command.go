package commitlist

import (
    "fmt"
    "os"

    "github.com/eujoy/gitpr/internal/config"
    "github.com/eujoy/gitpr/internal/domain"
    "github.com/urfave/cli/v2"
)

type service interface {
    GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error)
    PrintCommitList(commitList []domain.Commit, useTmpl string) (string, error)
}

// NewCmd creates a new command to retrieve the commits between 2 provided tags or commits.
func NewCmd(cfg config.Config, service service) *cli.Command {
    var authToken, repoOwner, repository, startTag, endTag string

    commitListCmd := cli.Command{
        Name:    "commit-list",
        Aliases: []string{"c"},
        Usage:   "Retrieves and prints the list of commits between two provided tags or commits.",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:        "auth_token, t",
                Usage:       "Github authorization token.",
                Value:       cfg.Clients.Github.Token.DefaultValue,
                Destination: &authToken,
                Required:    false,
            },
            &cli.StringFlag{
                Name:        "owner, o",
                Usage:       "Owner of the repository to retrieve pull requests for.",
                Value:       "",
                Destination: &repoOwner,
                Required:    true,
            },
            &cli.StringFlag{
                Name:        "repository, r",
                Usage:       "Repository name to check.",
                Value:       "",
                Destination: &repository,
                Required:    true,
            },
            &cli.StringFlag{
                Name:        "start_tag, s",
                Usage:       "The starting tag/commit to compare against.",
                Value:       "",
                Destination: &startTag,
                Required:    true,
            },
            &cli.StringFlag{
                Name:        "end_tag, e",
                Usage:       "The ending/latest tag/commit to compare against.",
                Value:       "HEAD",
                Destination: &endTag,
                Required:    true,
            },
        },
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

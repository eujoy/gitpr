package createrelease

import (
    "fmt"
    "os"
    "strings"

    "github.com/AlecAivazis/survey/v2"
    "github.com/eujoy/gitpr/internal/config"
    "github.com/eujoy/gitpr/internal/domain"
    "github.com/urfave/cli/v2"
)

var promptMessages = map[bool]string{
    true: "Are you sure you want to create the draft release?",
    false: "Are you sure you want to create the release?",
}

type service interface {
    CreateRelease(authToken, repoOwner, repository, tagName string, draftRelease bool, name, body string) error
    GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error)
    GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error)
    PrintCommitList(commitList []domain.Commit, useTmpl string) (string, error)
}

// NewCmd creates a new command to retrieve the commits between 2 provided tags or commits.
func NewCmd(cfg config.Config, service service) *cli.Command {
    var authToken, repoOwner, repository, latestTag, releaseTag, releaseName string
    var draftRelease bool
    var checkPattern cli.StringSlice

    forceCreate := false

    commitListCmd := cli.Command{
        Name:    "create-release",
        Aliases: []string{"cr"},
        Usage:   "Retrieves and prints the list of commits between two provided tags or commits.",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:        "auth_token",
                Aliases:     []string{"t"},
                Usage:       "Github authorization token.",
                Value:       cfg.Clients.Github.Token.DefaultValue,
                Destination: &authToken,
                Required:    false,
            },
            &cli.StringFlag{
                Name:        "owner",
                Aliases:     []string{"o"},
                Usage:       "Owner of the repository to retrieve pull requests for.",
                Value:       "",
                Destination: &repoOwner,
                Required:    true,
            },
            &cli.StringFlag{
                Name:        "repository",
                Aliases:     []string{"r"},
                Usage:       "Repository name to check.",
                Value:       "",
                Destination: &repository,
                Required:    true,
            },
            &cli.StringFlag{
                Name:        "release_name",
                Aliases:     []string{"n"},
                Usage:       "Define the release name to be set. You can use a string pattern to set the place where the new release tag will be set.",
                Value:       "Release version : %v",
                Destination: &releaseName,
                Required:    false,
            },
            &cli.StringFlag{
                Name:        "latest_tag",
                Aliases:     []string{"l"},
                Usage:       "The latest tag to compare against.",
                Value:       "",
                Destination: &latestTag,
                Required:    true,
            },
            &cli.StringFlag{
                Name:        "release_tag",
                Aliases:     []string{"v"},
                Usage:       "Repository name to check.",
                Value:       "HEAD",
                Destination: &releaseTag,
                Required:    true,
            },
            &cli.BoolFlag{
                Name:        "draft_release",
                Aliases:     []string{"d"},
                Usage:       "Defines if the release will be a draft or published. (default: false)",
                Destination: &draftRelease,
                Required:    false,
            },
            &cli.BoolFlag{
                Name:        "force_create",
                Aliases:     []string{"f"},
                Usage:       "Forces the creation of the release without asking for confirmation. (default: false)",
                Destination: &forceCreate,
                Required:    false,
            },
            &cli.StringSliceFlag{
                Name:        "check_pattern",
                Aliases:     []string{"p"},
                Usage:       "Define the pattern to check the files modified against.",
                Destination: &checkPattern,
                Required:    false,
            },
        },
        Action: func(c *cli.Context) error {
            commitList, err := service.GetDiffBetweenTags(authToken, repoOwner, repository, latestTag, "HEAD")
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            listOfCommitsToPrint := commitList.Commits
            checkPatternValues := checkPattern.Value()
            if len(checkPatternValues) > 0 {
                includeCommits := make(map[string]domain.Commit)

                for _, commitItem := range commitList.Commits {
                    commitDetails, err := service.GetCommitDetails(authToken, repoOwner, repository, commitItem.Sha)
                    if err != nil {
                        fmt.Println(err)
                        os.Exit(1)
                    }

                    for _, fileInCommit := range commitDetails.Files {
                        for _, pattern := range checkPatternValues {
                            if strings.Contains(fileInCommit.Filename, pattern) {
                                includeCommits[commitItem.Sha] = commitItem
                                break
                            }
                        }
                    }

                    listOfCommitsToPrint = []domain.Commit{}
                    for _, cmt := range includeCommits {
                        listOfCommitsToPrint = append(listOfCommitsToPrint, cmt)
                    }
                }
            }

            commitListPrintout, err := service.PrintCommitList(listOfCommitsToPrint, domain.CommitListReleaseTemplate)
            if err != nil {
                os.Exit(1)
            }

            releaseName = fmt.Sprintf(releaseName, releaseTag)

            fmt.Println(releaseName)
            fmt.Println(commitListPrintout)

            createRelease := false
            if !forceCreate {
                createReleasePrompt := &survey.Confirm{
                    Message: promptMessages[draftRelease],
                }

                err = survey.AskOne(createReleasePrompt, &createRelease)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
            }

            if forceCreate || createRelease {
                err = service.CreateRelease(authToken, repoOwner, repository, releaseTag, draftRelease, releaseName, commitListPrintout)
                if err != nil {
                    fmt.Printf("Failed to create release with error : %v\n", err)
                    os.Exit(1)
                }

                fmt.Printf("Created release: '%v' \n", releaseName)
            }

            return nil
        },
    }

    return &commitListCmd
}

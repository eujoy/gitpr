package createrelease

import (
    "fmt"
    "os"
    "strings"

    "github.com/AlecAivazis/survey/v2"
    "github.com/eujoy/gitpr/internal/config"
    "github.com/eujoy/gitpr/internal/domain"
    "github.com/eujoy/gitpr/internal/infra/flag"
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

    flagBuilder := flag.New(cfg)

    commitListCmd := cli.Command{
        Name:    "create-release",
        Aliases: []string{"cr"},
        Usage:   "Retrieves all the commits between two tags and creates a list of them to be used a release description..",
        Flags: flagBuilder.
            AppendAuthFlag(&authToken).
            AppendOwnerFlag(&repoOwner).
            AppendRepositoryFlag(&repository).
            AppendReleaseNameFlag(&releaseName).
            AppendLatestTagFlag(&latestTag).
            AppendReleaseTagFlag(&releaseTag).
            AppendCheckPatternFlag(&checkPattern).
            AppendDraftReleaseFlag(&draftRelease).
            AppendForceCreateFlag(&forceCreate).
            GetFlags(),
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

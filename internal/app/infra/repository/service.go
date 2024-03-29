package repository

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/eujoy/gitpr/internal/domain"
)

var commitListTerminalTemplate = `Commit List :
{{- range .}}
- Author         : {{ .Author.Username }}
  Time           : {{ .Details.Committer.Date }}
  Commit Message : {{ .Details.Message }}
{{- end}}`

var commitListReleaseTemplate = `Commits included from last release :

{{- range .}}
- [ ] (@{{ .Author.Username }}) | {{ .Details.Message }}
{{- end}}`

type resource interface {
	CreateRelease(authToken, repoOwner, repository, tagName string, draftRelease bool, name, body string) error
	GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error)
	GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error)
	GetReleaseList(authToken, repoOwner, repository string, pageSize, pageNumber int) ([]domain.Release, error)
}

// Service describes the user repositories service.
type Service struct {
	resource     resource
	templateList map[string]string
}

// NewService creates and returns a service instance.
func NewService(resource resource) *Service {
	return &Service{
		resource: resource,
		templateList: map[string]string{
			domain.CommitListTerminalTemplate: commitListTerminalTemplate,
			domain.CommitListReleaseTemplate:  commitListReleaseTemplate,
		},
	}
}

// GetCommitDetails to get the details of a commit.
func (s *Service) GetCommitDetails(authToken, repoOwner, repository, commitSha string) (domain.Commit, error) {
	commitDetails, err := s.resource.GetCommitDetails(authToken, repoOwner, repository, commitSha)
	return commitDetails, err
}

// GetDiffBetweenTags to get a list of commits.
func (s *Service) GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag string) (domain.CompareTagsResponse, error) {
	commitList, err := s.resource.GetDiffBetweenTags(authToken, repoOwner, repository, existingTag, latestTag)
	return commitList, err
}

// CreateRelease makes a post request to github api to create a new release with description.
func (s *Service) CreateRelease(authToken, repoOwner, repository, tagName string, draftRelease bool, name, body string) error {
	err := s.resource.CreateRelease(authToken, repoOwner, repository, tagName, draftRelease, name, body)
	return err
}

// GetReleaseList fetches the releases that have taken place in a repository.
func (s *Service) GetReleaseList(authToken, repoOwner, repository string, pageSize, pageNumber int) ([]domain.Release, error) {
	releaseList, err := s.resource.GetReleaseList(authToken, repoOwner, repository, pageSize, pageNumber)
	return releaseList, err
}

// PrintCommitList converts a list of commits to a checklist.
func (s *Service) PrintCommitList(commitList []domain.Commit, useTmpl string) (string, error) {
	t, err := template.New("outputTemplate").Parse(s.templateList[useTmpl])
	if err != nil {
		fmt.Printf("Failed to prepare template with error : %v\n", err)
		return "", err
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, commitList)
	if err != nil {
		fmt.Printf("Failed to print text with error : %v\n", err)
		return "", err
	}

	return tpl.String(), nil
}

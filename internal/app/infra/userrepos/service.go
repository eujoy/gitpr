package userrepos

import (
	"github.com/Angelos-Giannis/gitpr/internal/domain"
)

type resource interface {
	GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error)
}

// Service describes the user repositories service.
type Service struct {
	resource resource
}

// NewService creates and returns a service instance.
func NewService(resource resource) *Service {
	return &Service{
		resource: resource,
	}
}

// GetUserRepos retrieves all the user repositories from github.
func (s *Service) GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error) {
	userRepos, err := s.resource.GetUserRepos(authToken, pageSize, pageNumber)
	return userRepos, err
}

package actions

import "github.com/eujoy/gitpr/internal/domain"

type resource interface {
    GetWorkflowExecutions(authToken, repoOwner, repository, startDateStr, endDateStr string, pageSize, pageNumber int) ([]domain.Workflow, error)
    GetWorkflowsOfRepository(authToken, repoOwner, repository string) ([]domain.Workflow, error)
    GetWorkflowTiming(authToken, repoOwner, repository string, runID int) (domain.WorkflowTiming, error)
    GetWorkflowUsage(authToken, repoOwner, repository string, workflowID int) (domain.WorkflowTiming, error)
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

// GetWorkflowExecutions retrieves the executions of the workflows of a repository.
func (s *Service) GetWorkflowExecutions(authToken, repoOwner, repository, startDateStr, endDateStr string, pageSize, pageNumber int) ([]domain.Workflow, error) {
    workflows, err := s.resource.GetWorkflowExecutions(authToken, repoOwner, repository, startDateStr, endDateStr, pageSize, pageNumber)
    if err != nil {
        return []domain.Workflow{}, err
    }

    return workflows, nil
}

// GetWorkflowsOfRepository retrieves and returns all the workflows of a repository.
func (s *Service) GetWorkflowsOfRepository(authToken, repoOwner, repository string) ([]domain.Workflow, error) {
    workflows, err := s.resource.GetWorkflowsOfRepository(authToken, repoOwner, repository)
    if err != nil {
        return []domain.Workflow{}, err
    }

    return workflows, nil
}

// GetWorkflowTiming retrieves the timing details of a workflow.
func (s *Service) GetWorkflowTiming(authToken, repoOwner, repository string, runID int) (domain.WorkflowTiming, error) {
    workflowTiming, err := s.resource.GetWorkflowTiming(authToken, repoOwner, repository, runID)
    if err != nil {
        return domain.WorkflowTiming{}, err
    }

    return workflowTiming, nil
}

// GetWorkflowUsage retrieves the timing details of a workflow.
func (s *Service) GetWorkflowUsage(authToken, repoOwner, repository string, workflowID int) (domain.WorkflowTiming, error) {
    workflowTiming, err := s.resource.GetWorkflowUsage(authToken, repoOwner, repository, workflowID)
    if err != nil {
        return domain.WorkflowTiming{}, err
    }

    return workflowTiming, nil
}

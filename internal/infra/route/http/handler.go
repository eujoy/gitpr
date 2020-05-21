package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Angelos-Giannis/gitpr/internal/config"
	"github.com/Angelos-Giannis/gitpr/internal/domain"
)

// badRequestErrorMessage is used to prettify the bad request error message returned.
type badRequestErrorMessage struct {
	ErrorMessage string `json:"error_message"`
}

type userReposService interface {
	GetUserRepos(authToken string, pageSize int, pageNumber int) ([]domain.Repository, error)
}

type pullRequestsService interface {
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) ([]domain.PullRequest, error)
}

// Handler describes the route handler.
type Handler struct {
	cfg                 config.Config
	userReposService    userReposService
	pullRequestsService pullRequestsService
}

// NewHandler creates and returns a new http handler.
func NewHandler(cfg config.Config, userReposService userReposService, pullRequestsService pullRequestsService) *Handler {
	return &Handler{
		cfg:                 cfg,
		userReposService:    userReposService,
		pullRequestsService: pullRequestsService,
	}
}

// GetDefaultSettings returns the default values for the service.
func (h *Handler) GetDefaultSettings(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	b, err := json.Marshal(h.cfg.Defaults)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// GetUserRepos retrieves and returns all the repositories of a user.
func (h *Handler) GetUserRepos(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	var err error

	authToken := ""
	if _, ok := req.URL.Query()["authToken"]; ok {
		authToken = req.URL.Query().Get("authToken")
	} else {
		err = errors.New("missing required field authToken")
		errBadRequest(w, err.Error())
		return
	}

	pageSize := h.cfg.Defaults.PageSize
	if _, ok := req.URL.Query()["pageSize"]; ok {
		pageSize, err = strconv.Atoi(req.URL.Query().Get("pageSize"))
		if err != nil {
			errBadRequest(w, err.Error())
			return
		}
	}

	currentPage := 1
	if _, ok := req.URL.Query()["page"]; ok {
		currentPage, err = strconv.Atoi(req.URL.Query().Get("page"))
		if err != nil {
			errBadRequest(w, err.Error())
			return
		}
	}

	userRepos, err := h.userReposService.GetUserRepos(authToken, pageSize, currentPage)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	b, err := json.Marshal(userRepos)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// GetPullRequestsOfRepository retrieves and returns all the pull requests of a repository.
func (h *Handler) GetPullRequestsOfRepository(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	var err error

	authToken := ""
	if _, ok := req.URL.Query()["authToken"]; ok {
		authToken = req.URL.Query().Get("authToken")
	} else {
		err = errors.New("missing required field authToken")
		errBadRequest(w, err.Error())
		return
	}

	repoOwner := ""
	if _, ok := req.URL.Query()["repoOwner"]; ok {
		repoOwner = req.URL.Query().Get("repoOwner")
	} else {
		err = errors.New("missing required field repoOwner")
		errBadRequest(w, err.Error())
		return
	}

	repository := ""
	if _, ok := req.URL.Query()["repository"]; ok {
		repository = req.URL.Query().Get("repository")
	} else {
		err = errors.New("missing required field repository")
		errBadRequest(w, err.Error())
		return
	}

	baseBranch := h.cfg.Defaults.BaseBranch
	if _, ok := req.URL.Query()["baseBranch"]; ok {
		baseBranch = req.URL.Query().Get("baseBranch")
	}

	prState := h.cfg.Defaults.PullRequestState
	if _, ok := req.URL.Query()["prState"]; ok {
		prState = req.URL.Query().Get("prState")
	}
	prState = validatePrStateAndGetDefault(h.cfg, prState)

	pageSize := h.cfg.Defaults.PageSize
	if _, ok := req.URL.Query()["pageSize"]; ok {
		pageSize, err = strconv.Atoi(req.URL.Query().Get("pageSize"))
		if err != nil {
			errBadRequest(w, err.Error())
			return
		}
	}

	currentPage := 1
	if _, ok := req.URL.Query()["page"]; ok {
		currentPage, err = strconv.Atoi(req.URL.Query().Get("page"))
		if err != nil {
			errBadRequest(w, err.Error())
			return
		}
	}

	pullRequests, err := h.pullRequestsService.GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState, pageSize, currentPage)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	b, err := json.Marshal(pullRequests)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// validatePrStateAndGetDefault checks if the requested state of pull requests is valid and returns
// it in case it is, otherwise it returns the default pull request state.
func validatePrStateAndGetDefault(cfg config.Config, prState string) string {
	for _, prs := range cfg.Defaults.AllowedPullRequestStates {
		if prState == prs {
			return prState
		}
	}

	return cfg.Defaults.PullRequestState
}

// errBadRequest generates and error request with status "Bad Request" including the error message as well.
func errBadRequest(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	badRequestError := badRequestErrorMessage{ErrorMessage: errMsg}
	b, err := json.Marshal(badRequestError)
	if err != nil {
		return
	}

	w.Write(b)
}

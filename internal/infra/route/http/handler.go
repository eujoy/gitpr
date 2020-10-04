package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/domain"
)

const (
	returnContentType = "application/json"
)

// badRequestErrorMessage is used to prettify the bad request error message returned.
type badRequestErrorMessage struct {
	ErrorMessage string `json:"error"`
}

type userReposService interface {
	GetUserRepos(authToken string, pageSize int, pageNumber int) (domain.UserReposResponse, error)
}

type pullRequestsService interface {
	GetPullRequestsOfRepository(authToken, repoOwner, repository, baseBranch, prState string, pageSize int, pageNumber int) (domain.RepoPullRequestsResponse, error)
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

// GetSettings returns the default values for the service.
func (h *Handler) GetSettings(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if req.Method != http.MethodGet {
		w.Header().Set("Content-Type", returnContentType)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	b, err := json.Marshal(h.cfg.Settings)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", returnContentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// GetUserRepos retrieves and returns all the repositories of a user.
func (h *Handler) GetUserRepos(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if req.Method != http.MethodGet {
		w.Header().Set("Content-Type", returnContentType)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	var err error

	authToken, err := parseRequiredStringFromRequestWithDefault(req, "authToken", h.cfg.Clients.Github.Token.DefaultValue, true, true)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	pageSize := h.cfg.Settings.PageSize
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

	w.Header().Set("Content-Type", returnContentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// GetPullRequestsOfRepository retrieves and returns all the pull requests of a repository.
func (h *Handler) GetPullRequestsOfRepository(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if req.Method != http.MethodGet {
		w.Header().Set("Content-Type", returnContentType)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	var err error

	authToken, err := parseRequiredStringFromRequestWithDefault(req, "authToken", h.cfg.Clients.Github.Token.DefaultValue, true, true)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	repoOwner, err := parseRequiredStringFromRequestWithDefault(req, "repoOwner", "", true, false)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	repository, err := parseRequiredStringFromRequestWithDefault(req, "repository", "", true, false)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	baseBranch, err := parseRequiredStringFromRequestWithDefault(req, "baseBranch", h.cfg.Settings.BaseBranch, false, true)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	prState, err := parseRequiredStringFromRequestWithDefault(req, "prState", h.cfg.Settings.PullRequestState, false, true)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	prState = validatePrStateAndGetDefault(h.cfg, prState)

	pageSize, err := parseIntFieldsFromRequest(req, "pageSize", h.cfg.Settings.PageSize, false, true)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}

	currentPage, err := parseIntFieldsFromRequest(req, "page", 1, false, true)
	if err != nil {
		errBadRequest(w, err.Error())
		return
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

	w.Header().Set("Content-Type", returnContentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// parseRequiredStringFromRequestWithDefault parses the request fields that are of string type and returns
// the respective value provided.
func parseRequiredStringFromRequestWithDefault(req *http.Request, urlParam, defaultValue string, isRequired, hasDefault bool) (requestValue string, err error) {
	if hasDefault {
		requestValue = defaultValue
	}

	if _, ok := req.URL.Query()[urlParam]; ok {
		requestValue = req.URL.Query().Get(urlParam)
	} else {
		if isRequired && requestValue == "" {
			err = fmt.Errorf("missing required field %v", urlParam)
			return "", err
		}
	}

	return requestValue, nil
}

// parseIntFieldsFromRequest parses and returns an integer field from the request.
func parseIntFieldsFromRequest(req *http.Request, urlParam string, defaultValue int, isRequired, hasDefault bool) (requestValue int, err error) {
	if hasDefault {
		requestValue = defaultValue
	}

	if _, ok := req.URL.Query()[urlParam]; ok {
		requestValue, err = strconv.Atoi(req.URL.Query().Get(urlParam))
		if err != nil {
			return 0, err
		}
	} else {
		if isRequired {
			err := fmt.Errorf("missing required field %v", urlParam)
			return 0, err
		}
	}

	return requestValue, err
}

// validatePrStateAndGetDefault checks if the requested state of pull requests is valid and returns
// it in case it is, otherwise it returns the default pull request state.
func validatePrStateAndGetDefault(cfg config.Config, prState string) string {
	for _, prs := range cfg.Settings.AllowedPullRequestStates {
		if prState == prs {
			return prState
		}
	}

	return cfg.Settings.PullRequestState
}

// errBadRequest generates and error request with status "Bad Request" including the error message as well.
func errBadRequest(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", returnContentType)
	w.WriteHeader(http.StatusBadRequest)

	badRequestError := badRequestErrorMessage{ErrorMessage: errMsg}
	b, err := json.Marshal(badRequestError)
	if err != nil {
		return
	}

	_, _ = w.Write(b)
}

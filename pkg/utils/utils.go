package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Angelos-Giannis/gitpr/internal/config"
)

// Utils describes the common utilities package.
type Utils struct {
	cfg config.Config
}

// NewUtils create and return a new utilities struct.
func NewUtils(cfg config.Config) *Utils {
	return &Utils{
		cfg: cfg,
	}
}

// ClearTerminalScreen clears up the screen to get the new data in.
func (u *Utils) ClearTerminalScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to clear screen!")
		os.Exit(1)
	}
}

// GetPageOptions prepares and returns a list of available options for the user repos list.
func (u *Utils) GetPageOptions(respLength int, pageSize int, currentPage int) []string {
	var options []string

	if respLength == pageSize {
		options = append(options, u.cfg.Pagination.Next)
	}
	if currentPage > 1 {
		options = append(options, u.cfg.Pagination.Previous)
	}
	options = append(options, u.cfg.Pagination.Exit)

	return options
}

// GetNextPageNumberOrExit reads user input and returns if the process shall continue or not.
func (u *Utils) GetNextPageNumberOrExit(surveySelection string, currentPage int) (int, bool) {
	switch surveySelection {
	case u.cfg.Pagination.Next:
		currentPage++
		return currentPage, true
	case u.cfg.Pagination.Previous:
		currentPage--
		if currentPage < 0 {
			currentPage = 0
		}
		return currentPage, true
	case u.cfg.Pagination.Exit:
		return 0, false
	}

	return 0, false
}

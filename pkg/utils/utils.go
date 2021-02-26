package utils

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/eujoy/gitpr/internal/config"
)

// Utils describes the common utilities package.
type Utils struct {
	cfg config.Config
}

// New create and return a new utilities struct.
func New(cfg config.Config) *Utils {
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

func (u *Utils) ConvertDurationToString(dur time.Duration) string {
	if dur == time.Duration(0) {
		return ""
	}

	days := int64(dur.Hours() / 24)
	hours := int64(math.Mod(dur.Hours(), 24))
	minutes := int64(math.Mod(dur.Minutes(), 60))
	seconds := int64(math.Mod(dur.Seconds(), 60))

	formattedDuration := fmt.Sprintf("%d days & %02d:%02d:%02d", days, hours, minutes, seconds)

	return formattedDuration
}

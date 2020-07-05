package utils_test

import (
	"reflect"
	"testing"

	"github.com/Angelos-Giannis/gitpr/internal/config"
	"github.com/Angelos-Giannis/gitpr/pkg/utils"
)

func TestNew(t *testing.T) {
	var cfg config.Config
	actualUtils := utils.New(cfg)

	var expectedUtils *utils.Utils
	if reflect.TypeOf(actualUtils) != reflect.TypeOf(expectedUtils) {
		t.Errorf("Expected to get type as '%v' but got '%v'", reflect.TypeOf(expectedUtils), reflect.TypeOf(actualUtils))
	}
}

// @todo Add test for GetPageOptions

func TestGetNextPageNumberOrExit(t *testing.T) {
	var cfg config.Config
	cfg.Pagination.Next = "Next"
	cfg.Pagination.Previous = "Previous"
	cfg.Pagination.Exit = "Exit"

	utils := utils.New(cfg)

	type input struct {
		surveySelection string
		currentPage     int
	}

	type expected struct {
		newPage        int
		shouldContinue bool
	}

	testCases := map[string]struct {
		input    input
		expected expected
	}{
		"Normal execution with getting next page as 1 - selecting 'Next'": {
			input{"Next", 0},
			expected{1, true},
		},
		"Normal execution with getting next page as 1 - selecting 'Previous'": {
			input{"Previous", 2},
			expected{1, true},
		},
		"Trying to get to previous page while being on first page - expecting first page as next": {
			input{"Previous", 0},
			expected{0, true},
		},
		"Selecting to exit - expecting first page as next and false as should continue": {
			input{"Exit", 0},
			expected{0, false},
		},
		"Providing random selection option - expecting first page as next and false as should continue": {
			input{"Some Random Option", 0},
			expected{0, false},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualNewPage, actualShouldContinue := utils.GetNextPageNumberOrExit(tc.input.surveySelection, tc.input.currentPage)

			if !reflect.DeepEqual(tc.expected.newPage, actualNewPage) {
				t.Errorf("Expected to get '%v' as new page, but got '%v'", tc.expected.newPage, actualNewPage)
			}

			if !reflect.DeepEqual(tc.expected.shouldContinue, actualShouldContinue) {
				t.Errorf("Expected to get '%v' as should continue, but got '%v'", tc.expected.shouldContinue, actualShouldContinue)
			}
		})
	}
}

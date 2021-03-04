package utils_test

import (
	"reflect"
	"testing"

	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/pkg/utils"
)

func TestNew(t *testing.T) {
	var cfg config.Config
	actualUtils := utils.New(cfg)

	var expectedUtils *utils.Utils
	if reflect.TypeOf(actualUtils) != reflect.TypeOf(expectedUtils) {
		t.Errorf("Expected to get type as '%v' but got '%v'", reflect.TypeOf(expectedUtils), reflect.TypeOf(actualUtils))
	}
}

func TestGetPageOptions(t *testing.T) {
	var cfg config.Config
	cfg.Pagination.Next = "Next"
	cfg.Pagination.Previous = "Previous"
	cfg.Pagination.Exit = "Exit"

	utilities := utils.New(cfg)

	type input struct {
		respLength  int
		pageSize    int
		currentPage int
	}

	type expected struct {
		pageOptions []string
	}

	testCases := map[string]struct {
		input    input
		expected expected
	}{
		"Normal execution being on the first page": {
			input{5, 5, 0},
			expected{[]string{"Next", "Exit"}},
		},
		"Normal execution being on the last page": {
			input{4, 5, 2},
			expected{[]string{"Previous", "Exit"}},
		},
		"Normal execution being on a mid page": {
			input{3, 3, 2},
			expected{[]string{"Next", "Previous", "Exit"}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualPageOptions := utilities.GetPageOptions(tc.input.respLength, tc.input.pageSize, tc.input.currentPage)

			if !reflect.DeepEqual(tc.expected.pageOptions, actualPageOptions) {
				t.Errorf("Expected to get '%v' as page options, but got '%v'", tc.expected.pageOptions, actualPageOptions)
			}
		})
	}
}

func TestGetNextPageNumberOrExit(t *testing.T) {
	var cfg config.Config
	cfg.Pagination.Next = "Next"
	cfg.Pagination.Previous = "Previous"
	cfg.Pagination.Exit = "Exit"

	utilities := utils.New(cfg)

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
			actualNewPage, actualShouldContinue := utilities.GetNextPageNumberOrExit(tc.input.surveySelection, tc.input.currentPage)

			if !reflect.DeepEqual(tc.expected.newPage, actualNewPage) {
				t.Errorf("Expected to get '%v' as new page, but got '%v'", tc.expected.newPage, actualNewPage)
			}

			if !reflect.DeepEqual(tc.expected.shouldContinue, actualShouldContinue) {
				t.Errorf("Expected to get '%v' as should continue, but got '%v'", tc.expected.shouldContinue, actualShouldContinue)
			}
		})
	}
}

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

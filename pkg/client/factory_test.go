package client_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/pkg/client"
	"github.com/eujoy/gitpr/pkg/client/github"
)

func TestNewFactory(t *testing.T) {
	// client := mock.Client{}
	// m := &mock.App{}

	var cfg config.Config

	t.Run("Normal instantiation using github client", func(t *testing.T) {
		actualFactory, actualError := client.NewFactory("github", cfg)

		if reflect.TypeOf(actualFactory) != reflect.TypeOf(&client.Factory{}) {
			t.Errorf("Expected to get factory type '%v', but got '%v'", reflect.TypeOf(&client.Factory{}), reflect.TypeOf(actualFactory))
		}
		if actualError != nil {
			t.Errorf("Expected to get nil as error, but got '%v'", actualError)
		}
	})

	t.Run("Try to instatiate a client using invalid client type - expecting an error", func(t *testing.T) {
		expectedError := errors.New("failed to initialize client")
		actualFactory, actualError := client.NewFactory("invalid client", cfg)

		if actualFactory != nil {
			t.Errorf("Expected to get nil as factory, but got '%v'", reflect.TypeOf(actualFactory))
		}
		if !reflect.DeepEqual(actualError, expectedError) {
			t.Errorf("Expected to get '%v' as error, but got '%v'", expectedError, actualError)
		}
	})

	t.Run("Try to instatiate a client using an empty client type - expecting an error", func(t *testing.T) {
		expectedError := errors.New("failed to initialize client")
		actualFactory, actualError := client.NewFactory("invalid client", cfg)

		if actualFactory != nil {
			t.Errorf("Expected to get nil as factory, but got '%v'", reflect.TypeOf(actualFactory))
		}
		if !reflect.DeepEqual(actualError, expectedError) {
			t.Errorf("Expected to get '%v' as error, but got '%v'", expectedError, actualError)
		}
	})
}

func TestGetClient(t *testing.T) {
	var cfg config.Config
	actualFactory, actualError := client.NewFactory("github", cfg)

	if reflect.TypeOf(actualFactory) != reflect.TypeOf(&client.Factory{}) {
		t.Errorf("Expected to get factory type '%v', but got '%v'", reflect.TypeOf(&client.Factory{}), reflect.TypeOf(actualFactory))
	}
	if actualError != nil {
		t.Errorf("Expected to get nil as error, but got '%v'", actualError)
	}

	actualClient := actualFactory.GetClient()

	if reflect.TypeOf(actualClient) != reflect.TypeOf(&github.Resource{}) {
		t.Errorf("Expected to get github client resource type '%v', but got '%v'", reflect.TypeOf(&github.Resource{}), reflect.TypeOf(actualClient))
	}
}

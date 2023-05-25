package app_test

import (
	"fmt"
	"testing"

	"github.com/bongnv/sen/app"
)

type mockComponent struct {
	Data int `inject:"data"`
}

func TestInjector(t *testing.T) {
	t.Run("should return error if there is no registered dependency", func(t *testing.T) {
		injector := app.NewInjector()
		err := injector.Inject(&mockComponent{})
		if fmt.Sprintf("%v", err) != "injector: data is not registered" {
			t.Errorf("Unexpected err %v", err)
		}
	})
}

func TestInjector_Retrieve(t *testing.T) {
	t.Run("should return the component if it's registered", func(t *testing.T) {
		injector := app.NewInjector()

		err := injector.Register("data", 10)
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}

		data, err := injector.Retrieve("data")
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}

		if 10 != data {
			t.Errorf("Unexpected data: %v", data)
		}
	})

	t.Run("should return an error if the component isn't registered", func(t *testing.T) {
		injector := app.NewInjector()
		data, err := injector.Retrieve("data")
		if err != app.ErrComponentNotRegistered {
			t.Errorf("Unexpected err %v", err)
		}

		if data != nil {
			t.Errorf("Expected nil but got %v", data)
		}
	})
}

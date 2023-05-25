package plugin_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bongnv/sen/app"
	"github.com/bongnv/sen/plugin"
)

type mockFactory struct {
	Data int `inject:"data"`
	err  error
}

func (f mockFactory) New() (int, error) {
	return f.Data, f.err
}

func TestFactory(t *testing.T) {
	t.Run("should return an error when couldn't initialise a new component", func(t *testing.T) {
		p := plugin.Provider[int]("mock-component", &mockFactory{
			err: errors.New("random error"),
		})

		mockApp := app.New()
		err := mockApp.With(
			plugin.Component("data", 10),
			p,
		)
		if fmt.Sprintf("%v", err) != "random error" {
			t.Errorf("Unexpected err %v", err)
		}
	})

	t.Run("should register successfully when there is no error", func(t *testing.T) {
		p := plugin.Provider[int]("mock-component", &mockFactory{})

		mockApp := app.New()
		err := mockApp.With(
			plugin.Component("data", 10),
			p,
		)
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}

		component, errGet := mockApp.Retrieve("mock-component")
		if errGet != nil {
			t.Errorf("Unexpected err %v", err)
		}
		if 10 != component {
			t.Errorf("Expected 10 by got %v", component)
		}
	})

	t.Run("should propergate error when couldn't populate depdencies", func(t *testing.T) {
		p := plugin.Provider[int]("mock-component", &mockFactory{})

		mockApp := app.New()
		err := mockApp.With(
			p,
		)
		if fmt.Sprintf("%v", err) != dataInjectErrMsg {
			t.Errorf("Unexpected err %v", err)
		}
	})
}

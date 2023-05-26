package sen_test

import (
	"fmt"
	"testing"

	"github.com/bongnv/sen"
)

const dataInjectErrMsg = "injector: data is not registered"

type mockPlugin struct {
	Data int `inject:"data"`
}

func (mockPlugin) Initialize() error {
	return nil
}

func TestComponent(t *testing.T) {
	t.Run("should register a component into the application", func(t *testing.T) {
		component := &mockComponent{}
		app := sen.New()
		err := app.With(
			sen.Component("data", 10),
			sen.Component("need-data", component),
		)
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}
		if component.Data != 10 {
			t.Errorf("Unexpected data %v", component.Data)
		}
	})

	t.Run("should propergate error if a component cannot be registered", func(t *testing.T) {
		component := &mockComponent{}
		app := sen.New()
		err := app.With(sen.Component("need-data", component))
		if fmt.Sprintf("%v", err) != dataInjectErrMsg {
			t.Errorf("Unexpected err %v", err)
		}
	})
}

func TestModule(t *testing.T) {
	t.Run("should apply all plugins into the application", func(t *testing.T) {
		component := &mockComponent{}
		m := sen.Module(
			sen.Component("data", 10),
			sen.Component("need-data", component),
		)
		app := sen.New()
		err := app.With(m)
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}
		if component.Data != 10 {
			t.Errorf("Unexpected data %v", component.Data)
		}
	})

	t.Run("should propagate error if one plugin returns error", func(t *testing.T) {
		m := sen.Module(
			sen.Component("error-plugin", &mockComponent{}),
			sen.Component("ok-plugin", 10),
		)
		app := sen.New()
		err := app.With(m)
		if fmt.Sprintf("%v", err) != dataInjectErrMsg {
			t.Errorf("Unexpected err %v", err)
		}
	})

	t.Run("should propagate error if one plugin doesn't have enough dependencies", func(t *testing.T) {
		m := sen.Module(
			&mockPlugin{},
		)

		app := sen.New()
		err := app.With(m)
		if fmt.Sprintf("%v", err) != dataInjectErrMsg {
			t.Errorf("Unexpected err %v", err)
		}
	})
}

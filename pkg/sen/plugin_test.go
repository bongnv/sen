package sen_test

import (
	"fmt"
	"testing"

	"github.com/bongnv/sen/pkg/sen"
)

const dataInjectErrMsg = "hub: data is not registered"

type mockPlugin struct {
	Data int `inject:"data"`
}

func (mockPlugin) Initialize() error {
	return nil
}

func TestComponent(t *testing.T) {
	t.Run("should register a component into the application", func(t *testing.T) {
		component := &mockComponent{}
		_, err := sen.New(
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
		_, err := sen.New(sen.Component("need-data", component))
		if fmt.Sprintf("%v", err) != dataInjectErrMsg {
			t.Errorf("Unexpected err %v", err)
		}
	})
}

func TestBundle(t *testing.T) {
	t.Run("should apply all plugins into the application", func(t *testing.T) {
		component := &mockComponent{}
		m := sen.Bundle(
			sen.Component("data", 10),
			sen.Component("need-data", component),
		)
		_, err := sen.New(m)
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}
		if component.Data != 10 {
			t.Errorf("Unexpected data %v", component.Data)
		}
	})

	t.Run("should propagate error if one plugin returns error", func(t *testing.T) {
		m := sen.Bundle(
			sen.Component("error-plugin", &mockComponent{}),
			sen.Component("ok-plugin", 10),
		)
		_, err := sen.New(m)
		if fmt.Sprintf("%v", err) != dataInjectErrMsg {
			t.Errorf("Unexpected err %v", err)
		}
	})

	t.Run("should propagate error if one plugin doesn't have enough dependencies", func(t *testing.T) {
		m := sen.Bundle(
			&mockPlugin{},
		)

		_, err := sen.New(m)
		if fmt.Sprintf("%v", err) != dataInjectErrMsg {
			t.Errorf("Unexpected err %v", err)
		}
	})
}

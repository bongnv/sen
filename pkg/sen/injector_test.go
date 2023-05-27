package sen_test

import (
	"fmt"
	"testing"

	"github.com/bongnv/sen/pkg/sen"
)

type mockComponent struct {
	Data int `inject:"data"`
}

type mockOptionalComponent struct {
	Data int `inject:"data,optional"`
}

type mockErrorComponent struct {
	Data int `inject:"data,required"`
}

func TestInjector(t *testing.T) {
	t.Run("should return error if there is no registered dependency", func(t *testing.T) {
		injector := sen.NewInjector()
		err := injector.Inject(&mockComponent{})
		if fmt.Sprintf("%v", err) != "injector: data is not registered" {
			t.Errorf("Unexpected err %v", err)
		}
	})

	t.Run("should not return if it's marked as optional", func(t *testing.T) {
		injector := sen.NewInjector()
		err := injector.Inject(&mockOptionalComponent{})
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
	})

	t.Run("should return an error for an unsupported option", func(t *testing.T) {
		injector := sen.NewInjector()
		err := injector.Inject(&mockErrorComponent{})
		if fmt.Sprintf("%v", err) != "injector: required is unexpected" {
			t.Errorf("Unexpected error %v", err)
		}
	})

	t.Run("should return an error if it isn't injectable", func(t *testing.T) {
		injector := sen.NewInjector()
		err := injector.Inject(mockErrorComponent{})
		if fmt.Sprintf("%v", err) != "injector: sen_test.mockErrorComponent is not injectable, a pointer is expected" {
			t.Errorf("Unexpected error %v", err)
		}
	})
}

func TestInjector_Retrieve(t *testing.T) {
	t.Run("should return the component if it's registered", func(t *testing.T) {
		injector := sen.NewInjector()

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
		injector := sen.NewInjector()
		data, err := injector.Retrieve("data")
		if err != sen.ErrComponentNotRegistered {
			t.Errorf("Unexpected err %v", err)
		}

		if data != nil {
			t.Errorf("Expected nil but got %v", data)
		}
	})
}

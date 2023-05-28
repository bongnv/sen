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

func TestHub(t *testing.T) {
	t.Run("should return error if there is no registered dependency", func(t *testing.T) {
		hub := sen.NewHub()
		err := hub.Inject(&mockComponent{})
		if fmt.Sprintf("%v", err) != "hub: data is not registered" {
			t.Errorf("Unexpected err %v", err)
		}
	})

	t.Run("should not return if it's marked as optional", func(t *testing.T) {
		hub := sen.NewHub()
		err := hub.Inject(&mockOptionalComponent{})
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
	})

	t.Run("should return an error for an unsupported option", func(t *testing.T) {
		hub := sen.NewHub()
		err := hub.Inject(&mockErrorComponent{})
		if fmt.Sprintf("%v", err) != "hub: required is unexpected" {
			t.Errorf("Unexpected error %v", err)
		}
	})

	t.Run("should return an error if it isn't injectable", func(t *testing.T) {
		hub := sen.NewHub()
		err := hub.Inject(mockErrorComponent{})
		if fmt.Sprintf("%v", err) != "hub: sen_test.mockErrorComponent is not injectable, a pointer is expected" {
			t.Errorf("Unexpected error %v", err)
		}
	})
}

func TestHub_Retrieve(t *testing.T) {
	t.Run("should return the component if it's registered", func(t *testing.T) {
		hub := sen.NewHub()

		err := hub.Register("data", 10)
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}

		data, err := hub.Retrieve("data")
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}

		if 10 != data {
			t.Errorf("Unexpected data: %v", data)
		}
	})

	t.Run("should return an error if the component isn't registered", func(t *testing.T) {
		hub := sen.NewHub()
		data, err := hub.Retrieve("data")
		if err != sen.ErrComponentNotRegistered {
			t.Errorf("Unexpected err %v", err)
		}

		if data != nil {
			t.Errorf("Expected nil but got %v", data)
		}
	})
}

package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/sen/app"
)

type mockComponent struct {
	Data int `inject:"data"`
}

func TestInjector(t *testing.T) {
	t.Run("should return error if there is no registered dependency", func(t *testing.T) {
		injector := app.NewInjector()
		err := injector.Inject(&mockComponent{})
		require.EqualError(t, err, "injector: data is not registered")
	})
}

package app_test

import (
	"testing"

	"github.com/bongnv/sen/app"
	"github.com/stretchr/testify/require"
)

type mockComponent struct {
	Data int `inject:"data"`
}

func TestInjector(t *testing.T) {
	injector := app.NewInjector()

	t.Run("should return error if there is no registered dependency", func(t *testing.T) {
		err := injector.Inject(&mockComponent{})
		require.EqualError(t, err, "injector: data is not registered")
	})
}

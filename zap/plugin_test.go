package zap_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	uberZap "go.uber.org/zap"

	"github.com/bongnv/sen"
	"github.com/bongnv/sen/zap"
)

func TestPlugin(t *testing.T) {
	t.Run("should inject logger to the app", func(t *testing.T) {
		app := sen.New()
		err := app.Apply(zap.Plugin())
		require.NoError(t, err)

		logger, err := app.Retrieve("logger")
		require.NoError(t, err)
		require.IsType(t, &uberZap.Logger{}, logger)
	})
}

package grpc_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/sen"
	"github.com/bongnv/sen/grpc"
)

type configRetriever struct {
	Config *grpc.Config `inject:"grpc.config"`
}

func getConfigFrom(app *sen.Application) (*grpc.Config, error) {
	r := &configRetriever{}
	if err := app.Inject(r); err != nil {
		return nil, err
	}

	return r.Config, nil
}

func TestEnvConfigPlugin(t *testing.T) {
	t.Run("should be able to load address from GRPC_ADDR", func(t *testing.T) {
		require.NoError(t, os.Setenv("GRPC_ADDR", ":8080"))
		defer os.Unsetenv("GRPC_ADDR")

		app := sen.New()
		require.NoError(t, app.Apply(grpc.EnvConfigPlugin()))
		cfg, err := getConfigFrom(app)
		require.NoError(t, err)
		require.Equal(t, ":8080", cfg.Address)
	})
}

func TestPlugin(t *testing.T) {
	t.Run("should be able to start the server", func(t *testing.T) {
		require.NoError(t, os.Setenv("GRPC_ADDR", ":8080"))
		defer os.Unsetenv("GRPC_ADDR")

		app := sen.New()
		require.NoError(t, app.Apply(
			grpc.EnvConfigPlugin(),
			grpc.Plugin(),
		))

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		shutdownDoneCh := make(chan struct{})
		runDoneCh := make(chan struct{})

		go func() {
			require.NoError(t, app.Run(ctx))
			close(runDoneCh)
		}()

		go func() {
			require.NoError(t, app.Shutdown(ctx))
			<-runDoneCh
			close(shutdownDoneCh)
		}()

		select {
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "test timed out")
		case <-shutdownDoneCh:
		}
	})
}

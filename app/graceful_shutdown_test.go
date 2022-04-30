package app_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/sen/app"
)

type mockWaitingPlugin struct {
	App   *app.Application `inject:"app"`
	ready chan struct{}
}

func (p *mockWaitingPlugin) Apply(ctx context.Context) error {
	waitCh := make(chan error)
	p.App.OnRun(func(ctx context.Context) error {
		close(p.ready)
		return <-waitCh
	})

	p.App.OnShutdown(func(ctx context.Context) error {
		close(waitCh)
		return nil
	})

	return nil
}

func makeMockWaitingPlugin() *mockWaitingPlugin {
	return &mockWaitingPlugin{
		ready: make(chan struct{}),
	}
}

func TestGracefulShutdown(t *testing.T) {
	t.Run("should exit the app if no background job", func(t *testing.T) {
		app := app.New(app.GracefulShutdown())
		doneCh := make(chan struct{})
		go func() {
			err := app.Run(context.Background())
			require.NoError(t, err)
			close(doneCh)
		}()

		select {
		case <-doneCh:
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "test timed out")
		}
	})

	t.Run("should exit the app if shutdown is called", func(t *testing.T) {
		app := app.New(
			app.GracefulShutdown(),
			makeMockWaitingPlugin(),
		)

		doneCh := make(chan struct{})
		go func() {
			err := app.Run(context.Background())
			require.NoError(t, err)
			close(doneCh)
		}()

		require.NoError(t, app.Shutdown(context.Background()))
		select {
		case <-doneCh:
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "test timed out")
		}
	})

	t.Run("should exit the app if SIGTERM signal is received", func(t *testing.T) {
		m := makeMockWaitingPlugin()
		app := app.New(
			app.GracefulShutdown(),
			m,
		)

		doneCh := make(chan struct{})
		go func() {
			err := app.Run(context.Background())
			require.NoError(t, err)
			close(doneCh)
		}()

		<-m.ready
		require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))
		select {
		case <-doneCh:
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "test timed out")
		}
	})
}

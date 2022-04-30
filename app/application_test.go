package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/sen/app"
)

func TestApplication(t *testing.T) {
	t.Run("should run all hooks for OnRun stage", func(t *testing.T) {
		hook1Called := 0
		hook2Called := 0

		lc := app.New()
		lc.OnRun(func(_ context.Context) error {
			hook1Called++
			return nil
		})

		lc.OnRun(func(ctx context.Context) error {
			hook2Called++
			return nil
		})

		require.NoError(t, lc.Run(context.Background()))
		require.Equal(t, 1, hook1Called)
		require.Equal(t, 1, hook2Called)
	})

	t.Run("should run all hooks for OnShutdown stage", func(t *testing.T) {
		hook1Called := 0
		hook2Called := 0
		doneCh := make(chan struct{})

		lc := app.New()
		lc.OnShutdown(func(_ context.Context) error {
			hook1Called++
			return nil
		})

		lc.OnShutdown(func(_ context.Context) error {
			hook2Called++
			return nil
		})

		go func() {
			require.NoError(t, lc.Run(context.Background()))
			close(doneCh)
		}()

		require.NoError(t, lc.Shutdown(context.Background()))
		select {
		case <-doneCh:
			require.Equal(t, 1, hook1Called)
			require.Equal(t, 1, hook2Called)
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "test timed out")
		}
	})

	t.Run("should propergate the error if a hook returns an error", func(t *testing.T) {
		hook1Called := 0
		doneCh := make(chan struct{})

		lc := app.New()
		lc.OnRun(func(_ context.Context) error {
			hook1Called++
			close(doneCh)
			return nil
		})

		lc.OnRun(func(_ context.Context) error {
			return errors.New("random error")
		})

		lc.OnRun(func(_ context.Context) error {
			return errors.New("random error")
		})

		require.EqualError(t, lc.Run(context.Background()), "random error")
		select {
		case <-doneCh:
			require.Equal(t, 1, hook1Called)
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "test timed out")
		}
	})

	t.Run("should propergate the error from OnShutdown hook", func(t *testing.T) {
		hook1Called := 0
		doneCh := make(chan struct{})

		lc := app.New()
		lc.OnShutdown(func(_ context.Context) error {
			hook1Called++
			return errors.New("random error")
		})

		go func() {
			require.EqualError(t, lc.Run(context.Background()), "random error")
			close(doneCh)
		}()

		require.EqualError(t, lc.Shutdown(context.Background()), "random error")

		select {
		case <-doneCh:
			require.Equal(t, 1, hook1Called)
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "test timed out")
		}
	})
}

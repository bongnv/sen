package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/sen/app"
)

func TestLifecycle(t *testing.T) {
	t.Run("should run all hooks for OnRun stage", func(t *testing.T) {
		hook1Called := 0
		hook2Called := 0

		lc := &app.DefaultLifecycle{}
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

		lc := &app.DefaultLifecycle{}
		lc.OnShutdown(func(_ context.Context) error {
			hook1Called++
			return nil
		})

		lc.OnShutdown(func(_ context.Context) error {
			hook2Called++
			return nil
		})

		require.NoError(t, lc.Shutdown(context.Background()))
		require.Equal(t, 1, hook1Called)
		require.Equal(t, 1, hook2Called)
	})

	t.Run("should propergate the error if a hook returns an error", func(t *testing.T) {
		hook1Called := 0

		lc := &app.DefaultLifecycle{}
		lc.OnShutdown(func(_ context.Context) error {
			hook1Called++
			return nil
		})

		lc.OnShutdown(func(_ context.Context) error {
			return errors.New("random error")
		})

		lc.OnShutdown(func(_ context.Context) error {
			return errors.New("random error")
		})

		require.EqualError(t, lc.Shutdown(context.Background()), "random error")
		require.Equal(t, 1, hook1Called)
	})

	t.Run("should return an error if Shutdown is called twice", func(t *testing.T) {
		hook1Called := 0

		lc := &app.DefaultLifecycle{}
		lc.OnShutdown(func(_ context.Context) error {
			hook1Called++
			return nil
		})

		require.NoError(t, lc.Shutdown(context.Background()))
		require.Equal(t, 1, hook1Called)
		require.EqualError(t, lc.Shutdown(context.Background()), "app: Shutdown has been called")
	})
}

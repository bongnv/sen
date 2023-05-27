package sen_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bongnv/sen/pkg/sen"
)

func TestLifecycle(t *testing.T) {
	t.Run("should run all hooks for OnRun stage", func(t *testing.T) {
		hook1Called := 0
		hook2Called := 0

		lc := sen.NewLifecycle()
		lc.OnRun(func(_ context.Context) error {
			hook1Called++
			return nil
		})

		lc.OnRun(func(ctx context.Context) error {
			hook2Called++
			return nil
		})

		err := lc.Run(context.Background())
		if err != nil {
			t.Errorf("Expected no error")
		}

		if hook1Called != 1 {
			t.Errorf("Expected hook1 is called but got %d", hook1Called)
		}

		if hook2Called != 1 {
			t.Errorf("Expected hook2 is called once but got %d", hook2Called)
		}
	})

	t.Run("should run all hooks for OnShutdown stage", func(t *testing.T) {
		hook1Called := 0
		hook2Called := 0
		doneCh := make(chan struct{})

		lc := sen.NewLifecycle()
		lc.OnShutdown(func(_ context.Context) error {
			hook1Called++
			return nil
		})

		lc.OnShutdown(func(_ context.Context) error {
			hook2Called++
			return nil
		})

		go func() {
			err := lc.Run(context.Background())
			if err != nil {
				t.Errorf("Expected no error")
			}
			close(doneCh)
		}()

		err := lc.Shutdown(context.Background())
		if err != nil {
			t.Errorf("Expected no error")
		}
		select {
		case <-doneCh:
			if hook1Called != 1 {
				t.Errorf("Expected hook1 is called but got %d", hook1Called)
			}

			if hook2Called != 1 {
				t.Errorf("Expected hook2 is called once but got %d", hook2Called)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("test timed out")
		}
	})

	t.Run("should propergate the error if a hook returns an error", func(t *testing.T) {
		hook1Called := 0
		doneCh := make(chan struct{})

		lc := sen.NewLifecycle()
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

		err := lc.Run(context.Background())
		if fmt.Sprintf("%v", err) != "random error" {
			t.Errorf("Unexpected err %v", err)
		}
		select {
		case <-doneCh:
			if hook1Called != 1 {
				t.Errorf("Expected hook1 is called but got %d", hook1Called)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("test timed out")
		}
	})

	t.Run("should call Shutdown if a hook from OnRun returns an error", func(t *testing.T) {
		hook1Called := 0
		hook2Called := 0

		lc := sen.NewLifecycle()
		lc.OnRun(func(_ context.Context) error {
			hook1Called++
			return errors.New("run error")
		})

		lc.OnShutdown(func(ctx context.Context) error {
			hook2Called++
			return nil
		})

		err := lc.Run(context.Background())
		if fmt.Sprintf("%v", err) != "run error" {
			t.Errorf("Unexpected error: %v", err)
		}

		if hook1Called != 1 {
			t.Errorf("Expected hook1 is called but got %d", hook1Called)
		}

		if hook2Called != 1 {
			t.Errorf("Expected hook2 is called once but got %d", hook2Called)
		}
	})

	t.Run("should run PostRun hook", func(t *testing.T) {
		hook1Called := 0
		doneCh := make(chan struct{})

		lc := sen.NewLifecycle()
		lc.PostRun(func(_ context.Context) error {
			hook1Called++
			return nil
		})

		go func() {
			err := lc.Run(context.Background())
			if err != nil {
				t.Errorf("Expected no error")
			}
			close(doneCh)
		}()

		err := lc.Shutdown(context.Background())
		if err != nil {
			t.Errorf("Expected no error")
		}

		select {
		case <-doneCh:
			if hook1Called != 1 {
				t.Errorf("Expected hook1 is called but got %d", hook1Called)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("test timed out")
		}
	})
}

package echo_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bongnv/sen/pkg/sen"
	"github.com/labstack/echo/v4"

	echoPlugin "github.com/bongnv/sen/pkg/plugins/echo"
)

type mockPlugin struct {
	Echo *echo.Echo `inject:"echo"`
}

func (p mockPlugin) Initialize() error {
	return nil
}

func TestPlugin(t *testing.T) {
	t.Run("should inject *echo.Echo to the app", func(t *testing.T) {
		app := sen.New()
		m := &mockPlugin{}

		err := app.With(echoPlugin.Bundle(), m)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}

		if m.Echo == nil {
			t.Errorf("Expected Echo to be populated")
		}
	})

	t.Run("should call Shutdown if a hook from OnRun returns an error", func(t *testing.T) {
		hook1Called := 0
		doneCh := make(chan struct{})

		app := sen.New()
		err := app.With(echoPlugin.Bundle())
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}

		_ = app.With(sen.OnRun(func(_ context.Context) error {
			hook1Called++
			return errors.New("run error")
		}))

		go func() {
			err := app.Run(context.Background())
			if fmt.Sprintf("%v", err) != "run error" {
				t.Errorf("Unexpected error: %v", err)
			}
			close(doneCh)
		}()

		select {
		case <-doneCh:
			if hook1Called != 1 {
				t.Errorf("Expected hook1 is called once but got %d", hook1Called)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("test timed out")
		}
	})

	t.Run("should not return an error when it's stopped gracefully", func(t *testing.T) {
		serverIsStarted := make(chan struct{})
		doneCh := make(chan struct{})

		app := sen.New()
		err := app.With(echoPlugin.Bundle())
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}

		_ = app.With(sen.OnRun(func(_ context.Context) error {
			close(serverIsStarted)
			return nil
		}))

		go func() {
			err := app.Run(context.Background())
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			close(doneCh)
		}()

		go func() {
			<-serverIsStarted
			_ = app.Shutdown(context.Background())
		}()

		select {
		case <-doneCh:
		case <-time.After(100 * time.Millisecond):
			t.Errorf("test timed out")
		}
	})
}

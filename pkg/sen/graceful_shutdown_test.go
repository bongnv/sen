package sen_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/bongnv/sen/pkg/sen"
)

type mockWaitingPlugin struct {
	App   *sen.Application `inject:"app"`
	ready chan struct{}
}

func (p *mockWaitingPlugin) Initialize() error {
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
		a, err := sen.New(sen.GracefulShutdown())
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}
		doneCh := make(chan struct{})
		go func() {
			err := a.Run(context.Background())
			if err != nil {
				t.Errorf("Unexpected err %v", err)
			}
			close(doneCh)
		}()

		select {
		case <-doneCh:
		case <-time.After(100 * time.Millisecond):
			t.Errorf("test timed out")
		}
	})

	t.Run("should exit the app if shutdown is called", func(t *testing.T) {
		app, err := sen.New(
			sen.GracefulShutdown(),
			makeMockWaitingPlugin(),
		)
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}

		doneCh := make(chan struct{})
		go func() {
			err := app.Run(context.Background())
			if err != nil {
				t.Errorf("Unexpected err %v", err)
			}
			close(doneCh)
		}()

		err = app.Shutdown(context.Background())
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}
		select {
		case <-doneCh:
		case <-time.After(100 * time.Millisecond):
			t.Errorf("test timed out")
		}
	})

	t.Run("should exit the app if SIGTERM signal is received", func(t *testing.T) {
		m := makeMockWaitingPlugin()
		app, err := sen.New(
			sen.GracefulShutdown(),
			m,
		)
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}

		doneCh := make(chan struct{})
		go func() {
			err := app.Run(context.Background())
			if err != nil {
				t.Errorf("Unexpected err %v", err)
			}
			close(doneCh)
		}()

		<-m.ready
		err = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}
		select {
		case <-doneCh:
		case <-time.After(100 * time.Millisecond):
			t.Errorf("test timed out")
		}
	})
}

package app

import (
	"context"
	"sync"
)

// there are 2 ways it can be implemented:
// Using factory: nicer, isolated
// Injecting immediately, how to handle error?
// we need to inject immediately so application can inject further components while doing it
// how do we handle error? returning error always is annoying

// Hook represents a hook which allows to customise the application life cycle.
type Hook func(ctx context.Context) error

type Application struct {
	Injector

	plugins           Plugin
	runHooks          []Hook
	shutdownHooks     []Hook
	postShutdownHooks []Hook
	shutdownDoneCh    chan struct{}
	shutdownErr       error
	shutdownOnce      sync.Once
	shutdownCh        chan struct{}
}

func New(plugins ...Plugin) *Application {
	app := &Application{
		plugins:        Module(plugins...),
		Injector:       newInjector(),
		shutdownDoneCh: make(chan struct{}),
		shutdownCh:     make(chan struct{}),
	}

	return app
}

// OnRun adds additional logic when the app runs. For a long lasting service
// it should only exit the function when the service no longer runs.
func (app *Application) OnRun(h Hook) {
	app.runHooks = append(app.runHooks, h)
}

// OnShutdown adds additional logic when the app shuts down.
func (app *Application) OnShutdown(h Hook) {
	app.shutdownHooks = append(app.shutdownHooks, h)
}

// Run runs the application by executing all the registered hooks for this phase.
func (app *Application) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := app.Register("app", app); err != nil {
		return err
	}

	if err := app.ApplyPlugin(ctx, app.plugins); err != nil {
		return err
	}

	runErrCh := executeHooks(ctx, app.runHooks)
	var runErr error

	select {
	case <-app.shutdownCh:
		app.shutdownErr = <-executeHooks(ctx, app.shutdownHooks)
		close(app.shutdownDoneCh)
		runErr = <-runErrCh
	case runErr = <-runErrCh:
		app.shutdownErr = <-executeHooks(ctx, app.shutdownHooks)
		close(app.shutdownDoneCh)
	}

	return firstErr(
		runErr,
		app.shutdownErr,
		<-executeHooks(ctx, app.postShutdownHooks),
	)
}

// Shutdown runs the application by executing all the registered hooks for this phase.
// It will include OnShutdown & AfterShutdown.
func (app *Application) Shutdown(ctx context.Context) error {
	app.shutdownOnce.Do(func() {
		close(app.shutdownCh)
	})

	select {
	case <-app.shutdownDoneCh:
		return app.shutdownErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (app *Application) ApplyPlugin(ctx context.Context, p Plugin) error {
	if err := app.Inject(p); err != nil {
		return err
	}

	return p.Apply(ctx)
}

func executeHooks(ctx context.Context, hooks []Hook) <-chan error {
	errCh := make(chan error)
	wg := sync.WaitGroup{}

	for _, h := range hooks {
		wg.Add(1)
		h := h
		go func() {
			if err := h(ctx); err != nil {
				select {
				case errCh <- err:
				case <-ctx.Done():
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return errCh
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

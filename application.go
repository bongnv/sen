package sen

import (
	"context"
	"sync"
)

// Hook represents a hook which allows to customise the application life cycle.
type Hook func(ctx context.Context) error

// Application represents an application.
// To construct an application from plugins use Apply. For example:
// app := sen.New()
// if err := app.Apply(plugin1, plugin2); err != nil {
//    handleError(err)
// }
type Application struct {
	*defaultInjector

	runHooks          []Hook
	shutdownHooks     []Hook
	postShutdownHooks []Hook
	shutdownDoneCh    chan struct{}
	shutdownErr       error
	shutdownOnce      sync.Once
	shutdownCh        chan struct{}
}

// New creates a new Application.
func New() *Application {
	app := &Application{
		defaultInjector: newInjector(),
		shutdownDoneCh:  make(chan struct{}),
		shutdownCh:      make(chan struct{}),
	}

	_ = app.Register("app", app)
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

// Apply applies a plugin or multiple plugins.
// While applying the plugin, Init method will be called.
func (app *Application) Apply(plugins ...Plugin) error {
	p := Module(plugins...)
	if err := app.Inject(p); err != nil {
		return err
	}

	return p.Init()
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

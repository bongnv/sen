package echo

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v8"
	"github.com/labstack/echo/v4"

	"github.com/bongnv/sen"
)

// Module is a sen.Plugin that provides both Config and echo.Echo for convenience.
//
// # Usage
//
// app.With(echo.Module())
func Module() sen.Plugin {
	return sen.Module(
		&ConfigProvider{},
		&Plugin{},
	)
}

// Config includes configuration to initialize an echo server.
type Config struct {
	Port string `env:"PORT,required" envDefault:"1323"`
}

// ConfigProvider is a plugin that provides Config for initializing echo.
//
// # Usage
//
// app.With(&echo.ConfigProvider{})
type ConfigProvider struct {
	Injector sen.Injector `inject:"injector"`
}

// Initialize loads Config from environment variables and
// registers it to the application.
func (p ConfigProvider) Initialize() error {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return err
	}

	return p.Injector.Register("echo.config", cfg)
}

// Plugin is a plugin that provides an instance of echo.Echo.
// The plugin requires Config is registered in advance.
//
// # Usage
//
//	app.With(&echo.Plugin{
//		Middlewares: middlewaresFuncs,
//	})
type Plugin struct {
	Middlewares []echo.MiddlewareFunc

	// will be injected
	LC       sen.Lifecycle `inject:"lifecycle"`
	Injector sen.Injector  `inject:"injector"`
	Cfg      *Config       `inject:"echo.config"`
}

// Initialize initializes and registers the echo.Echo instance with the provided middlewares.
func (p Plugin) Initialize() error {
	e := echo.New()
	e.Use(p.Middlewares...)

	shutdownFn := runOnce(e.Shutdown)

	p.LC.OnRun(func(ctx context.Context) error {
		// since echo doesn't take context, we will need to handle it manually
		// in case the context is cancelled, e.Shutdown() will be called.
		go func() {
			<-ctx.Done()
			_ = shutdownFn(context.Background())
		}()

		return e.Start(fmt.Sprintf(":%s", p.Cfg.Port))
	})

	p.LC.OnShutdown(func(ctx context.Context) error {
		return shutdownFn(ctx)
	})

	return p.Injector.Register("echo", e)
}

// runOnce allows creates a function that will call fn only once.
// It's different from sync.Once that, all calls will be blocked and returns
// the error from the single call of fn.
func runOnce(fn func(ctx context.Context) error) func(ctx context.Context) error {
	var err error
	once := &sync.Once{}
	done := make(chan struct{})
	return func(ctx context.Context) error {
		once.Do(func() {
			err = fn(ctx)
			close(done)
		})
		select {
		case <-done:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

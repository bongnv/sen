package main

import (
	"context"
	"log"
	"net/http"

	echoPlugin "github.com/bongnv/sen/pkg/plugins/echo"
	zapPlugin "github.com/bongnv/sen/pkg/plugins/zap"
	"github.com/bongnv/sen/pkg/sen"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	app := sen.New()
	err := app.With(
		sen.GracefulShutdown(),
		&zapPlugin.Plugin{},
		echoPlugin.Bundle(),
		&Service{},
	)
	if err != nil {
		log.Fatalf("Failed to initialize the app due to %v\n", err)
	}

	if err := app.Run(context.Background()); err != nil {
		log.Printf("Error while running, err: %v\n", err)
	}
}

// Service is an example implementation of a service.
// Its dependencies are injected when registering to sen.
type Service struct {
	Injector sen.Injector  `inject:"injector"`
	Echo     *echo.Echo    `inject:"echo"`
	Logger   *zap.Logger   `inject:"logger"`
	LC       sen.Lifecycle `inject:"lifecycle"`
}

// Initialize initializes the service.
// It also implements sen.Plugin interface so
// it can be used as a plugin.
func (s *Service) Initialize() error {
	s.Logger.Info("The service is initializing")

	// Registering handlers with echo.Echo.
	// Following is an example of registering the handler for GET /hello
	s.Echo.GET("/hello", s.Hello)

	// Registers hooks at OnRun and OnShutdown.
	// Not all services need these.
	s.LC.OnRun(s.Run)
	s.LC.OnShutdown(s.Shutdown)

	// Registering the service under the name "my-service" so
	// it can be injected as a dependency later on.
	return s.Injector.Register("my-service", s)
}

// Run is a hook when the application starts to run.
// If is a long-running service, it should block the function from returning
// until it finishes.
//
// For example, echo.Plugin will block its OnRun hook until the HTTP service stops.
func (s *Service) Run(_ context.Context) error {
	s.Logger.Info("The service is starting")
	return nil
}

// Shutdown is a hook when the application is shutting down.
func (s *Service) Shutdown(_ context.Context) error {
	s.Logger.Info("The service is shutting down")
	return nil
}

// Hello is a handler that handles requests to GET /hello.
func (s *Service) Hello(c echo.Context) error {
	s.Logger.Info("Hello is called")
	return c.String(http.StatusOK, "OK")
}

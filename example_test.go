package sen_test

import (
	"context"
	"fmt"

	"github.com/bongnv/sen"
)

func Example() {
	app := sen.New()

	runHook := sen.OnRun(func(_ context.Context) error {
		fmt.Println("OnRun is executed")
		return nil
	})

	shutdownHook := sen.OnShutdown(func(_ context.Context) error {
		fmt.Println("OnShutdown is executed")
		return nil
	})

	afterRunHook := sen.AfterRun(func(_ context.Context) error {
		fmt.Println("AfterRun is executed")
		return nil
	})

	_ = app.With(runHook, shutdownHook, afterRunHook)
	err := app.Run(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// OnRun is executed
	// OnShutdown is executed
	// AfterRun is executed
}

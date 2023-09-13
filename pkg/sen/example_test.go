package sen_test

import (
	"context"
	"fmt"

	"github.com/bongnv/sen/pkg/sen"
)

func Example() {
	runHook := sen.OnRun(func(_ context.Context) error {
		fmt.Println("OnRun is executed")
		return nil
	})

	shutdownHook := sen.OnShutdown(func(_ context.Context) error {
		fmt.Println("OnShutdown is executed")
		return nil
	})

	postRunHook := sen.PostRun(func(_ context.Context) error {
		fmt.Println("PostRun is executed")
		return nil
	})

	app, err := sen.New(runHook, shutdownHook, postRunHook)
	if err != nil {
		fmt.Println("Failed to initialize the app:", err)
	}
	err = app.Run(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// OnRun is executed
	// OnShutdown is executed
	// PostRun is executed
}

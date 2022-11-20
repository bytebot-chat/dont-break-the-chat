package main

import (
	"fmt"
	"os"

	"github.com/bytebot-chat/dont-break-the-chat/app"
	dbtc "github.com/bytebot-chat/dont-break-the-chat/app"
)

func main() {
	// Parse arguments from command line and environment variables
	config, err := parseArgs()
	if err != nil {
		fmt.Errorf("failed to parse arguments: %w", err)
		os.Exit(1)
	}

	// Create a new app instance
	app, err := dbtc.NewApp(config)
	if err != nil {
		fmt.Errorf("failed to create app: %w", err)
		os.Exit(1)
	}

	// Start the app
	if err := app.Start(); err != nil {
		fmt.Errorf("failed to start app: %w", err)
		os.Exit(1)
	}
}

func parseArgs() (app.Config, error) {
	return app.Config{}, nil
}

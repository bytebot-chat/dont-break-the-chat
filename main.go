package main

import (
	"fmt"
	"os"

	"github.com/bytebot-chat/dont-break-the-chat/app"
)

func main() {
	// Parse arguments from command line and environment variables
	config, err := parseArgs()
	if err != nil {
		fmt.Errorf("failed to parse arguments: %w", err)
		os.Exit(1)
	}

	// Create a new app instance
	app := app.NewApp(config)

	// Get a logger from the app
	logger := app.NewLogger()

	// Connect to redis
	logger.Info().
		Str("host", config.RedisHost()).
		Msg("connecting to redis")
	r, err := app.NewRedis()

}

func parseArgs() (app.Config, error) {
	return app.Args{}, nil
}

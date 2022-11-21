package main

import (
	"fmt"
	"os"

	dbtc "github.com/bytebot-chat/dont-break-the-chat/app"
	"github.com/rs/zerolog"
)

func main() {
	// Parse arguments from command line and environment variables
	config := dbtc.Config{
		RedisHost:     "localhost",
		RedisPort:     6379,
		InboundTopic:  "discord:inbound",
		OutboundTopic: "discord:outbound",
		LogLevel:      zerolog.Level(zerolog.DebugLevel),
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

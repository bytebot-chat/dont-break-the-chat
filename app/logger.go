package app

import (
	"os"

	"github.com/rs/zerolog"
)

func (*App) NewLogger() zerolog.Logger {

	// Create a new logger with the configured log level
	// One thing per line to make it easier to read and track changes
	return zerolog.New(os.Stdout).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Logger()
}

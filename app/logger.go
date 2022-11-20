package app

import (
	"os"

	"github.com/rs/zerolog"
)

func (*App) NewLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

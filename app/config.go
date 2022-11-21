package app

import "github.com/rs/zerolog"

// Config is the configuration for the app.
type Config struct {
	RedisHost     string
	RedisPort     int
	InboundTopic  string
	OutboundTopic string
	LogLevel      zerolog.Level
}

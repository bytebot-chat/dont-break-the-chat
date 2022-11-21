package app

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
)

// App is the struct that represents the configuration and state of the app.
type App struct {
	Config Config

	redis   *Redis
	context context.Context
	logger  zerolog.Logger
}

// The start method wraps all the tasks necessary to start and manage the app.
// It should be considered the "main" method of the app.
// It also manages its own state, so it can be called multiple times.
func (a *App) Start() error {
	// Create a new logger
	a.logger = a.NewLogger()

	// Connect to Redis
	a.logger.Info().
		Msg("connecting to redis")
	redis, err := newRedis(a.Config.RedisHost, a.Config.RedisPort, "", 0)
	if err != nil {
		return err
	}
	a.redis = redis
	a.logger.Info().
		Msg("connected to redis!")

	// Start the inbound listener in an anonymous goroutine
	go func() {
		a.logger.Info().
			Msg("starting inbound listener")

		// Subscribe to the inbound topic on the redis pubsub
		topic := a.redis.Subscribe(a.context, a.Config.InboundTopic)

		// Create a go channel to receive messages from the topic
		channel := topic.Channel()

		// Iterate over the messages from the channel
		for msg := range channel {
			a.logger.Debug().
				Msg("received message")

			// Unmarshal the message into a Message struct
			m, err := unmarshalIncomingMessage(msg)
			if err != nil {
				a.logger.Error().
					Err(err).
					Msg("failed to unmarshal message")
				continue // Skip this message and continue if we can't unmarshal it
			}

			// Handle the message. We are only interested in message that start with a command prefix for now.
			if strings.HasPrefix(m.Content, "!") {
				a.logger.Debug().
					Msg("message is a command")

				// the entrypoint for handling commands. located in app/commands.go
				err = handleCommand(a, m)
				if err != nil {
					a.logger.Error().
						Err(err).
						Msg("error handling command")
				}
			} else {
				a.logger.Debug().
					Msg("message is not a command")
			}
		}
	}()

	a.logger.Info().
		Msg("Holding the app open")
	for {
		//
	}
	return nil
}

// NewApp creates a new app instance with the given configuration.
func NewApp(config Config) (*App, error) {
	return &App{
		Config:  config,
		context: context.Background(),
	}, nil
}

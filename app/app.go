package app

import (
	"context"

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
	redis, err := newRedis()
	if err != nil {
		return err
	}
	a.redis = redis

	// Connect to discord
	a.logger.Info().
		Msg("connecting to discord")

	// Start the inbound listener in an anonymous goroutine
	go func() {
		a.logger.Info().
			Msg("starting inbound listener")
		topic := a.redis.Subscribe(a.context, a.Config.InboundTopic)
		channel := topic.Channel()
		for msg := range channel {
			handleIncomingMessage(msg)
			a.logger.Info().
				Msg("received message")
		}
	}()

	return nil
}

// Stop stops the app.
func (*App) Stop() error {
	return nil
}

// Health returns a health check for the app.
// May be polled by a load balancer or other monitoring service.
func (*App) Health() error {
	return nil
}

// Ready returns a readiness check for the app. This is used to indicate whether
// the app is ready to receive traffic.
func (*App) Ready() error {
	return nil
}

// NewApp creates a new app instance with the given configuration.
func NewApp(config Config) (*App, error) {
	return &App{
		Config:  config,
		context: context.Background(),
	}, nil
}

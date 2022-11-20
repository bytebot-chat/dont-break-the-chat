package app

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type Config struct {
	RedisHost string
}

type App struct {
	Config Config
}

type Redis struct {
	*redis.Client
}

func (*App) Start() error {
	return nil
}

func (*App) Stop() error {
	return nil
}

func (*App) Health() error {
	return nil
}

func (*App) Ready() error {
	return nil
}

func (*App) NewLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func NewApp(config Config) *App {
	return &App{
		Config: config,
	}
}

func NewRedis() (*Redis, error) {
	return &Redis{}, nil
}

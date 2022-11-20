package app

import "github.com/go-redis/redis/v8"

type Redis struct {
	*redis.Client
}

func newRedis() (*Redis, error) {
	return &Redis{}, nil
}

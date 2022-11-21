package app

import (
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	*redis.Client
}

func newRedis(host string, port int, password string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: password,
		DB:       db,
	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}
	return &Redis{client}, nil
}

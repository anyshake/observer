package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func OpenRedis(host string, port int, password string, db int, enable bool) (*redis.Client, error) {
	if !enable {
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CloseRedis(d *redis.Client) error {
	return d.Close()
}

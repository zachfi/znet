package znet

import (
	"fmt"

	"github.com/go-redis/redis"
)

func NewRedisClient(host string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", host),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client, nil
}

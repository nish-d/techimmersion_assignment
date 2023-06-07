package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	cli *redis.Client
}

func (client *RedisClient) InitRedisClient(context context.Context, address string, password string) error {
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,
	})

	//test and return err if any
	err := redisClient.Ping(context).Err()
	if err != nil {
		return err
	}

	//set client and return
	client.cli = redisClient
	return nil

}

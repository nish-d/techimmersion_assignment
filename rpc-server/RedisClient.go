package main

import (
	"context"
	"encoding/json"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	cli *redis.Client
}

var (
	redisClient = &RedisClient{}
)

func (cli *RedisClient) InitRedisClient(context context.Context, address string, password string) error {
	redisClient.cli = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,
	})

	//test and return err if any
	err := redisClient.cli.Ping(context).Err()
	if err != nil {
		return err
	}

	////set client and return
	//redisClient.cli = redisClient
	return nil

}

func (cli *RedisClient) SaveMessageToRedis(context context.Context, groupId string, message *rpc.Message) {

	text, _ := json.Marshal(message)

	item := &redis.Z{
		Score:  float64(message.SendTime),
		Member: text,
	}

	redisClient.cli.ZAdd(context, groupId, *item).Result()

}

func (cli *RedisClient) GetMessagesByGroupID(ctx context.Context, groupID string, start, end int64, reverse bool) ([]*rpc.Message, error) {
	var (
		rawMessages []string
		messages    []*rpc.Message
		err         error
	)

	if reverse {
		// Desc order with time -> first message is the latest message
		rawMessages, err = redisClient.cli.ZRevRange(ctx, groupID, start, end).Result()
	} else {
		// Asc order with time -> first message is the earliest message
		rawMessages, err = redisClient.cli.ZRange(ctx, groupID, start, end).Result()
	}
	if err != nil {
		return nil, err
	}

	for _, msg := range rawMessages {
		temp := &rpc.Message{}
		err := json.Unmarshal([]byte(msg), temp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, temp)
	}

	return messages, nil
}

package utils

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const (
	RedisAddr     = "43.139.116.74:6379"
	RedisPassword = "fanqie"
	RedisDB       = 0
)

var Red *redis.Client

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPassword, // no password set
		DB:       RedisDB,       // use default DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("Redis connect error, " + err.Error())
	} else {
		Red = client
	}
}

const (
	PublishKey = "websocket"
)

// 发送消息到redis
func Publish(ctx context.Context, channel, msg string) error {
	fmt.Println("Publish: ", msg)
	err := Red.Publish(ctx, channel, msg).Err()
	return err
}

// 订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe: ", sub)
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("msg: ", msg.Payload)
	return msg.Payload, err
}

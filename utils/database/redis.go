package database

import (
	"context"
	"fmt"
	"gain-v2/configs"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

func InitRedis(c configs.ProgrammingConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", c.DBRedisAddress, c.DBRedisPort),
		Password: c.DBRedisPassword,
		DB:       c.DBRedisDatabase,
	})

	ctx := context.Background()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Error("Terjadi kesalahan pada Redis, error:", err.Error())
		return nil, err
	}

	fmt.Println("Redis result: ", pong)

	return rdb, nil
}

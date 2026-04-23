package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var RedisClient *redis.Client
var RedisContext = context.Background()

func InitRedis() {

	redisConfig := Conf.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: "", // 密码
		DB:       0,  // 数据库编号

		// 连接池配置
		PoolSize:     redisConfig.PoolSize, // 最大连接数
		MinIdleConns: 10,                   // 最小空闲连接数
		DialTimeout:  5 * time.Second,      // 连接超时
		ReadTimeout:  3 * time.Second,      // 读超时
		WriteTimeout: 3 * time.Second,      // 写超时
		PoolTimeout:  4 * time.Second,      // 连接池超时

	})

	pong, err := client.Ping(RedisContext).Result()

	if err != nil {
		fmt.Println(err)
		return
	}

	RedisClient = client
	fmt.Printf("redis connect success:%s\n", pong)

}

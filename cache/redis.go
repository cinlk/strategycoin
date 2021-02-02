package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

// 扩展 TODO
func initRedis() {

	r := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.address"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.db"),
		WriteTimeout: viper.GetDuration("redis.wirteTimeout") * time.Second,
		IdleTimeout:  viper.GetDuration("redis.IdleTimeout") * time.Second,
	})

	_, err := r.Ping(context.Background()).Result()

	if err != nil {
		panic(err)
	}
	return

}

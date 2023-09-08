package redis

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// var server = "invoice_redis:6379"
// var password = "invoice_1234"
func NewRedisDb() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_HOST") + ":" + viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"), // no password set
		DB:       viper.GetInt("REDIS_DB"),          // use default DB
	})
}

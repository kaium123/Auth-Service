package redis

import (
	"auth/common/logger"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func NewRedisDb() *redis.Client {
    redisClient := redis.NewClient(&redis.Options{
        Addr:     "social_media_redis:6379", // Use the service name directly
        Password: viper.GetString("REDIS_PASSWORD"), // no password set
        DB:       viper.GetInt("REDIS_DB"),          // use default DB
    })
    if s, err := redisClient.Ping().Result(); err != nil {
        logger.LogError("Error in petty-cash redis, ", err)
    } else {
        logger.LogInfo(s)
    }

    return redisClient
}


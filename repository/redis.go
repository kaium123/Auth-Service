package repository

import (
	"auth/common/logger"
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RedisRepositoryInterface interface {
	Set(ctx context.Context, key string, value map[string]interface{}) error
	GetAll(ctx context.Context, key string) (map[string]string, error)
	GetSingleData(ctx context.Context, key string, field string) (string, error)
	GetByFields(ctx context.Context, key string, fields []string) ([]string, error)
	Delete(ctx context.Context, key string, fields []string) error
	SetExpire(ctx context.Context, key string, expireTime time.Duration) error
	Exists(ctx context.Context, key string, field string) bool
	SetValue(ctx context.Context, key string, value string, expiration time.Duration) error
}

type RedisRepository struct {
	RedisClient *redis.Client
	logger      logger.LoggerInterface
}

func NewRedisRepository(redisClient *redis.Client, logger logger.LoggerInterface) RedisRepositoryInterface {
	return &RedisRepository{redisClient, logger}
}

func (r RedisRepository) SetValue(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := r.RedisClient.Set(key, value, expiration).Err()
	if err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

func (r RedisRepository) Set(ctx context.Context, key string, value map[string]interface{}) error {
	logger.LogInfo("set to redis ", key)
	err := r.RedisClient.HMSet(key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisRepository) Exists(ctx context.Context, key string, field string) bool {
	var exists bool
	if len(field) != 0 {
		exists = r.RedisClient.HExists(key, field).Val()
	} else {
		check := r.RedisClient.Exists(key).Val()
		exists = check == 1
	}
	return exists
}

func (r RedisRepository) SetExpire(ctx context.Context, key string, expireTime time.Duration) error {
	err := r.RedisClient.Expire(key, expireTime).Err()

	return err
}

func (r RedisRepository) Delete(ctx context.Context, key string, fields []string) error {
	var err error
	if len(fields) == 0 || fields == nil {
		err = r.RedisClient.Del(key).Err()
	} else {
		err = r.RedisClient.HDel(key, fields...).Err()
	}
	if err != nil {
		return err
	}

	return nil
}

func (r RedisRepository) GetAll(ctx context.Context, key string) (map[string]string, error) {
	logger.LogInfo("get from redis key ", key)
	data, err := r.RedisClient.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RedisRepository) GetSingleData(ctx context.Context, key string, field string) (string, error) {
	logger.LogInfo("get from redis repository ", field, " ", key)
	data, err := r.RedisClient.HGet(key, field).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

func (r RedisRepository) GetByFields(ctx context.Context, key string, fields []string) ([]string, error) {
	logger.LogInfo("get from redis key ", key)
	pipe := r.RedisClient.Pipeline()
	cmds := map[string]*redis.StringCmd{}
	for _, field := range fields {
		cmds[field] = pipe.HGet(key, field)
	}
	_, err := pipe.Exec()
	if err != nil {
		return nil, err
	}

	data := []string{}
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			return nil, err
		}
		data = append(data, val)
	}

	return data, nil
}

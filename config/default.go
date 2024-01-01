package config

import (
	"time"
)

var Config Conf

type Conf struct {
	PostgreDriver  string `mapstructure:"POSTGRES_DRIVER"`
	PostgresSource string `mapstructure:"POSTGRES_SOURCE"`

	Port string `mapstructure:"PORT"`

	Origin string `mapstructure:"ORIGIN"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
	QueueService           string        `mapstructure:"QUEUE_SERVICE"`
	RedisHost              string        `mapstructure:"REDIS_HOST"`
	RedisDb                int           `mapstructure:"REDIS_PORT"`
}

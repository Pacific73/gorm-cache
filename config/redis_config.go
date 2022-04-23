package config

import (
	"github.com/go-redis/redis"
)

type RedisConfigMode int

const (
	RedisConfigModePort RedisConfigMode = 0
	RedisConfigModeRaw  RedisConfigMode = 1
)

type RedisConfig struct {
	Mode RedisConfigMode

	Host string
	Port int
	Password string

	Client *redis.Client
}

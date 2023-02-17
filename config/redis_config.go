package config

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

type RedisConfigMode int

const (
	RedisConfigModeOptions RedisConfigMode = 0
	RedisConfigModeRaw     RedisConfigMode = 1
)

type RedisConfig struct {
	Mode RedisConfigMode

	Options *redis.Options
	Client  *redis.Client

	once sync.Once
}

func (c *RedisConfig) InitClient() *redis.Client {
	c.once.Do(func() {
		if c.Mode == RedisConfigModeOptions {
			c.Client = redis.NewClient(c.Options)
		}
	})
	return c.Client
}

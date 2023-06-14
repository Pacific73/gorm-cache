package cache

import (
	"fmt"

	"github.com/Pacific73/gorm-cache/config"
	"github.com/go-redis/redis/v8"
)

func NewGorm2Cache(cacheConfig *config.CacheConfig) (*Gorm2Cache, error) {
	if cacheConfig == nil {
		return nil, fmt.Errorf("you pass a nil config")
	}
	cache := &Gorm2Cache{
		Config: cacheConfig,
	}
	err := cache.Init()
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewRedisConfigWithOptions(options *redis.Options) *config.RedisConfig {
	return &config.RedisConfig{
		Mode:    config.RedisConfigModeOptions,
		Options: options,
	}
}

func NewRedisConfigWithClient(client *redis.Client) *config.RedisConfig {
	return &config.RedisConfig{
		Mode:   config.RedisConfigModeRaw,
		Client: client,
	}
}

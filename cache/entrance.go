package cache

import (
	"github.com/Pacific73/gorm-cache/config"
	"github.com/go-redis/redis"
)

func NewGorm2Cache(cacheConfig *config.CacheConfig) (*Gorm2Cache, error) {
	logger := config.DefaultLogger
	if cacheConfig.DebugLogger != nil {
		logger = cacheConfig.DebugLogger
	}
	cache := &Gorm2Cache{
		Config: cacheConfig,
		Logger: logger,
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

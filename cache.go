package gorm_cache

import (
	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/go-redis/redis"
)

func NewGorm2Cache(cacheConfig *config.CacheConfig) *cache.Gorm2Cache {
	logger := config.DefaultLogger
	if cacheConfig.DebugLogger != nil {
		logger = cacheConfig.DebugLogger
	}

	return &cache.Gorm2Cache{
		Config: cacheConfig,
		Logger: logger,
	}
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

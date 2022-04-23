package gorm_cache

import (
	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/go-redis/redis"
)

func NewGormCache(config *config.CacheConfig) *cache.Gorm2Cache {
	// TODO
	return nil
}

func NewRedisConfigWithHostPort(host string, port int) *config.RedisConfig {
	return &config.RedisConfig{
		Mode:     config.RedisConfigModePort,
		Host:     host,
		Port:     port,
	}
}

func NewRedisConfigWithPassword(host string, port int, password string) *config.RedisConfig {
	return &config.RedisConfig{
		Mode:     config.RedisConfigModePort,
		Host:     host,
		Port:     port,
		Password: password,
	}
}

func NewRedisConfigWithClient(client *redis.Client) *config.RedisConfig {
	return &config.RedisConfig{
		Mode:     config.RedisConfigModeRaw,
		Client:   client,
	}
}

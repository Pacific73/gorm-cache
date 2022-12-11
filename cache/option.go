package cache

import "github.com/Pacific73/gorm-cache/config"

type Option func(p *config.CacheConfig)

func WithLevel(level config.CacheLevel) Option {
	return func(p *config.CacheConfig) {
		p.CacheLevel = level
	}
}
func WithStorage(storage config.CacheStorage) Option {
	return func(p *config.CacheConfig) {
		p.CacheStorage = storage
	}
}
func WithRedisConfig(redisClient *config.RedisConfig) Option {
	return func(p *config.CacheConfig) {
		p.RedisConfig = redisClient
	}
}
func WithTables(tables []string) Option {
	return func(p *config.CacheConfig) {
		p.Tables = tables
	}
}
func WithInvalidateWhenUpdate(isBool bool) Option {
	return func(p *config.CacheConfig) {
		p.InvalidateWhenUpdate = isBool
	}
}
func WithCacheTTL(ttl int64) Option {
	return func(p *config.CacheConfig) {
		p.CacheTTL = ttl
	}
}
func WithCacheMaxItemCnt(cnt int64) Option {
	return func(p *config.CacheConfig) {
		p.CacheMaxItemCnt = cnt
	}
}

func WithCacheSize(size int) Option {
	return func(p *config.CacheConfig) {
		p.CacheSize = size
	}
}
func WithDebugMode(debug bool) Option {
	return func(p *config.CacheConfig) {
		p.DebugMode = debug
	}
}
func WithDebugLogger(log config.LoggerInterface) Option {
	return func(p *config.CacheConfig) {
		p.DebugLogger = log
	}
}

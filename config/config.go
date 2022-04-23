package config

type CacheConfig struct {
	// CacheLevel there're 2 types of cache and 4 kinds of cache option
	CacheLevel CacheLevel

	// CacheStorage choose proper storage medium
	CacheStorage CacheStorage

	// RedisConfig if storage is redis, then this config needs to be setup
	RedisConfig *RedisConfig

	// Tables only cache data within given data tables (cache all if empty)
	Tables []string

	// InvalidateWhenUpdate
	// if user update/delete/create something in DB, we invalidate all cached data to ensure consistency,
	// else we do nothing to outdated cache.
	InvalidateWhenUpdate bool

	// CacheTTL cache ttl in ms, where 0 represents forever
	CacheTTL int64

	// PrimaryCacheSize cache maximal size for primary cache, in MB
	PrimaryCacheSize int

	// SearchCacheSize cache maximal size for search cache, in MB
	SearchCacheSize int

	// DebugMode indicate if we're in debug mode (will print access log)
	DebugMode bool

	// DebugLogger
	DebugLogger LoggerInterface
}

type CacheLevel int

const (
	CacheLevelOff         CacheLevel = 0
	CacheLevelOnlyPrimary CacheLevel = 1
	CacheLevelOnlySearch  CacheLevel = 2
	CacheLevelAll         CacheLevel = 3
)

type CacheStorage int

const (
	CacheStorageMemory CacheStorage = 0
	CacheStorageRedis  CacheStorage = 1
)

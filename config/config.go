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

	// CacheMaxItemCnt for given query, if objects retrieved are more than this cnt,
	// then we choose not to cache for this query. 0 represents caching all queries.
	CacheMaxItemCnt int64

	// CacheSize maximal items in primary cache (only works in MEMORY storage)
	CacheSize int

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

package cache

import (
	"context"
	"sync/atomic"

	"github.com/Pacific73/gorm-cache/util"

	"github.com/Pacific73/gorm-cache/config"
	"github.com/Pacific73/gorm-cache/data_layer"
	"gorm.io/gorm"
)

type Gorm2Cache struct {
	Config     *config.CacheConfig
	Logger     config.LoggerInterface
	InstanceId string

	db           *gorm.DB
	primaryCache data_layer.DataLayerInterface
	searchCache  data_layer.DataLayerInterface
	hitCount     int64
}

func (c *Gorm2Cache) AttachToDB(db *gorm.DB) {

	_ = db.Callback().Create().After("*").Register("gorm:cache:after_create", AfterCreate(c))
	_ = db.Callback().Delete().After("*").Register("gorm:cache:after_delete", AfterDelete(c))
	_ = db.Callback().Update().After("*").Register("gorm:cache:after_update", AfterUpdate(c))

	_ = db.Callback().Query().Before("gorm:query").Register("gorm:cache:before_query", BeforeQuery(c))
	_ = db.Callback().Query().After("*").Register("gorm:cache:after_query", AfterQuery(c))
}

func (c *Gorm2Cache) Init() error {
	if c.Config.CacheStorage == config.CacheStorageRedis {
		if c.Config.RedisConfig == nil {
			panic("please init redis config!")
		}
		c.Config.RedisConfig.InitClient()
	}
	c.InstanceId = util.GenInstanceId()

	prefix := util.GormCachePrefix + ":" + c.InstanceId

	if c.Config.CacheStorage == config.CacheStorageRedis {
		c.primaryCache = &data_layer.RedisLayer{}
		c.searchCache = &data_layer.RedisLayer{}
	}

	if c.Config.DebugLogger == nil {
		c.Config.DebugLogger = &config.DefaultLoggerImpl{}
	}
	c.Logger = c.Config.DebugLogger
	c.Logger.SetIsDebug(c.Config.DebugMode)

	err := c.primaryCache.Init(c.Config, prefix)
	if err != nil {
		c.Logger.CtxError(context.Background(), "[Init] primary cache init error: %v", err)
		return err
	}
	err = c.searchCache.Init(c.Config, prefix)
	if err != nil {
		c.Logger.CtxError(context.Background(), "[Init] search cache init error: %v", err)
		return err
	}
	return nil
}

func (c *Gorm2Cache) GetHitCount() int64 {
	return atomic.LoadInt64(&c.hitCount)
}

func (c *Gorm2Cache) ResetHitCount() {
	atomic.StoreInt64(&c.hitCount, 0)
}

func (c *Gorm2Cache) IncrHitCount() {
	atomic.AddInt64(&c.hitCount, 1)
}

func (c *Gorm2Cache) ResetCache() error {
	c.ResetHitCount()
	ctx := context.Background()
	err := c.searchCache.CleanCache(ctx)
	if err != nil {
		c.Logger.CtxError(ctx, "[ResetCache] reset search cache error: %v", err)
		return err
	}
	return c.primaryCache.CleanCache(ctx)
}

func (c *Gorm2Cache) InvalidateSearchCache(ctx context.Context, tableName string) error {
	return c.searchCache.DeleteKeysWithPrefix(ctx, util.GenSearchCachePrefix(c.InstanceId, tableName))
}

func (c *Gorm2Cache) InvalidatePrimaryCache(ctx context.Context, tableName string, primaryKey string) error {
	return c.primaryCache.DeleteKey(ctx, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
}

func (c *Gorm2Cache) BatchInvalidatePrimaryCache(ctx context.Context, tableName string, primaryKeys []string) error {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.primaryCache.BatchDeleteKeys(ctx, cacheKeys)
}

func (c *Gorm2Cache) InvalidateAllPrimaryCache(ctx context.Context, tableName string) error {
	return c.primaryCache.DeleteKeysWithPrefix(ctx, util.GenPrimaryCachePrefix(c.InstanceId, tableName))
}

func (c *Gorm2Cache) BatchPrimaryKeyExists(ctx context.Context, tableName string, primaryKeys []string) (bool, error) {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.primaryCache.BatchKeyExist(ctx, cacheKeys)
}

func (c *Gorm2Cache) SearchKeyExists(ctx context.Context, tableName string, SQL string, vars ...interface{}) (bool, error) {
	cacheKey := util.GenSearchCacheKey(c.InstanceId, tableName, SQL, vars...)
	return c.searchCache.KeyExists(ctx, cacheKey)
}

func (c *Gorm2Cache) BatchSetPrimaryKeyCache(ctx context.Context, tableName string, kvs []util.Kv) error {
	for idx, kv := range kvs {
		kvs[idx].Key = util.GenPrimaryCacheKey(c.InstanceId, tableName, kv.Key)
	}
	return c.primaryCache.BatchSetKeys(ctx, kvs)
}

func (c *Gorm2Cache) SetSearchCache(ctx context.Context, cacheValue string, tableName string,
	sql string, vars ...interface{}) error {
	key := util.GenSearchCacheKey(c.InstanceId, tableName, sql, vars...)
	return c.searchCache.SetKey(ctx, util.Kv{
		Key:   key,
		Value: cacheValue,
	})
}

func (c *Gorm2Cache) GetSearchCache(ctx context.Context, tableName string, sql string, vars ...interface{}) (string, error) {
	key := util.GenSearchCacheKey(c.InstanceId, tableName, sql, vars...)
	return c.searchCache.GetValue(ctx, key)
}

func (c *Gorm2Cache) BatchGetPrimaryCache(ctx context.Context, tableName string, primaryKeys []string) ([]string, error) {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.primaryCache.BatchGetValues(ctx, cacheKeys)
}

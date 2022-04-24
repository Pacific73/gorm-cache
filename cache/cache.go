package cache

import (
	"context"
	"sync/atomic"

	"github.com/Pacific73/gorm-cache/util"

	"github.com/Pacific73/gorm-cache/callback"
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
	c.db = db

	c.db.Callback().Create().After("gorm:after_create").Register("gorm:cache:after_create", callback.AfterCreate(c))

	c.db.Callback().Delete().After("gorm:after_delete").Register("gorm:cache:after_delete", callback.AfterDelete(c))

	c.db.Callback().Update().After("gorm:after_update").Register("gorm:cache:after_update", callback.AfterUpdate(c))

	c.db.Callback().Query().Before("gorm:query").Register("gorm:cache:before_query", callback.BeforeQuery(c))
	c.db.Callback().Query().After("gorm:after_query").Register("gorm:cache:after_query", callback.AfterQuery(c))

	c.db.Callback().Row().Before("gorm:row").Register("gorm:cache:before_row_query", callback.BeforeRow(c))
	c.db.Callback().Row().After("gorm:row").Register("gorm:cache:after_row_query", callback.AfterRow(c))

	c.init()
}

func (c *Gorm2Cache) init() {
	if c.Config.CacheStorage == config.CacheStorageRedis {
		if c.Config.RedisConfig == nil {
			panic("please init redis config!")
		}
		c.Config.RedisConfig.InitClient()
	}
	c.InstanceId = util.GenInstanceId()

	prefix := util.GormCachePrefix + ":" + c.InstanceId

	c.primaryCache.Init(c.Config, prefix)
	c.searchCache.Init(c.Config, prefix)

}

func (c *Gorm2Cache) GetHitCount() int64 {
	return atomic.LoadInt64(&c.hitCount)
}

func (c *Gorm2Cache) ResetHitCount() {
	atomic.StoreInt64(&c.hitCount, 0)
}

func (c *Gorm2Cache) ResetCache() {
	c.searchCache.CleanCache()
	c.primaryCache.CleanCache()
}

func (c *Gorm2Cache) InvalidateSearchCache(ctx context.Context, tableName string) {
	c.searchCache.DeleteKeysWithPrefix(ctx, util.GenSearchCachePrefix(c.InstanceId, tableName))
}

func (c *Gorm2Cache) InvalidatePrimaryCache(ctx context.Context, tableName string, primaryKey string) {
	c.primaryCache.DeleteKey(ctx, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
}

func (c *Gorm2Cache) BatchInvalidatePrimaryCache(ctx context.Context, tableName string, primaryKeys []string) {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	c.primaryCache.BatchDeleteKeys(ctx, cacheKeys)
}

func (c *Gorm2Cache) InvalidateAllPrimaryCache(ctx context.Context, tableName string) {
	c.primaryCache.DeleteKeysWithPrefix(ctx, util.GenPrimaryCachePrefix(c.InstanceId, tableName))
}

func (c *Gorm2Cache) BatchPrimaryKeyExists(ctx context.Context, tableName string, primaryKeys []string) bool {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.primaryCache.BatchKeyExist(ctx, cacheKeys)
}

func (c *Gorm2Cache) SearchKeyExists(ctx context.Context, tableName string, SQL string, vars ...interface{}) bool {
	cacheKey := util.GenSearchCacheKey(c.InstanceId, tableName, SQL, vars...)
	return c.searchCache.KeyExists(ctx, cacheKey)
}

func (c *Gorm2Cache) BatchSetPrimaryKeyCache(ctx context.Context, tableName string, kvs []util.Kv) {
	for _, kv := range kvs {
		kv.Key = util.GenPrimaryCacheKey(c.InstanceId, tableName, kv.Key)
	}
	c.primaryCache.BatchSetKeys(ctx, kvs)
}

func (c *Gorm2Cache) SetSearchCache(ctx context.Context, cacheValue string, tableName string, sql string, vars ...interface{}) {
	key := util.GenSearchCacheKey(c.InstanceId, tableName, sql, vars...)
	c.searchCache.SetKey(ctx, util.Kv{
		Key:   key,
		Value: cacheValue,
	})
}

func (c *Gorm2Cache) GetSearchCache(ctx context.Context, tableName string, sql string, vars ...interface{}) string {
	key := util.GenSearchCacheKey(c.InstanceId, tableName, sql, vars...)
	return c.searchCache.GetValue(ctx, key)
}

func (c *Gorm2Cache) BatchGetPrimaryCache(ctx context.Context, tableName string, primaryKeys []string) []string {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.primaryCache.BatchGetValues(ctx, cacheKeys)
}

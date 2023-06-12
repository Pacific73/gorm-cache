package cache

import (
	"context"
	"sync/atomic"

	"github.com/herisy/gorm-cache/config"
	"github.com/herisy/gorm-cache/data_layer"
	"github.com/herisy/gorm-cache/util"
	"gorm.io/gorm"
)

var (
	_ gorm.Plugin = &Gorm2Cache{}
)

type Gorm2Cache struct {
	Config     *config.CacheConfig
	Logger     config.LoggerInterface
	InstanceId string

	db       *gorm.DB
	cache    data_layer.DataLayerInterface
	hitCount int64
}

func (c *Gorm2Cache) Name() string {
	return util.GormCachePrefix
}

func (c *Gorm2Cache) Initialize(db *gorm.DB) (err error) {
	err = db.Callback().Create().After("*").Register("gorm:cache:after_create", AfterCreate(c))
	if err != nil {
		return err
	}

	err = db.Callback().Delete().After("*").Register("gorm:cache:after_delete", AfterDelete(c))
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("*").Register("gorm:cache:after_update", AfterUpdate(c))
	if err != nil {
		return err
	}

	err = db.Callback().Query().Before("gorm:query").Register("gorm:cache:before_query", BeforeQuery(c))
	if err != nil {
		return err
	}

	err = db.Callback().Query().After("*").Register("gorm:cache:after_query", AfterQuery(c))
	if err != nil {
		return err
	}

	return
}

func (c *Gorm2Cache) AttachToDB(db *gorm.DB) {
	c.Initialize(db)
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
		c.cache = &data_layer.RedisLayer{}
	} else if c.Config.CacheStorage == config.CacheStorageMemory {
		c.cache = &data_layer.MemoryLayer{}
	}

	if c.Config.DebugLogger == nil {
		c.Config.DebugLogger = &config.DefaultLoggerImpl{}
	}
	c.Logger = c.Config.DebugLogger
	c.Logger.SetIsDebug(c.Config.DebugMode)

	err := c.cache.Init(c.Config, prefix)
	if err != nil {
		c.Logger.CtxError(context.Background(), "[Init] cache init error: %v", err)
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
	err := c.cache.CleanCache(ctx)
	if err != nil {
		c.Logger.CtxError(ctx, "[ResetCache] reset cache error: %v", err)
		return err
	}
	return nil
}

func (c *Gorm2Cache) InvalidateSearchCache(ctx context.Context, tableName string) error {
	return c.cache.DeleteKeysWithPrefix(ctx, util.GenSearchCachePrefix(c.InstanceId, tableName))
}

func (c *Gorm2Cache) InvalidatePrimaryCache(ctx context.Context, tableName string, primaryKey string) error {
	return c.cache.DeleteKey(ctx, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
}

func (c *Gorm2Cache) BatchInvalidatePrimaryCache(ctx context.Context, tableName string, primaryKeys []string) error {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.cache.BatchDeleteKeys(ctx, cacheKeys)
}

func (c *Gorm2Cache) InvalidateAllPrimaryCache(ctx context.Context, tableName string) error {
	return c.cache.DeleteKeysWithPrefix(ctx, util.GenPrimaryCachePrefix(c.InstanceId, tableName))
}

func (c *Gorm2Cache) BatchPrimaryKeyExists(ctx context.Context, tableName string, primaryKeys []string) (bool, error) {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.cache.BatchKeyExist(ctx, cacheKeys)
}

func (c *Gorm2Cache) SearchKeyExists(ctx context.Context, tableName string, SQL string, vars ...interface{}) (bool, error) {
	cacheKey := util.GenSearchCacheKey(c.InstanceId, tableName, SQL, vars...)
	return c.cache.KeyExists(ctx, cacheKey)
}

func (c *Gorm2Cache) BatchSetPrimaryKeyCache(ctx context.Context, tableName string, kvs []util.Kv) error {
	for idx, kv := range kvs {
		kvs[idx].Key = util.GenPrimaryCacheKey(c.InstanceId, tableName, kv.Key)
	}
	return c.cache.BatchSetKeys(ctx, kvs)
}

func (c *Gorm2Cache) SetSearchCache(ctx context.Context, cacheValue string, tableName string,
	sql string, vars ...interface{}) error {
	key := util.GenSearchCacheKey(c.InstanceId, tableName, sql, vars...)
	return c.cache.SetKey(ctx, util.Kv{
		Key:   key,
		Value: cacheValue,
	})
}

func (c *Gorm2Cache) GetSearchCache(ctx context.Context, tableName string, sql string, vars ...interface{}) (string, error) {
	key := util.GenSearchCacheKey(c.InstanceId, tableName, sql, vars...)
	return c.cache.GetValue(ctx, key)
}

func (c *Gorm2Cache) BatchGetPrimaryCache(ctx context.Context, tableName string, primaryKeys []string) ([]string, error) {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, util.GenPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.cache.BatchGetValues(ctx, cacheKeys)
}

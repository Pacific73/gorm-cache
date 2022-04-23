package cache

import (
	"context"

	"github.com/Pacific73/gorm-cache/key"

	"github.com/Pacific73/gorm-cache/callback"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/Pacific73/gorm-cache/data_layer"
	"gorm.io/gorm"
)

type Gorm2Cache struct {
	Config *config.CacheConfig
	Logger config.LoggerInterface

	db        *gorm.DB
	dataLayer data_layer.DataLayerInterface
	hitCount  int64
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

}

func (c *Gorm2Cache) GetHitCount() int64 {
	return c.hitCount
}

func (c *Gorm2Cache) InvalidateSearchCache(ctx context.Context, tableName string) {
	err := c.dataLayer.DeleteKeysWithPrefix(ctx, key.GenSearchCachePrefix(tableName))
	if err != nil {
		c.Logger.CtxDebug(ctx, "[InvalidateSearchCache] invalidating search cache for table %s error: %v",
			tableName, err)
	}
}

func (c *Gorm2Cache) InvalidatePrimaryCache(ctx context.Context, tableName string, primaryKey string) {
	err := c.dataLayer.DeleteKey(ctx, key.GenPrimaryCacheKey(tableName, primaryKey))
	if err != nil {
		c.Logger.CtxDebug(ctx, "[InvalidatePrimaryCache] invalidating primary cache for key %s:%s error: %v",
			tableName, primaryKey, err)
	}
}

func (c *Gorm2Cache) BatchInvalidatePrimaryCache(ctx context.Context, tableName string, primaryKeys []string) {

	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, key.GenPrimaryCacheKey(tableName, primaryKey))
	}

	err := c.dataLayer.BatchDeleteKeys(ctx, cacheKeys)
	if err != nil {
		c.Logger.CtxDebug(ctx, "[BatchInvalidatePrimaryCache] batch invalidating primary cache for keys %v error: %v",
			primaryKeys, err)
	}
}

func (c *Gorm2Cache) InvalidateAllPrimaryCache(ctx context.Context, tableName string) {
	err := c.dataLayer.DeleteKeysWithPrefix(ctx, key.GenPrimaryCachePrefix(tableName))
	if err != nil {
		c.Logger.CtxDebug(ctx, "[InvalidateAllPrimaryCache] invalidating all primary cache for table %s error: %v",
			tableName, err)
	}
}

func (c *Gorm2Cache) BatchPrimaryKeyExists(ctx context.Context, tableName string, primaryKeys []string) bool {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, key.GenPrimaryCacheKey(tableName, primaryKey))
	}
	allExists, err := c.dataLayer.BatchKeyExists(ctx, cacheKeys)
	if err != nil {
		c.Logger.CtxDebug(ctx, "[BatchPrimaryKeyExists] checking batch primary keys %v exists error: %v",
			primaryKeys, err)
		return false
	}
	return allExists
}

func (c *Gorm2Cache) SQLKeyExists(ctx context.Context, tableName string, SQL string, vars ...interface{}) bool {
	cacheKey := key.GenSearchCacheKey(tableName, SQL, vars...)
	exists, err := c.dataLayer.KeyExists(ctx, cacheKey)
	if err != nil {
		c.Logger.CtxDebug(ctx, "[SQLKeyExists] checking SQL key %s exists error: %v", cacheKey, err)
		return false
	}
	return exists
}

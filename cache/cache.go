package cache

import (
	"context"

	"github.com/Pacific73/gorm-cache/callback"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/Pacific73/gorm-cache/data_layer"
	"gorm.io/gorm"
)

type Gorm2Cache struct {
	Config    *config.CacheConfig
	Db        *gorm.DB
	DataLayer data_layer.DataLayerInterface
	Logger    config.LoggerInterface
}

func (c *Gorm2Cache) AttachToDB() {
	if c.Config.InvalidateWhenUpdate {
		c.Db.Callback().Create().After("gorm:after_create").Register("gorm:cache:after_create", callback.AfterCreate(c))

		c.Db.Callback().Delete().After("gorm:after_delete").Register("gorm:cache:after_delete", callback.AfterDelete(c))

		c.Db.Callback().Update().After("gorm:after_update").Register("gorm:cache:after_update", callback.AfterUpdate(c))
	}

	c.Db.Callback().Query().Before("gorm:query").Register("gorm:cache:before_query", callback.BeforeQuery(c))
	c.Db.Callback().Query().After("gorm:after_query").Register("gorm:cache:after_query", callback.AfterQuery(c))

	c.Db.Callback().Row().Before("gorm:row").Register("gorm:cache:before_row_query", callback.BeforeRow(c))
	c.Db.Callback().Row().After("gorm:row").Register("gorm:cache:after_row_query", callback.AfterRow(c))
}

func (c *Gorm2Cache) InvalidateSearchCache(ctx context.Context, tableName string) {
	err := c.DataLayer.DeleteKeysWithPrefix(ctx, GenSearchCachePrefix(tableName))
	if err != nil {
		c.Logger.CtxDebug(ctx, "[InvalidateSearchCache] invalidating search cache for table %s error: %+v",
			tableName, err)
	}
}

func (c *Gorm2Cache) InvalidatePrimaryCache(ctx context.Context, tableName string, primaryKey string) {
	err := c.DataLayer.DeleteKey(ctx, GenPrimaryCacheKey(tableName, primaryKey))
	if err != nil {
		c.Logger.CtxDebug(ctx, "[InvalidatePrimaryCache] invalidating primary cache for key %s:%s error: %+v",
			tableName, primaryKey, err)
	}
}

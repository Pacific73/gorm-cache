package cache

import (
	"github.com/Pacific73/gorm-cache/config"
	"github.com/Pacific73/gorm-cache/util"
	"gorm.io/gorm"
)

func AfterCreate(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if cache.Config.InvalidateWhenUpdate && util.ShouldCache(tableName, cache.Config.Tables) {
			if cache.Config.CacheLevel == config.CacheLevelAll || cache.Config.CacheLevel == config.CacheLevelOnlySearch {
				// We invalidate search cache here,
				// because any newly created objects may cause search cache results to be outdated and invalid.
				cache.Logger.CtxInfo(ctx, "[AfterCreate] now start to invalidate search cache for table: %s", tableName)
				err := cache.InvalidateSearchCache(ctx, tableName)
				if err != nil {
					cache.Logger.CtxError(ctx, "[AfterCreate] invalidating search cache for table %s error: %v",
						tableName, err)
					return
				}
				cache.Logger.CtxInfo(ctx, "[AfterCreate] invalidating search cache for table: %s finished.", tableName)
			}
		}
	}
}

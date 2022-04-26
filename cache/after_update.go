package cache

import (
	"sync"

	"github.com/Pacific73/gorm-cache/util"

	"github.com/Pacific73/gorm-cache/config"

	"gorm.io/gorm"
)

func AfterUpdate(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Table
		ctx := db.Statement.Context
		cache.Logger.CtxInfo(ctx, "here")

		if db.Error == nil && cache.Config.InvalidateWhenUpdate && util.ShouldCache(tableName, cache.Config.Tables) {
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()

				if cache.Config.CacheLevel == config.CacheLevelAll || cache.Config.CacheLevel == config.CacheLevelOnlyPrimary {
					primaryKeys := getPrimaryKeysFromWhereClause(db)
					cache.Logger.CtxInfo(ctx, "[AfterUpdate] parse primary keys = %v", primaryKeys)

					if len(primaryKeys) > 0 {
						cache.Logger.CtxInfo(ctx, "[AfterUpdate] now start to invalidate cache for primary keys: %+v",
							primaryKeys)
						err := cache.BatchInvalidatePrimaryCache(ctx, tableName, primaryKeys)
						if err != nil {
							cache.Logger.CtxError(ctx, "[AfterUpdate] invalidating primary cache for key %v error: %v",
								primaryKeys, err)
							return
						}
						cache.Logger.CtxInfo(ctx, "[AfterUpdate] invalidating cache for primary keys: %+v finished.", primaryKeys)
					} else {
						cache.Logger.CtxInfo(ctx, "[AfterUpdate] now start to invalidate all primary cache for table: %s", tableName)
						err := cache.InvalidateAllPrimaryCache(ctx, tableName)
						if err != nil {
							cache.Logger.CtxError(ctx, "[AfterUpdate] invalidating primary cache for table %s error: %v",
								tableName, err)
							return
						}
						cache.Logger.CtxInfo(ctx, "[AfterUpdate] invalidating all primary cache for table: %s finished.", tableName)
					}
				}
			}()

			go func() {
				defer wg.Done()

				if cache.Config.CacheLevel == config.CacheLevelAll || cache.Config.CacheLevel == config.CacheLevelOnlySearch {
					cache.Logger.CtxInfo(ctx, "[AfterUpdate] now start to invalidate search cache for table: %s", tableName)
					err := cache.InvalidateSearchCache(ctx, tableName)
					if err != nil {
						cache.Logger.CtxError(ctx, "[AfterUpdate] invalidating search cache for table %s error: %v",
							tableName, err)
						return
					}
					cache.Logger.CtxInfo(ctx, "[AfterUpdate] invalidating search cache for table: %s finished.", tableName)
				}
			}()

			wg.Wait()
		}
	}
}

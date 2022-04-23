package callback

import (
	"sync"

	"github.com/Pacific73/gorm-cache/config"

	"github.com/Pacific73/gorm-cache/cache"
	"gorm.io/gorm"
)

func AfterDelete(cache *cache.Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Table
		ctx := db.Statement.Context

		if cache.Config.InvalidateWhenUpdate && ContainString(tableName, cache.Config.Tables) {
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()

				if cache.Config.CacheLevel == config.CacheLevelAll || cache.Config.CacheLevel == config.CacheLevelOnlyPrimary {
					primaryKeys := getPrimaryKeysFromWhereClause(db)
					cache.Logger.CtxDebug(ctx, "[AfterDelete] now start to invalidate cache for primary keys: %+v",
						primaryKeys)
					cache.BatchInvalidatePrimaryCache(ctx, tableName, primaryKeys)
					cache.Logger.CtxDebug(ctx, "[AfterDelete] invalidating cache for primary keys: %+v finished.", primaryKeys)
				}
			}()

			go func() {
				defer wg.Done()

				if cache.Config.CacheLevel == config.CacheLevelAll || cache.Config.CacheLevel == config.CacheLevelOnlySearch {
					cache.Logger.CtxDebug(ctx, "[AfterDelete] now start to invalidate search cache for table: %s", tableName)
					cache.InvalidateSearchCache(ctx, tableName)
					cache.Logger.CtxDebug(ctx, "[AfterDelete] invalidating search cache for table: %s finished.", tableName)
				}
			}()

			wg.Wait()
		}
	}
}

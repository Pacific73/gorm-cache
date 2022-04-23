package callback

import (
	"sync"

	"github.com/Pacific73/gorm-cache/cache"
	"gorm.io/gorm"
)

func AfterCreate(cache *cache.Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table

		var outerWg sync.WaitGroup
		outerWg.Add(2)

		go func() {
			defer outerWg.Done()

			// We invalidate related cached objects here, not caching newly created objects.
			// Because gorm only back fill object's id without other possible default values set in database side,
			// which makes our cached objects (if any) inconsistent.
			primaryKeys := GetPrimaryKeysAfterCreate(db)
			cache.Logger.CtxDebug(db.Statement.Context,
				"[AfterCreate] now start to invalidate cache for primary keys: %+v", primaryKeys)
			var wg sync.WaitGroup
			for _, p := range primaryKeys {
				wg.Add(1)
				go func(primaryKey string) {
					defer wg.Done()
					cache.InvalidatePrimaryCache(tableName, primaryKey)
				}(p)
			}
			cache.Logger.CtxDebug(db.Statement.Context,
				"[AfterCreate] invalidating cache for primary keys: %+v finished.", primaryKeys)
			wg.Wait()
		}()

		go func() {
			defer outerWg.Done()
			// We invalidate search cache here,
			// because any newly created objects may cause search cache results to be outdated and invalid.
			cache.Logger.CtxDebug(db.Statement.Context,
				"[AfterCreate] now start to invalidate search cache for table: %s", tableName)
			cache.InvalidateSearchCache(tableName)
			cache.Logger.CtxDebug(db.Statement.Context,
				"[AfterCreate] invalidating search cache for table: %s finished.", tableName)
		}()

		outerWg.Wait()

	}
}

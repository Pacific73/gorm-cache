package callback

import (
	"github.com/Pacific73/gorm-cache/cache"
	"gorm.io/gorm"
)

func AfterCreate(cache *cache.Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table

		if cache.Config.InvalidateWhenUpdate {
			// We invalidate search cache here,
			// because any newly created objects may cause search cache results to be outdated and invalid.
			cache.Logger.CtxDebug(db.Statement.Context,
				"[AfterCreate] now start to invalidate search cache for table: %s", tableName)
			cache.InvalidateSearchCache(db.Statement.Context, tableName)
			cache.Logger.CtxDebug(db.Statement.Context,
				"[AfterCreate] invalidating search cache for table: %s finished.", tableName)
		}
	}
}

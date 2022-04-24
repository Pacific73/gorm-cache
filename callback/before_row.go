package callback

import (
	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/util"
	"gorm.io/gorm"
)

func BeforeRow(cache *cache.Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if util.ShouldCache(tableName, cache.Config.Tables) {

			keyExists := cache.SearchKeyExists(ctx, tableName, db.Statement.SQL.String(), db.Statement.Vars...)
			if keyExists {
				db.Error = util.SearchCacheHit
				return
			}
		}
	}
}

package cache

import (
	"github.com/Pacific73/gorm-cache/util"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func BeforeRow(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		callbacks.BuildQuerySQL(db)
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if util.ShouldCache(tableName, cache.Config.Tables) {
			sql := db.Statement.SQL.String()
			keyExists, err := cache.SearchKeyExists(ctx, tableName, sql, db.Statement.Vars...)
			if err != nil {
				cache.Logger.CtxError(ctx, "[BeforeRow] check key exists for sql %s error: %v", sql, err)
				return
			}
			if keyExists {
				db.Error = util.SearchCacheHit
				return
			}
		}
	}
}

package cache

import (
	"github.com/Pacific73/gorm-cache/util"
	"gorm.io/gorm"
)

func BeforeQuery(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if util.ShouldCache(tableName, cache.Config.Tables) {

			keyExists, err := cache.SearchKeyExists(ctx, tableName, db.Statement.SQL.String(), db.Statement.Vars...)
			if err != nil {
				cache.Logger.CtxError(ctx, "[BeforeQuery] check search key exists for key %s error: %v",
					db.Statement.SQL.String(), err)
				return
			}
			if keyExists {
				db.Error = util.SearchCacheHit
				return
			}

			primaryKeys := getPrimaryKeysFromWhereClause(db)

			// if (IN primaryKeys)/(Eq primaryKey) are the only clauses
			hasOtherClauseInWhere := hasOtherClauseExceptPrimaryField(db)
			if hasOtherClauseInWhere {
				// if query has other clauses, it can only query the database
				return
			}

			allKeyExist, err := cache.BatchPrimaryKeyExists(ctx, tableName, primaryKeys)
			if err != nil {
				cache.Logger.CtxError(ctx, "[BeforeQuery] check primary key exists for key %v error: %v", primaryKeys, err)
				return
			}
			if allKeyExist {
				db.InstanceSet("gorm:cache:primary_keys", primaryKeys)
				db.Error = util.PrimaryCacheHit
				// if part or none of the objects are cached, query the database
				return
			}
		}
	}
}

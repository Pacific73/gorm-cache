package callback

import (
	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/util"
	"gorm.io/gorm"
)

func BeforeQuery(cache *cache.Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if util.ShouldCache(tableName, cache.Config.Tables) {

			keyExists := cache.SearchKeyExists(ctx, tableName, db.Statement.SQL.String(), db.Statement.Vars...)
			if keyExists {
				db.AddError(util.SearchCacheHit)
				return
			}

			primaryKeys := getPrimaryKeysFromWhereClause(db)

			// if (IN primaryKeys)/(Eq primaryKey) are the only clauses
			hasOtherClauseInWhere := hasOtherClauseExceptPrimaryField(db)
			if hasOtherClauseInWhere {
				// if query has other clauses, it can only query the database
				return
			}

			allKeyExist := cache.BatchPrimaryKeyExists(ctx, tableName, primaryKeys)
			if allKeyExist {
				db.InstanceSet("gorm:cache:primary_keys", primaryKeys)
				db.AddError(util.PrimaryCacheHit)
				// if part or none of the objects are cached, query the database
				return
			}
		}
	}
}

package callback

import (
	"github.com/Pacific73/gorm-cache/cache"
	"gorm.io/gorm"
)

func BeforeQuery(cache *cache.Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if ContainString(tableName, cache.Config.Tables) {
			primaryKeys := getPrimaryKeysFromWhereClause(db)

			// if (IN primaryKeys)/(Eq primaryKey) are the only clauses
			hasOtherClauseInWhere := hasOtherClauseExceptPrimaryField(db)
			if hasOtherClauseInWhere {
				// if query has other clauses, it can only query the database
				return
			}

			keyExists := cache.SQLKeyExists(ctx, tableName, db.Statement.SQL.String(), db.Statement.Vars...)
			if keyExists {
				db.AddError(searchCacheHit)
				return
			}

			allKeyExist := cache.BatchPrimaryKeyExists(ctx, tableName, primaryKeys)
			if allKeyExist {
				db.InstanceSet("primary_keys", primaryKeys)
				db.AddError(primaryCacheHit)
				// if part or none of the objects are cached, query the database
				return
			}
		}
	}
}

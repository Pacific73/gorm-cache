package cache

import (
	"github.com/Pacific73/gorm-cache/config"
	"github.com/Pacific73/gorm-cache/util"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func BeforeQuery(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		callbacks.BuildQuerySQL(db)
		tableName := ""
		if db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		} else {
			tableName = db.Statement.Table
		}
		ctx := db.Statement.Context

		sql := db.Statement.SQL.String()
		db.InstanceSet("gorm:cache:sql", sql)
		db.InstanceSet("gorm:cache:vars", db.Statement.Vars)

		if util.ShouldCache(tableName, cache.Config.Tables) {

			if cache.Config.CacheLevel == config.CacheLevelAll || cache.Config.CacheLevel == config.CacheLevelOnlySearch {
				keyExists, err := cache.SearchKeyExists(ctx, tableName, sql, db.Statement.Vars...)
				if err != nil {
					cache.Logger.CtxError(ctx, "[BeforeQuery] check search key exists for key %s error: %v",
						sql, err)
					return
				}
				cache.Logger.CtxInfo(ctx, "[BeforeQuery] search key exists ? %v", keyExists)
				if keyExists {
					db.Error = util.SearchCacheHit
					return
				}
			}

			if cache.Config.CacheLevel == config.CacheLevelAll || cache.Config.CacheLevel == config.CacheLevelOnlyPrimary {
				primaryKeys := getPrimaryKeysFromWhereClause(db)
				cache.Logger.CtxInfo(ctx, "[BeforeQuery] parse primary keys = %v", primaryKeys)

				if len(primaryKeys) == 0 {
					return
				}

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
				cache.Logger.CtxInfo(ctx, "[BeforeQuery] all primary key exists ? %v", allKeyExist)
				if allKeyExist {
					db.InstanceSet("gorm:cache:primary_keys", primaryKeys)
					db.Error = util.PrimaryCacheHit
					// if part or none of the objects are cached, query the database
					return
				}
			}
		}
	}
}

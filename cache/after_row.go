package cache

import (
	"encoding/json"
	"errors"

	"github.com/Pacific73/gorm-cache/util"
	"gorm.io/gorm"
)

func AfterRow(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if db.Error == nil {
			// error is nil -> cache not hit, we cache newly retrieved data
			_, objects := GetObjectsAfterLoad(db)

			go func() {
				// cache search data
				sql := db.Statement.SQL.String()
				cache.Logger.CtxInfo(ctx, "[AfterRow] start to set search cache for sql: %s", sql)
				cacheBytes, err := json.Marshal(objects)
				if err != nil {
					cache.Logger.CtxError(ctx, "[AfterRow] cannot marshal cache for sql: %s, not cached", sql)
					return
				}
				err = cache.SetSearchCache(ctx, string(cacheBytes), tableName, sql, db.Statement.Vars...)
				if err != nil {
					cache.Logger.CtxError(ctx, "[AfterRow] set search cache for sql %s error", sql, err)
					return
				}
				cache.Logger.CtxInfo(ctx, "[AfterRow] sql %s cached", sql)
			}()

			return
		}

		if errors.Is(db.Error, util.SearchCacheHit) {
			// search cache hit
			cacheValue, err := cache.GetSearchCache(ctx, tableName, db.Statement.SQL.String(), db.Statement.Vars...)
			if err != nil {

			}
			err = json.Unmarshal([]byte(cacheValue), db.Statement.Dest)
			if err != nil {
				cache.Logger.CtxError(ctx, "[AfterRow] unmarshal search cache error")
				db.Error = util.ErrCacheUnmarshal
				return
			}
			cache.IncrHitCount()
			db.Error = nil
			return
		}
	}
}

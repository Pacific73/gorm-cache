package callback

import (
	"encoding/json"
	"errors"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/util"
	"gorm.io/gorm"
)

func AfterRow(cache *cache.Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if db.Error == nil {
			// error is nil -> cache not hit, we cache newly retrieved data
			_, objects := GetObjectsAfterLoad(db)

			go func() {
				// cache search data
				cache.Logger.CtxInfo(ctx, "[AfterQuery] start to set search cache for sql: %s", db.Statement.SQL.String())
				cacheBytes, err := json.Marshal(objects)
				if err != nil {
					cache.Logger.CtxError(ctx, "[AfterQuery] cannot marshal cache for sql: %s, not cached", db.Statement.SQL.String())
					return
				}
				cache.SetSearchCache(ctx, string(cacheBytes), tableName, db.Statement.SQL.String(), db.Statement.Vars...)
				cache.Logger.CtxInfo(ctx, "[AfterQuery] sql %s cached", db.Statement.SQL.String())
			}()

			return
		}

		if errors.Is(db.Error, util.SearchCacheHit) {
			// search cache hit
			cacheValue := cache.GetSearchCache(ctx, tableName, db.Statement.SQL.String(), db.Statement.Vars...)
			err := json.Unmarshal([]byte(cacheValue), db.Statement.Dest)
			if err != nil {
				cache.Logger.CtxError(ctx, "[AfterQuery] unmarshal search cache error")
				db.Error = util.ErrCacheUnmarshal
				return
			}
			db.Error = nil
			return
		}
	}
}

package callback

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Pacific73/gorm-cache/util"

	"github.com/Pacific73/gorm-cache/cache"
	"gorm.io/gorm"
)

func AfterQuery(cache *cache.Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if db.Error == nil {
			// error is nil -> cache not hit, we cache newly retrieved data
			primaryKeys, objects := GetObjectsAfterLoad(db)

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

			go func() {
				// cache primary cache data
				if len(primaryKeys) != len(objects) {
					return
				}
				kvs := make([]util.Kv, 0, len(objects))
				for i := 0; i < len(objects); i++ {
					jsonStr, err := json.Marshal(objects[i])
					if err != nil {
						cache.Logger.CtxError(ctx, "[AfterQuery] object %v cannot marshal, not cached")
						continue
					}
					kvs = append(kvs, util.Kv{
						Key:   primaryKeys[i],
						Value: string(jsonStr),
					})
				}
				cache.BatchSetPrimaryKeyCache(ctx, tableName, kvs)
			}()
			return
		} else if errors.Is(db.Error, util.SearchCacheHit) {
			cacheValue := cache.GetSearchCache(ctx, tableName, db.Statement.SQL.String(), db.Statement.Vars...)
			err := json.Unmarshal([]byte(cacheValue), db.Statement.Dest)
			if err != nil {
				cache.Logger.CtxError(ctx, "[AfterQuery] unmarshal search cache error")
				db.AddError(util.ErrCacheUnmarshal)
				return
			}
			return
		} else if errors.Is(db.Error, util.PrimaryCacheHit) {
			primaryKeyObjs, ok := db.InstanceGet("gorm:cache:primary_keys")
			if !ok {
				cache.Logger.CtxError(ctx, "[AfterQuery] cannot get primary keys from db instance get")
				db.AddError(util.ErrCacheUnmarshal)
				return
			}
			primaryKeys := primaryKeyObjs.([]string)
			cacheValues := cache.BatchGetPrimaryCache(ctx, tableName, primaryKeys)
			finalValue := "[" + strings.Join(cacheValues, ",") + "]"
			err := json.Unmarshal([]byte(finalValue), db.Statement.Dest)
			if err != nil {
				cache.Logger.CtxError(ctx, "[AfterQuery] unmarshal final value error")
				db.AddError(util.ErrCacheUnmarshal)
				return
			}
			return
		}
	}
}

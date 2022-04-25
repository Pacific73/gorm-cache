package cache

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Pacific73/gorm-cache/util"

	"gorm.io/gorm"
)

func AfterQuery(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := db.Statement.Schema.Table
		ctx := db.Statement.Context

		if db.Error == nil {
			// error is nil -> cache not hit, we cache newly retrieved data
			primaryKeys, objects := GetObjectsAfterLoad(db)

			go func() {
				// cache search data
				sql := db.Statement.SQL.String()
				cache.Logger.CtxInfo(ctx, "[AfterQuery] start to set search cache for sql: %s", sql)
				cacheBytes, err := json.Marshal(objects)
				if err != nil {
					cache.Logger.CtxError(ctx, "[AfterQuery] cannot marshal cache for sql: %s, not cached", sql)
					return
				}
				err = cache.SetSearchCache(ctx, string(cacheBytes), tableName, sql, db.Statement.Vars...)
				if err != nil {
					cache.Logger.CtxError(ctx, "[AfterQuery] set search cache for sql: %s error: %v", sql, err)
					return
				}
				cache.Logger.CtxInfo(ctx, "[AfterQuery] sql %s cached", sql)
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
				err := cache.BatchSetPrimaryKeyCache(ctx, tableName, kvs)
				if err != nil {
					cache.Logger.CtxError(ctx, "[AfterQuery] batch set primary key cache for key %v error: %v",
						primaryKeys, err)
				}
			}()
			return
		}

		if errors.Is(db.Error, util.SearchCacheHit) {
			// search cache hit
			cacheValue, err := cache.GetSearchCache(ctx, tableName, db.Statement.SQL.String(), db.Statement.Vars...)
			if err != nil {
				cache.Logger.CtxError(ctx, "[AfterQuery] get cache value for sql %s error: %v", db.Statement.SQL.String(), err)
				db.Error = util.ErrCacheLoadFailed
				return
			}
			err = json.Unmarshal([]byte(cacheValue), db.Statement.Dest)
			if err != nil {
				cache.Logger.CtxError(ctx, "[AfterQuery] unmarshal search cache error")
				db.Error = util.ErrCacheUnmarshal
				return
			}
			cache.IncrHitCount()
			db.Error = nil
			return
		}

		if errors.Is(db.Error, util.PrimaryCacheHit) {
			// primary cache hit
			primaryKeyObjs, ok := db.InstanceGet("gorm:cache:primary_keys")
			if !ok {
				cache.Logger.CtxError(ctx, "[AfterQuery] cannot get primary keys from db instance get")
				db.Error = util.ErrCacheUnmarshal
				return
			}
			primaryKeys := primaryKeyObjs.([]string)
			cacheValues, err := cache.BatchGetPrimaryCache(ctx, tableName, primaryKeys)
			if err != nil {
				cache.Logger.CtxError(ctx, "[AfterQuery] get primary cache value for key %v error %v", primaryKeys, err)
				db.Error = util.ErrCacheLoadFailed
				return
			}
			finalValue := "[" + strings.Join(cacheValues, ",") + "]"
			err = json.Unmarshal([]byte(finalValue), db.Statement.Dest)
			if err != nil {
				cache.Logger.CtxError(ctx, "[AfterQuery] unmarshal final value error")
				db.Error = util.ErrCacheUnmarshal
				return
			}
			cache.IncrHitCount()
			db.Error = nil
			return
		}
	}
}

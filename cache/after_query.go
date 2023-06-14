package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/Pacific73/gorm-cache/config"
	"github.com/Pacific73/gorm-cache/util"
	"gorm.io/gorm"
)

func AfterQuery(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tableName := ""
		if db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		} else {
			tableName = db.Statement.Table
		}
		ctx := db.Statement.Context
		sqlObj, _ := db.InstanceGet("gorm:cache:sql")
		sql := sqlObj.(string)
		varObj, _ := db.InstanceGet("gorm:cache:vars")
		vars := varObj.([]interface{})

		if util.ShouldCache(tableName, cache.Config.Tables) {
			if db.Error == nil {
				// error is nil -> cache not hit, we cache newly retrieved data
				primaryKeys, objects := getObjectsAfterLoad(db)

				var wg sync.WaitGroup
				wg.Add(2)

				go func() {
					defer wg.Done()

					if cache.Config.CacheLevel == config.CacheLevelAll || cache.Config.CacheLevel == config.CacheLevelOnlySearch {
						// cache search data
						if cache.Config.CacheMaxItemCnt > 0 && int64(len(objects)) > cache.Config.CacheMaxItemCnt {
							return
						}

						cache.Logger.CtxInfo(ctx, "[AfterQuery] start to set search cache for sql: %s", sql)
						cacheBytes, err := json.Marshal(db.Statement.Dest)
						if err != nil {
							cache.Logger.CtxError(ctx, "[AfterQuery] cannot marshal cache for sql: %s, not cached", sql)
							return
						}
						cache.Logger.CtxInfo(ctx, "[AfterQuery] set cache: %v", string(cacheBytes))
						err = cache.SetSearchCache(ctx, fmt.Sprintf("%d|", db.RowsAffected)+string(cacheBytes), tableName, sql, vars...)
						if err != nil {
							cache.Logger.CtxError(ctx, "[AfterQuery] set search cache for sql: %s error: %v", sql, err)
							return
						}
						cache.Logger.CtxInfo(ctx, "[AfterQuery] sql %s cached", sql)
					}
				}()

				go func() {
					defer wg.Done()

					if cache.Config.CacheLevel == config.CacheLevelAll || cache.Config.CacheLevel == config.CacheLevelOnlyPrimary {
						// cache primary cache data
						if len(primaryKeys) != len(objects) {
							return
						}
						if cache.Config.CacheMaxItemCnt > 0 && int64(len(objects)) > cache.Config.CacheMaxItemCnt {
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
						cache.Logger.CtxInfo(ctx, "[AfterQuery] start to set primary cache for kvs: %+v", kvs)
						err := cache.BatchSetPrimaryKeyCache(ctx, tableName, kvs)
						if err != nil {
							cache.Logger.CtxError(ctx, "[AfterQuery] batch set primary key cache for key %v error: %v",
								primaryKeys, err)
						}
					}
				}()
				wg.Wait()
				return
			}
		}

		if errors.Is(db.Error, util.SearchCacheHit) {
			// search cache hit
			db.Error = nil
			return
		}

		if errors.Is(db.Error, util.PrimaryCacheHit) {
			// primary cache hit
			db.Error = nil
			return
		}
	}
}

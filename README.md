[![GoVersion](https://img.shields.io/github/go-mod/go-version/Pacific73/gorm-cache)](https://github.com/Pacific73/gorm-cache/blob/master/go.mod)
[![Release](https://img.shields.io/github/v/release/Pacific73/gorm-cache)](https://github.com/Pacific73/gorm-cache/releases)
[![Apache-2.0 license](https://img.shields.io/badge/license-Apache2.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)

[English Version](./README.md) | [中文版本](./README.ZH_CN.md)

# gorm-cache

`gorm-cache` aims to provide a look-aside, almost-no-code-modification cache solution for gorm v2 users. It only applys to situations where database table has only one single primary key.

We provide 2 types of cache storage here:

1. Memory, where all cached data stores in memory of a single server
2. Redis, where cached data stores in Redis (if you have multiple servers running the same procedure, they don't share the same space in Redis)

## Usage

```go
import (
    "context"
    "github.com/Pacific73/gorm-cache/cache"
    "github.com/go-redis/redis"
)

func main() {
    dsn := "user:pass@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4"
    db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",    
    })
    
    cache, _ := cache.NewGorm2Cache(&config.CacheConfig{
        CacheLevel:           config.CacheLevelAll,
        CacheStorage:         config.CacheStorageRedis,
        RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
        InvalidateWhenUpdate: true, // when you create/update/delete objects, invalidate cache
        CacheTTL:             5000, // 5000 ms
        CacheMaxItemCnt:      5,    // if length of objects retrieved one single time 
                                    // exceeds this number, then don't cache
    })
    // More options in `config.config.go`
    db.Use(cache)    // use gorm plugin
    // cache.AttachToDB(db)

    var users []User
    
    db.Where("value > ?", 123).Find(&users) // search cache not hit, objects cached
    db.Where("value > ?", 123).Find(&users) // search cache hit
    
    db.Where("id IN (?)", []int{1, 2, 3}).Find(&users) // primary key cache not hit, users cached
    db.Where("id IN (?)", []int{1, 3}).Find(&users) // primary key cache hit
}
```

There're mainly 5 kinds of operations in gorm (gorm function names in brackets):

1. Query (First/Take/Last/Find/FindInBatches/FirstOrInit/FirstOrCreate/Count/Pluck)
2. Create (Create/CreateInBatches/Save)
3. Delete (Delete)
4. Update (Update/Updates/UpdateColumn/UpdateColumns/Save)
5. Row (Row/Rows/Scan)

We don't support caching in Row operations.

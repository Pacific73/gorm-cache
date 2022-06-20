[![GoVersion](https://img.shields.io/github/go-mod/go-version/Pacific73/gorm-cache)](https://github.com/Pacific73/gorm-cache/blob/master/go.mod)
[![Release](https://img.shields.io/github/v/release/Pacific73/gorm-cache)](https://github.com/Pacific73/gorm-cache/releases)
[![Apache-2.0 license](https://img.shields.io/badge/license-Apache2.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)

[English Version](./README.md) | [中文版本](./README.ZH_CN.md)

`gorm-cache` 旨在为gorm v2用户提供一个即插即用的旁路缓存解决方案。本缓存只适用于数据库表单主键时的场景。

我们提供2种存储介质：

1. 内存 (所有数据存储在单服务器的内存中)
2. Redis (所有数据存储在redis中，如果你有多个实例使用本缓存，那么他们不共享redis存储空间)

## 使用说明

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
    db.Plugin(cache)    // use gorm plugin
    // cache.AttachToDB(db)

    var users []User
    
    db.Where("value > ?", 123).Find(&users) // search cache not hit, objects cached
    db.Where("value > ?", 123).Find(&users) // search cache hit
    
    db.Where("id IN (?)", []int{1, 2, 3}).Find(&users) // primary key cache not hit, users cached
    db.Where("id IN (?)", []int{1, 3}).Find(&users) // primary key cache hit
}
```

在gorm中主要有5种操作（括号中是gorm中对应函数名）:

1. Query (First/Take/Last/Find/FindInBatches/FirstOrInit/FirstOrCreate/Count/Pluck)
2. Create (Create/CreateInBatches/Save)
3. Delete (Delete)
4. Update (Update/Updates/UpdateColumn/UpdateColumns/Save)
5. Row (Row/Rows/Scan)

我们不支持Row操作的缓存。

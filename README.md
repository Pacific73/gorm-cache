# gorm-cache (working in progress, don't use)

`gorm-cache` aims to provide a look-aside, almost-no-code-modification cache solution to gorm v2 users. It only applys to situations where database table has only one single primary key.

We provide 2 types of cache storage here:

1. Memory, where all cached data stores in memory of a single server
2. Redis, where cached data stores in redis (if you have multiple servers running the same procedure, they don't share the same redis storage space)

`gorm-cache` 旨在为gorm v2用户提供一个即插即用的旁路缓存解决方案。本缓存只适用于数据库表单主键时的场景。

我们提供2种存储介质：

1. 内存 (所有数据存储在单服务器的内存中)
2. Redis (所有数据存储在redis中，如果你有多个实例使用本缓存，那么他们不共享redis存储空间)

## Usage 使用说明

```go
import (
    "github.com/Pacific73/gorm-cache/cache"
    "github.com/go-redis/redis"
)

func main() {
    dsn = "user:pass@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4"
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
        CacheMaxItemCnt:      5, // if length of objects retrieved one single time 
                                 // exceeds this number, then don't cache
        PrimaryCacheSize:     200, // 200 MB
        SearchCacheSize:      200,
    })
    // More options in `config.config.go`
    
    cache.AttachToDB(db)

    var users []User
    
    db.Where("id = ?", 1).Find(&users) // cache not hit, users cached
    db.Where("id = ?", 1).Find(&users) // cache hit
    
    type queryStruct struct {
        ID    int
        Value int
    }
    var dests []queryStruct
    db.Select([]string{"a.id AS id", "b.value AS value"}).
        Table("a").
        Joins("JOIN b ON b.id = a.id").
        Where("b.value > ?", 0).Scan(&dests) // cache not hit, query result cached
    
    db.Select([]string{"a.id AS id", "b.value AS value"}).
    	Table("a").
    	Joins("JOIN b ON b.id = a.id").
    	Where("b.value > ?", 0).Scan(&dests) // cache hit
}
```

## Readme English Version (todo...)

There're mainly 5 kinds of operations in gorm (gorm function names in brackets):

1. Query (First/Take/Last/Find/FindInBatches/FirstOrInit/FirstOrCreate/Count/Pluck)
2. Create (Create/CreateInBatches/Save)
3. Delete (Delete)
4. Update (Update/Updates/UpdateColumn/UpdateColumns/Save)
5. Row (Row/Rows/Scan)

## Readme 中文版 (待补充...)

在gorm中主要有5种操作（括号中是gorm中对应函数名）：

1. Query (First/Take/Last/Find/FindInBatches/FirstOrInit/FirstOrCreate/Count/Pluck)
2. Create (Create/CreateInBatches/Save)
3. Delete (Delete)
4. Update (Update/Updates/UpdateColumn/UpdateColumns/Save)
5. Row (Row/Rows/Scan)



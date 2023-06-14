package testkit

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm/logger"

	"github.com/Pacific73/gorm-cache/cache"

	"github.com/Pacific73/gorm-cache/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	username     = "root"
	password     = "Zcydf741205,."
	databaseName = "site_reldb"
	ip           = "localhost"
	port         = "3306"
)

var (
	redisIp   = "localhost"
	redisPort = "6379"
)

var (
	searchCache  *cache.Gorm2Cache
	primaryCache *cache.Gorm2Cache
	allCache     *cache.Gorm2Cache

	searchDB   *gorm.DB
	primaryDB  *gorm.DB
	allDB      *gorm.DB
	originalDB *gorm.DB
)

var (
	testSize = 200 // minimum 200
)

func TestMain(m *testing.M) {
	log("test setup ...")

	var err error
	//logger.Default.LogMode(logger.Info)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, ip, port, databaseName)
	originalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          logger.Default,
	})
	if err != nil {
		log("open db error: %v", err)
		os.Exit(-1)
	}

	searchDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          logger.Default,
	})
	if err != nil {
		log("open db error: %v", err)
		os.Exit(-1)
	}

	primaryDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          logger.Default,
	})
	if err != nil {
		log("open db error: %v", err)
		os.Exit(-1)
	}

	allDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          logger.Default,
	})
	if err != nil {
		log("open db error: %v", err)
		os.Exit(-1)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: redisIp + ":" + redisPort})

	searchCache, err = cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelOnlySearch,
		CacheStorage:         config.CacheStorageMemory,
		RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
		InvalidateWhenUpdate: true,
		CacheTTL:             5000,
		CacheMaxItemCnt:      5000,
		CacheSize:            1000,
		DebugMode:            false,
	})
	if err != nil {
		log("setup search cache error: %v", err)
		os.Exit(-1)
	}

	primaryCache, err = cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelOnlyPrimary,
		CacheStorage:         config.CacheStorageMemory,
		RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
		InvalidateWhenUpdate: true,
		CacheTTL:             5000,
		CacheMaxItemCnt:      5000,
		CacheSize:            1000,
		DebugMode:            false,
	})
	if err != nil {
		log("setup primary cache error: %v", err)
		os.Exit(-1)
	}

	allCache, err = cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         config.CacheStorageMemory,
		RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
		InvalidateWhenUpdate: true,
		CacheTTL:             5000,
		CacheMaxItemCnt:      5000,
		CacheSize:            1000,
		DebugMode:            false,
	})
	if err != nil {
		log("setup all cache error: %v", err)
		os.Exit(-1)
	}

	primaryDB.Use(primaryCache)
	searchDB.Use(searchCache)
	allDB.Use(allCache)
	// primaryCache.AttachToDB(primaryDB)
	// searchCache.AttachToDB(searchDB)
	// allCache.AttachToDB(allDB)

	err = timer("prepare table and data", func() error {
		return PrepareTableAndData(originalDB)
	})
	if err != nil {
		log("setup table and data error: %v", err)
		os.Exit(-1)
	}

	result := m.Run()

	err = timer("clean table and data", func() error {
		return CleanTable(originalDB)
	})
	if err != nil {
		log("clean table and data error: %v", err)
		os.Exit(-1)
	}

	log("integration test end.")
	os.Exit(result)
}

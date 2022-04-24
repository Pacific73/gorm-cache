package data_layer

import (
	"context"

	"github.com/Pacific73/gorm-cache/util"

	"github.com/Pacific73/gorm-cache/config"
)

type DataLayerInterface interface {
	Init(config *config.CacheConfig)

	// read
	BatchKeyExists(ctx context.Context, keys []string) bool
	KeyExists(ctx context.Context, key string) bool
	GetValue(ctx context.Context, key string) string

	// write

	CleanCache()

	DeleteKeysWithPrefix(ctx context.Context, keyPrefix string)
	DeleteKey(ctx context.Context, key string)
	BatchDeleteKeys(ctx context.Context, keys []string)
	BatchSetKeys(ctx context.Context, kvs []util.Kv)
	SetKey(ctx context.Context, kv util.Kv)
}

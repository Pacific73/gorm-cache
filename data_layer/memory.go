package data_layer

import (
	"context"
	"fmt"
	"time"

	"github.com/Pacific73/gorm-cache/config"
	"github.com/Pacific73/gorm-cache/util"
	"github.com/karlseguin/ccache/v2"
)

type MemoryLayer struct {
	cache *ccache.Cache
	ttl   int64
}

func (m *MemoryLayer) Init(ctx context.Context, conf *config.CacheConfig, prefix string) error {
	c := ccache.New(ccache.Configure().MaxSize(int64(conf.CacheSize)))
	m.cache = c
	m.ttl = conf.CacheTTL
	return nil
}

func (m *MemoryLayer) CleanCache(ctx context.Context) error {
	m.cache.Clear()
	return nil
}

func (m *MemoryLayer) BatchKeyExist(ctx context.Context, keys []string) (bool, error) {
	for _, key := range keys {
		item := m.cache.Get(key)
		if item == nil || item.Expired() {
			return false, nil
		}
	}
	return true, nil
}

func (m *MemoryLayer) KeyExists(ctx context.Context, key string) (bool, error) {
	item := m.cache.Get(key)
	return item != nil && !item.Expired(), nil
}

func (m *MemoryLayer) GetValue(ctx context.Context, key string) (string, error) {
	item := m.cache.Get(key)
	if item == nil || item.Expired() {
		return "", fmt.Errorf("cannot get item")
	}
	return item.Value().(string), nil
}

func (m *MemoryLayer) BatchGetValues(ctx context.Context, keys []string) ([]string, error) {
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		item := m.cache.Get(key)
		if item != nil && !item.Expired() {
			values = append(values, item.Value().(string))
		}
	}
	if len(values) != len(keys) {
		return nil, fmt.Errorf("cannot get items")
	}
	return values, nil
}

func (m *MemoryLayer) DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error {
	m.cache.DeletePrefix(keyPrefix)
	return nil
}

func (m *MemoryLayer) DeleteKey(ctx context.Context, key string) error {
	m.cache.Delete(key)
	return nil
}

func (m *MemoryLayer) BatchDeleteKeys(ctx context.Context, keys []string) error {
	for _, key := range keys {
		m.cache.Delete(key)
	}
	return nil
}

func (m *MemoryLayer) BatchSetKeys(ctx context.Context, kvs []util.Kv) error {
	for _, kv := range kvs {
		if m.ttl > 0 {
			m.cache.Set(kv.Key, kv.Value, time.Duration(util.RandFloatingInt64(m.ttl))*time.Millisecond)
		} else {
			m.cache.Set(kv.Key, kv.Value, time.Duration(util.RandFloatingInt64(24))*time.Hour)
		}
	}
	return nil
}

func (m *MemoryLayer) SetKey(ctx context.Context, kv util.Kv) error {
	if m.ttl > 0 {
		m.cache.Set(kv.Key, kv.Value, time.Duration(util.RandFloatingInt64(m.ttl))*time.Millisecond)
	} else {
		m.cache.Set(kv.Key, kv.Value, time.Duration(util.RandFloatingInt64(24))*time.Hour)
	}
	return nil
}

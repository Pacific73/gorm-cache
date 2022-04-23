package data_layer

import "context"

type DataLayerInterface interface {
	DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error
	DeleteKey(ctx context.Context, key string) error
	BatchDeleteKeys(ctx context.Context, keys []string) error
	BatchKeyExists(ctx context.Context, keys []string) (bool, error)
	KeyExists(ctx context.Context, key string) (bool, error)
	GetValue(ctx context.Context, key string) (string, error)
}

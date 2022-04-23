package data_layer

import "context"

type DataLayerInterface interface {
	DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error
	DeleteKey(ctx context.Context, key string) error
	BatchDeleteKeys(ctx context.Context, keys []string) error
}

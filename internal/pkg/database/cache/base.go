package cache

import (
	"context"
	"time"
)

type Cacher interface {
	// Get get value from cache, if return nil, nil, then the key doesn't exist or expired
	Get(context context.Context, key string) ([]byte, error)

	// MGet get multiple keys' values from cache
	MGet(context context.Context, keys ...string) ([][]byte, error)

	// Set set key value with ttl, 0 for no ttl
	Set(context context.Context, key string, value []byte, expiration time.Duration) error

	// MSet set multiple key-value pairs, 0 for no ttl
	MSet(context context.Context, entries map[string][]byte, expiration time.Duration) error

	// Del delete key
	Del(context context.Context, keys ...string) error

	// Clear remove all keys
	Clear(context context.Context) error
}

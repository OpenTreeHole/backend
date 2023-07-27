package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/opentreehole/backend/pkg/utils"
)

type RedisCacher struct {
	client *redis.Client
}

func (r RedisCacher) Get(context context.Context, key string) ([]byte, error) {
	value, err := r.client.Get(context, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return utils.StringToBytes(value), err
}

func (r RedisCacher) MGet(context context.Context, keys ...string) ([][]byte, error) {
	values, err := r.client.MGet(context, keys...).Result()
	if err != nil {
		return nil, err
	}
	if len(values) != len(keys) {
		return nil, errors.New("[redis] mget values doesn't match keys in length")
	}
	bytesResults := make([][]byte, len(keys))
	for i := range values {
		switch v := values[i].(type) {
		case string:
			bytesResults[i] = utils.StringToBytes(v)
		case []byte:
			bytesResults[i] = v
		case nil:
			bytesResults[i] = nil
		default:
			return nil, errors.New("[redis] mget unknown types")
		}
	}
	return bytesResults, nil
}

func (r RedisCacher) Set(context context.Context, key string, value []byte, expiration time.Duration) error {
	return r.client.Set(context, key, utils.BytesToString(value), expiration).Err()
}

func (r RedisCacher) MSet(context context.Context, entries map[string][]byte, expiration time.Duration) error {
	if expiration == 0 {
		return r.client.MSet(context, entries).Err()
	}

	var err error
	for key, value := range entries {
		err = errors.Join(err, r.Set(context, key, value, expiration))
	}
	return err
}

func (r RedisCacher) Del(context context.Context, keys ...string) error {
	return r.client.Del(context, keys...).Err()
}

func (r RedisCacher) Clear(context context.Context) error {
	return r.client.FlushDB(context).Err()
}

func NewRedisCacher(opt *redis.Options) Cacher {
	return RedisCacher{client: redis.NewClient(opt)}
}

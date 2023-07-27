package cache

import (
	"context"
	"time"

	"github.com/puzpuzpuz/xsync/v2"
)

type cacheData struct {
	data []byte

	// end time in UnixNano, 0 for no ttl
	endTime int64
}

type MemoryCacher struct {
	data *xsync.MapOf[string, cacheData]
}

func (m MemoryCacher) Get(_ context.Context, key string) ([]byte, error) {
	value, ok := m.data.Load(key)
	if ok {
		if value.endTime != 0 && value.endTime <= time.Now().UnixNano() {
			m.data.Delete(key)
			return nil, nil
		} else {
			return value.data, nil
		}
	} else {
		return nil, nil
	}
}

func (m MemoryCacher) MGet(context context.Context, keys ...string) ([][]byte, error) {
	var values = make([][]byte, len(keys))
	for i, key := range keys {
		value, err := m.Get(context, key)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}

func (m MemoryCacher) Set(_ context.Context, key string, value []byte, expiration time.Duration) error {
	if len(value) == 0 {
		return nil
	}

	var endTime int64
	if expiration != 0 {
		endTime = time.Now().Add(expiration).UnixNano()
	}

	m.data.Store(key, cacheData{value, endTime})
	return nil
}

func (m MemoryCacher) MSet(context context.Context, entries map[string][]byte, expiration time.Duration) error {
	for key, value := range entries {
		err := m.Set(context, key, value, expiration)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m MemoryCacher) Del(_ context.Context, keys ...string) error {
	for _, key := range keys {
		m.data.Delete(key)
	}
	return nil
}

func (m MemoryCacher) Clear(_ context.Context) error {
	m.data.Clear()
	return nil
}

func NewMemoryCacher() Cacher {
	data := xsync.NewMapOf[cacheData]()
	// clear expired data
	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			data.Range(func(key string, value cacheData) bool {
				if value.endTime != 0 && value.endTime <= time.Now().UnixNano() {
					data.Delete(key)
				}
				return true
			})
		}
	}()
	return MemoryCacher{data}
}

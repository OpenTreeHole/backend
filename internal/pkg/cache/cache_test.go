package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/pkg/log"
)

func TestNewCache(t *testing.T) {
	var (
		err     error
		ctx     = context.Background()
		conf    config.Config
		allConf config.AtomicAllConfig
	)
	conf.Cache.Type = "memory"
	allConf.Store(&config.AllConfig{Config: &conf})
	logger, cleanup := log.NewLogger(&allConf)
	defer cleanup()
	var cache = NewCache(&allConf, logger)

	t.Run("set string", func(t *testing.T) {
		err = cache.Set(ctx, "key", "value")
		if err != nil {
			t.Fatal(err)
		}

		var value string
		_, err = cache.Get(ctx, "key", &value)
		if err != nil {
			t.Fatal(err)
		}
		assert.EqualValues(t, "value", value)
	})

	t.Run("set struct (struct value)", func(t *testing.T) {
		type testStruct struct {
			Name string
		}
		err = cache.Set(ctx, "key", testStruct{Name: "value"})
		if err != nil {
			t.Fatal(err)
		}

		var value testStruct
		_, err = cache.Get(ctx, "key", &value)
		if err != nil {
			t.Fatal(err)
		}
		assert.EqualValues(t, "value", value.Name)
	})

	t.Run("set struct (struct key)", func(t *testing.T) {
		type testStruct struct {
			Name string
		}
		err = cache.Set(ctx, testStruct{Name: "key"}, "value")
		if err != nil {
			t.Fatal(err)
		}

		var value string
		_, err = cache.Get(ctx, testStruct{Name: "key"}, &value)
		if err != nil {
			t.Fatal(err)
		}
		assert.EqualValues(t, "value", value)
	})
}

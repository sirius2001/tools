package cache

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/redis/go-redis/v9"
)

type CacheStorer interface {
	Load(key string) (any, bool)
	Store(key string, v any) error
	Delete(key string) error
}

type cache struct {
	rdb    *redis.Client
	data   map[string]any
	locker sync.Mutex
	prefix string
}

// Delete implements CacheStorer.
func (c *cache) Delete(key string) error {
	c.locker.Lock()
	defer c.locker.Unlock()
	if err := c.rdb.HDel(context.Background(), c.prefix, key).Err(); err != nil {
		return err
	}
	delete(c.data, key)
	return nil
}

// Load implements CacheStorer.
func (c *cache) Load(key string) (any, bool) {
	c.locker.Lock()
	defer c.locker.Unlock()
	v, exits := c.data[key]
	return v, exits
}

// Store implements CacheStorer.
func (c *cache) Store(key string, v any) error {
	c.locker.Lock()
	defer c.locker.Unlock()
	if err := c.rdb.HSet(context.Background(), c.prefix, key, v).Err(); err != nil {
		return err
	}
	c.data[key] = v
	return nil
}

func NewCache(rdb *redis.Client, prefix string) CacheStorer {
	data := make(map[string]any)
	objs, err := rdb.HGetAll(context.Background(), prefix).Result()
	if err != nil {
		panic(err)
	}

	for k, v := range objs {
		var obj any
		if err := json.Unmarshal([]byte(v), &obj); err != nil {
			panic(err)
		}
		data[k] = obj
	}

	return &cache{
		rdb:    rdb,
		prefix: prefix,
		data:   data,
	}
}

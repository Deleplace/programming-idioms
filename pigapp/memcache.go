package main

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
)

// It is expected that cached values are volatile and may go away anytime.
// We must not rely on any persistence.
// We must still work correctly if an object is not found in cache.

type memcacheCache struct{}

// If not found in cache, returns nil.
func (memcacheCache) read(c context.Context, key string) (value []byte, err error) {
	cacheItem, err := memcache.Get(c, key)
	if err == memcache.ErrCacheMiss {
		// Item not in the cache
		return nil, nil
	}
	if err != nil {
		// Memcache failure.
		return nil, err
	}
	// Found :)
	return cacheItem.Value, nil
}

func (memcacheCache) write(c context.Context, key string, value []byte, duration time.Duration) error {
	cacheItem := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: duration,
	}
	err := memcache.Set(c, cacheItem)
	return err
}

func (memcacheCache) evict(c context.Context, key string) error {
	err := memcache.Delete(c, key)
	return err
}

func (memcacheCache) flush(c context.Context) error {
	return memcache.Flush(c)
}

func (mc memcacheCache) readMulti(c context.Context, keys []string) (values [][]byte, err1 error) {
	// TODO more efficient batch?
	values = make([][]byte, len(keys))
	for i, key := range keys {
		value, err := mc.read(c, key)
		if err != nil {
			if err1 == nil {
				err1 = err
			}
			continue
		}
		values[i] = value
	}
	return values, err1
}

func (mc memcacheCache) writeMulti(c context.Context, keys []string, values [][]byte, duration time.Duration) error {
	n := len(keys)
	if len(values) != n {
		panic("Inconsistent keys/values")
	}

	items := make([]*memcache.Item, n)
	for i := range keys {
		items[i] = &memcache.Item{
			Key:        keys[i],
			Value:      values[i],
			Expiration: duration,
		}
	}
	return memcache.SetMulti(c, items)
}

func (mc memcacheCache) evictMulti(c context.Context, keys []string) error {
	return memcache.DeleteMulti(c, keys)
}

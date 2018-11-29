package main

import (
	"time"

	"golang.org/x/net/context"
)

// It is expected that cached values are volatile and may go away anytime.
// We must not rely on any persistence.
// We must still work correctly if an object is not found in cache.

// var cache picache = memcacheCache{}
var cache picache = makeRedisCache()

type picache interface {
	read(c context.Context, key string) (value []byte, err error)
	write(c context.Context, key string, value []byte, duration time.Duration) error
	evict(c context.Context, key string) error
	flush(c context.Context) error
	readMulti(c context.Context, keys []string) (values [][]byte, err1 error)
	writeMulti(c context.Context, keys []string, values [][]byte, duration time.Duration) error
	evictMulti(c context.Context, keys []string) error
}

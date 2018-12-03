package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/net/context"
)

//
// Redis Labs Memcached Cloud
// From https://cloud.google.com/appengine/docs/flexible/go/using-redislabs-redis
//
// (this is not MemoryStore)

func makeRedisCache() *redisCache {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisAddr)
			if redisPassword == "" {
				return conn, err
			}
			if err != nil {
				return nil, err
			}
			if _, err := conn.Do("AUTH", redisPassword); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
		// TODO: Tune other settings, like IdleTimeout, MaxActive, MaxIdle, TestOnBorrow.
	}
	return &redisCache{
		pool: redisPool,
	}
}

type redisCache struct {
	pool *redis.Pool
}

// If not found in cache, returns nil.
func (rc redisCache) read(c context.Context, key string) (value []byte, err error) {
	_, endSpan := startSpanf(c, "redisCache.read %q", key)
	defer endSpan()

	redisConn := rc.pool.Get()
	defer redisConn.Close()

	reply, err := redisConn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	if reply == nil {
		// Not found in cache!
		return nil, nil
	}

	data, ok := reply.([]byte)
	if !ok {
		panic(fmt.Sprintf("Can't convert cached %T to []byte", reply))
	}

	// Redis GET works only for string values??
	// https://redis.io/commands/get
	// datastr, ok := reply.(string)
	// if !ok {
	// 	panic(fmt.Sprintf("Can't convert cached %T to string", reply))
	// }
	// data := []byte(datastr)

	return data, nil
}

func (rc redisCache) write(c context.Context, key string, value []byte, duration time.Duration) error {
	redisConn := rc.pool.Get()
	defer redisConn.Close()

	// Redis SET works only for string values??
	// https://redis.io/commands/set
	valuestr := string(value)

	_, err := redisConn.Do("SET", key, valuestr)
	return err
}

func (rc redisCache) evict(c context.Context, key string) error {
	redisConn := rc.pool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("DEL", key)
	return err
}

func (rc redisCache) flush(c context.Context) error {
	redisConn := rc.pool.Get()
	defer redisConn.Close()

	// TODO is this ASYNC?  Can we please make it SYNC?
	_, err := redisConn.Do("FLUSHALL")
	return err
}

func (rc redisCache) readMulti(c context.Context, keys []string) (values [][]byte, err1 error) {
	// TODO more efficient batch?
	values = make([][]byte, len(keys))
	for i, key := range keys {
		value, err := rc.read(c, key)
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

func (rc redisCache) writeMulti(c context.Context, keys []string, values [][]byte, duration time.Duration) error {
	n := len(keys)
	if len(values) != n {
		panic("Inconsistent keys/values")
	}

	// TODO more efficient batch?
	for i := range keys {
		key := keys[i]
		value := values[i]
		err := rc.write(c, key, value, duration)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rc redisCache) evictMulti(c context.Context, keys []string) error {
	// TODO more efficient batch?
	var err1 error
	// Try to evict *all* keys
	for _, key := range keys {
		err := rc.evict(c, key)
		if err != nil && err1 == nil {
			err1 = err
		}
	}
	// then return first error, if any
	return err1
}

package pigae

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

// Sometimes we want to saved whole blocks of template-generated HTML
// into HTML, and serve it again later.

// htmlCacheRead returns previously saved bytes for this key,
// It returns nil if not found, or expired, or on memcache error.
//
// There is no guarantee that previously cached data will be found,
// because memcache entries may vanish anytime, even before expiration.
func htmlCacheRead(c context.Context, key string) []byte {
	var cacheItem *memcache.Item
	var err error
	if cacheItem, err = memcache.Get(c, key); err == memcache.ErrCacheMiss {
		// Item not in the cache
		return nil
	} else if err != nil {
		// Memcache failure. Ignore.
		return nil
	}
	// Found :)
	return cacheItem.Value
}

// htmlCacheWrite saves bytes for given key.
// Failures are ignored.
func htmlCacheWrite(c context.Context, key string, data []byte, duration time.Duration) {
	cacheItem := &memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: duration,
	}
	_ = memcache.Set(c, cacheItem)
}

// Data changes should lead to cache entries invalidation.
func htmlCacheEvict(c context.Context, key string) {
	_ = memcache.Delete(c, key)
	// See also htmlUncacheIdiomAndImpls
}

// When expected data may be >1MB.
func htmlCacheZipRead(c context.Context, key string) []byte {
	zipdata := htmlCacheRead(c, key)
	if zipdata == nil {
		return nil
	}
	zipbuffer := bytes.NewBuffer(zipdata)
	zipreader, err := gzip.NewReader(zipbuffer)
	if err != nil {
		log.Errorf(c, "Reading zip memcached entry %q: %v", key, err)
		// Ignore failure
		return nil
	}
	buffer, err := ioutil.ReadAll(zipreader)
	if err != nil {
		log.Errorf(c, "Reading zip memcached entry %q: %v", key, err)
	}
	log.Debugf(c, "Reading %d bytes out of %d gzip bytes for entry %q", len(buffer), len(zipdata), key)
	return buffer
}

// When expected data may be >1MB.
func htmlCacheZipWrite(c context.Context, key string, data []byte, duration time.Duration) {
	var zipbuffer bytes.Buffer
	zipwriter := gzip.NewWriter(&zipbuffer)
	_, err := zipwriter.Write(data)
	if err != nil {
		log.Errorf(c, "Writing zip memcached entry %q: %v", key, err)
		// Ignore failure
		return
	}
	_ = zipwriter.Close()
	log.Debugf(c, "Writing %d gzip bytes out of %d data bytes for entry %q", zipbuffer.Len(), len(data), key)
	htmlCacheWrite(c, key, zipbuffer.Bytes(), duration)
}

func htmlUncacheIdiomAndImpls(c context.Context, idiom *Idiom) {
	//
	// There are only two hard things in Computer Science: cache invalidation and naming things.
	//
	cachekeys := make([]string, 0, 1+len(idiom.Implementations))
	cachekeys = append(cachekeys, NiceIdiomRelativeURL(idiom))
	for _, impl := range idiom.Implementations {
		cachekeys = append(cachekeys, NiceImplRelativeURL(idiom, impl.Id, impl.LanguageName))
	}
	err := memcache.DeleteMulti(c, cachekeys)
	if err != nil {
		// log.Errorf(c, "Uncaching idiom %d: %v", idiom.Id, err)
		// A lot of impl HTML paages won't be in cache, which will cause
		// en error. Never mind, just ignore.
	}
}

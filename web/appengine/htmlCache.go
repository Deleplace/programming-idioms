package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"strconv"
	"time"

	. "github.com/Deleplace/programming-idioms/idioms"

	"context"

	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/taskqueue"
)

// Sometimes we want to saved whole blocks of template-generated HTML
// into HTML, and serve it again later.

// htmlCacheRead returns previously saved bytes for this key,
// It returns nil if not found, or expired, or on memcache error.
//
// There is no guarantee that previously cached data will be found,
// because memcache entries may vanish anytime, even before expiration.
func htmlCacheRead(ctx context.Context, key string) []byte {
	cacheItem, err := memcache.Get(ctx, key)
	if err == memcache.ErrCacheMiss {
		// Item not in the cache
		return nil
	}
	if err != nil {
		// Memcache failure. Ignore.
		return nil
	}
	// Found :)
	return cacheItem.Value
}

// htmlCacheWrite saves bytes for given key.
// Failures are ignored.
func htmlCacheWrite(ctx context.Context, key string, data []byte, duration time.Duration) {
	cacheItem := &memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: duration,
	}
	_ = memcache.Set(ctx, cacheItem)
}

// Data changes should lead to cache entries invalidation.
func htmlCacheEvict(ctx context.Context, key string) {
	_ = memcache.Delete(ctx, key)
	// See also htmlUncacheIdiomAndImpls
}

// When expected data may be >1MB.
func htmlCacheZipRead(ctx context.Context, key string) []byte {
	zipdata := htmlCacheRead(ctx, key)
	if zipdata == nil {
		return nil
	}
	zipbuffer := bytes.NewBuffer(zipdata)
	zipreader, err := gzip.NewReader(zipbuffer)
	if err != nil {
		log.Errorf(ctx, "Reading zip memcached entry %q: %v", key, err)
		// Ignore failure
		return nil
	}
	buffer, err := ioutil.ReadAll(zipreader)
	if err != nil {
		log.Errorf(ctx, "Reading zip memcached entry %q: %v", key, err)
	}
	log.Debugf(ctx, "Reading %d bytes out of %d gzip bytes for entry %q", len(buffer), len(zipdata), key)
	return buffer
}

// When expected data may be >1MB.
func htmlCacheZipWrite(ctx context.Context, key string, data []byte, duration time.Duration) {
	var zipbuffer bytes.Buffer
	zipwriter := gzip.NewWriter(&zipbuffer)
	_, err := zipwriter.Write(data)
	if err != nil {
		log.Errorf(ctx, "Writing zip memcached entry %q: %v", key, err)
		// Ignore failure
		return
	}
	_ = zipwriter.Close()
	log.Debugf(ctx, "Writing %d gzip bytes out of %d data bytes for entry %q", zipbuffer.Len(), len(data), key)
	htmlCacheWrite(ctx, key, zipbuffer.Bytes(), duration)
}

func htmlUncacheIdiomAndImpls(ctx context.Context, idiom *Idiom) {
	//
	// There are only two hard things in Computer Science: cache invalidation and naming things.
	//
	log.Infof(ctx, "Evicting HTML cached pages for idiom %d %q", idiom.Id, idiom.Title)

	cachekeys := make([]string, 0, 1+len(idiom.Implementations))
	cachekeys = append(cachekeys, NiceIdiomRelativeURL(idiom))
	for _, impl := range idiom.Implementations {
		cachekeys = append(cachekeys, NiceImplRelativeURL(idiom, impl.Id, impl.LanguageName))
	}
	err := memcache.DeleteMulti(ctx, cachekeys)
	if err != nil {
		// log.Errorf(ctx, "Uncaching idiom %d: %v", idiom.Id, err)
		// A lot of impl HTML paages won't be in cache, which will cause
		// en error. Never mind, just ignore.
	}
}

func htmlRecacheNowAndTomorrow(ctx context.Context, idiomID int) error {
	log.Debugf(ctx, "Creating html recache tasks for idiom %d", idiomID)
	// These 2 task submissions may take several 10s of ms,
	// thus we decide to submit them as a small batch.

	// Now
	t1, err1 := recacheHtmlIdiom.Task(idiomID)
	if err1 != nil {
		return err1
	}

	// Tomorrow
	t2, err2 := recacheHtmlIdiom.Task(idiomID)
	if err2 != nil {
		return err2
	}
	t2.Delay = 24*time.Hour + 10*time.Minute

	_, err := taskqueue.AddMulti(ctx, []*taskqueue.Task{t1, t2}, "")
	return err
}

var recacheHtmlIdiom, recacheHtmlImpl *delay.Function

func init() {
	recacheHtmlIdiom = delay.Func("recache-html-idiom", func(ctx context.Context, idiomID int) {
		log.Infof(ctx, "Start recaching HTML for idiom %d", idiomID)
		_, idiom, err := dao.getIdiom(ctx, idiomID)
		if err != nil {
			log.Errorf(ctx, "recacheHtmlIdiom: %v", err)
			return
		}

		path := NiceIdiomRelativeURL(idiom)
		var buffer bytes.Buffer
		vars := map[string]string{
			"idiomId":    strconv.Itoa(idiomID),
			"idiomTitle": uriNormalize(idiom.Title),
		}
		err = generateIdiomDetailPage(ctx, &buffer, vars)
		if err != nil {
			log.Errorf(ctx, "recacheHtmlIdiom: %v", err)
			return
		}
		htmlCacheWrite(ctx, path, buffer.Bytes(), 24*time.Hour)

		// Then, create async task for each impl to be HTML-recached
		for _, impl := range idiom.Implementations {
			implPath := NiceImplRelativeURL(idiom, impl.Id, impl.LanguageName)
			recacheHtmlImpl.Call(ctx, implPath, idiom.Id, idiom.Title, impl.Id, impl.LanguageName)
		}

	})

	recacheHtmlImpl = delay.Func("recache-html-impl", func(
		ctx context.Context,
		implPath string,
		idiomID int,
		idiomTitle string,
		implID int,
		implLang string,
	) {
		log.Infof(ctx, "Recaching HTML for %s", implPath)
		// TODO call idiomDetail(fakeWriter, fakeRequest)

		var buffer bytes.Buffer
		vars := map[string]string{
			"idiomId":    strconv.Itoa(idiomID),
			"idiomTitle": uriNormalize(idiomTitle),
			"implId":     strconv.Itoa(implID),
			"implLang":   implLang,
		}
		err := generateIdiomDetailPage(ctx, &buffer, vars)
		if err != nil {
			log.Errorf(ctx, "recacheHtmlImpl: %v", err)
			return
		}
		htmlCacheWrite(ctx, implPath, buffer.Bytes(), 24*time.Hour)
	})
}

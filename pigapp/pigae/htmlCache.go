package pigae

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"strconv"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
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
func htmlCacheRead(c context.Context, key string) []byte {
	cacheItem, err := memcache.Get(c, key)
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
	log.Infof(c, "Evicting HTML cached pages for idiom %d %q", idiom.Id, idiom.Title)

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

func htmlRecacheNowAndTomorrow(c context.Context, idiomID int) error {
	// Now
	err := recacheHtmlIdiom.Call(c, idiomID)
	if err != nil {
		return err
	}

	// Tomorrow
	t, _ := recacheHtmlIdiom.Task(idiomID)
	t.Delay = 24*time.Hour + 10*time.Minute
	_, err = taskqueue.Add(c, t, "")
	return err
}

var recacheHtmlIdiom, recacheHtmlImpl *delay.Function

func init() {
	recacheHtmlIdiom = delay.Func("recache-html-idiom", func(c context.Context, idiomID int) {
		log.Infof(c, "Start recaching HTML for idiom %d", idiomID)
		_, idiom, err := dao.getIdiom(c, idiomID)
		if err != nil {
			log.Errorf(c, "recacheHtmlIdiom: %v", err)
			return
		}

		path := NiceIdiomRelativeURL(idiom)
		var buffer bytes.Buffer
		vars := map[string]string{
			"idiomId":    strconv.Itoa(idiomID),
			"idiomTitle": uriNormalize(idiom.Title),
		}
		err = generateIdiomDetailPage(c, &buffer, vars)
		if err != nil {
			log.Errorf(c, "recacheHtmlIdiom: %v", err)
			return
		}
		htmlCacheWrite(c, path, buffer.Bytes(), 24*time.Hour)

		// Then, create async task for each impl to be HTML-recached
		for _, impl := range idiom.Implementations {
			implPath := NiceImplRelativeURL(idiom, impl.Id, impl.LanguageName)
			recacheHtmlImpl.Call(c, implPath, idiom.Id, idiom.Title, impl.Id, impl.LanguageName)
		}

	})

	recacheHtmlImpl = delay.Func("recache-html-impl", func(
		c context.Context,
		implPath string,
		idiomID int,
		idiomTitle string,
		implID int,
		implLang string,
	) {
		log.Infof(c, "Recaching HTML for %s", implPath)
		// TODO call idiomDetail(fakeWriter, fakeRequest)

		var buffer bytes.Buffer
		vars := map[string]string{
			"idiomId":    strconv.Itoa(idiomID),
			"idiomTitle": uriNormalize(idiomTitle),
			"implId":     strconv.Itoa(implID),
			"implLang":   implLang,
		}
		err := generateIdiomDetailPage(c, &buffer, vars)
		if err != nil {
			log.Errorf(c, "recacheHtmlImpl: %v", err)
			return
		}
		htmlCacheWrite(c, implPath, buffer.Bytes(), 24*time.Hour)
	})
}

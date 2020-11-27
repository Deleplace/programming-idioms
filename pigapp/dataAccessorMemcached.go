package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

// This source file has a lot of duplicated code : "if cached then return else datastore and cache".
// TODO: find a smarter design for this "proxy" type which applies basically the same behavior to
// all read methods, and the same behavior to all write methods.

// MemcacheDatastoreAccessor accessor uses a MemCache for standard CRUD.
//
// Some methods are not redefined : randomIdiom, nextIdiomID, nextImplID, processUploadFile, processUploadFiles
type MemcacheDatastoreAccessor struct {
	GaeDatastoreAccessor
}

func (a *MemcacheDatastoreAccessor) cacheValue(ctx context.Context, cacheKey string, data interface{}, expiration time.Duration) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(&data)
	if err != nil {
		log.Debugf(ctx, "Failed encoding for cache[%v] : %v", cacheKey, err)
		return err
	}
	cacheItem := &memcache.Item{
		Key:        cacheKey,
		Value:      buffer.Bytes(),
		Expiration: expiration,
	}
	// Set the item, unconditionally
	err = memcache.Set(ctx, cacheItem)
	if err != nil {
		log.Debugf(ctx, "Failed setting cache[%v] : %v", cacheKey, err)
	} else {
		// log.Debugf(ctx, "Successfully set cache[%v]", cacheKey)
	}
	return err
}

func (a *MemcacheDatastoreAccessor) cacheValues(ctx context.Context, cacheKeys []string, data []interface{}, expiration time.Duration) error {
	if len(cacheKeys) != len(data) {
		panic(fmt.Errorf("Wrong params length %d, %d", len(cacheKeys), len(data)))
	}
	N := len(cacheKeys)

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	items := make([]*memcache.Item, N)
	for i, cacheKey := range cacheKeys {
		cacheData := data[i]
		err := enc.Encode(&cacheData)
		if err != nil {
			log.Debugf(ctx, "Failed encoding for cache[%v] : %v", cacheKey, err)
			return err
		}
		cacheItem := &memcache.Item{
			Key:        cacheKey,
			Value:      buffer.Bytes(),
			Expiration: expiration,
		}
		items[i] = cacheItem
	}

	// Set the items, unconditionally, in 1 batch call
	err := memcache.SetMulti(ctx, items)
	if err != nil {
		log.Debugf(ctx, "Failed setting cache items: %v", cacheKeys, err)
	}
	return err
}

func (a *MemcacheDatastoreAccessor) cacheSameValues(ctx context.Context, cacheKeys []string, data interface{}, expiration time.Duration) error {
	N := len(cacheKeys)

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(&data)

	if err != nil {
		log.Debugf(ctx, "Failed encoding for cache keys [%v] : %v", cacheKeys, err)
		return err
	}

	items := make([]*memcache.Item, N)
	for i, cacheKey := range cacheKeys {
		cacheItem := &memcache.Item{
			Key:        cacheKey,
			Value:      buffer.Bytes(),
			Expiration: expiration,
		}
		items[i] = cacheItem
	}

	// Set the items, unconditionally, in 1 batch call
	err = memcache.SetMulti(ctx, items)
	if err != nil {
		log.Debugf(ctx, "Failed setting cache items: %v", cacheKeys, err)
	}
	return err
}

// A shortcut for caching the datastoreKey + value
func (a *MemcacheDatastoreAccessor) cacheKeyValue(ctx context.Context, cacheKey string, datastoreKey *datastore.Key, entity interface{}, expiration time.Duration) error {
	kae := &KeyAndEntity{datastoreKey, entity}
	return a.cacheValue(ctx, cacheKey, kae, expiration)
}

// A shortcut for caching the datastoreKeys + values
func (a *MemcacheDatastoreAccessor) cacheKeysValues(ctx context.Context, cacheKeys []string, datastoreKeys []*datastore.Key, entities []interface{}, expiration time.Duration) error {
	if len(cacheKeys) != len(datastoreKeys) || len(cacheKeys) != len(entities) {
		panic(fmt.Errorf("Wrong params length %d, %d, %d", len(cacheKeys), len(datastoreKeys), len(entities)))
	}
	N := len(cacheKeys)

	items := make([]interface{}, N)
	for i, datastoreKey := range datastoreKeys {
		kae := &KeyAndEntity{datastoreKey, entities[i]}
		items[i] = kae
	}
	return a.cacheValues(ctx, cacheKeys, items, expiration)
}

// A shortcut for caching the same datastore key and same value (to encode once) for each cacheKey
func (a *MemcacheDatastoreAccessor) cacheKeysSameValue(ctx context.Context, cacheKeys []string, datastoreKey *datastore.Key, entity interface{}, expiration time.Duration) error {
	kae := &KeyAndEntity{datastoreKey, entity}
	return a.cacheSameValues(ctx, cacheKeys, kae, expiration)
}

// Just a shortcut for caching the pair
func (a *MemcacheDatastoreAccessor) cachePair(ctx context.Context, cacheKey string, first interface{}, second interface{}, expiration time.Duration) error {
	pair := &pair{first, second}
	return a.cacheValue(ctx, cacheKey, pair, expiration)
}

func (a *MemcacheDatastoreAccessor) readCache(ctx context.Context, cacheKey string) (interface{}, error) {
	// Get the item from the memcache
	var cacheItem *memcache.Item
	var err error
	if cacheItem, err = memcache.Get(ctx, cacheKey); err == memcache.ErrCacheMiss {
		// Item not in the cache
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(cacheItem.Value) // todo avoid bytes copy ??
	dec := gob.NewDecoder(buffer)
	var data interface{}
	err = dec.Decode(&data)
	return data, err
}

func init() {
	gob.Register(&KeyAndEntity{})
	gob.Register(&pair{})
	gob.Register(&Idiom{})
	gob.Register([]*datastore.Key{})
	gob.Register([]*Idiom{})
	gob.Register([]string{})
	gob.Register(map[string]bool{})
	gob.Register(Toggles{})
	gob.Register(&ApplicationConfig{})
	gob.Register([]*MessageForUser{})
}

// KeyAndEntity is a specific pair wrapper.
type KeyAndEntity struct {
	Key    *datastore.Key
	Entity interface{}
}

type pair struct {
	First  interface{}
	Second interface{}
}

func (a *MemcacheDatastoreAccessor) recacheIdiom(ctx context.Context, datastoreKey *datastore.Key, idiom *Idiom, invalidateHTML bool) error {
	cacheKey := fmt.Sprintf("getIdiom(%v)", idiom.Id)
	err := a.cacheKeyValue(ctx, cacheKey, datastoreKey, idiom, 24*time.Hour)
	if err != nil {
		log.Errorf(ctx, err.Error())
		return err
	}

	// Batch memcache set call
	N := len(idiom.Implementations)
	cacheKeys := make([]string, N)
	for i, impl := range idiom.Implementations {
		cacheKeys[i] = fmt.Sprintf("getIdiomByImplID(%v)", impl.Id)
	}
	err = a.cacheKeysSameValue(ctx, cacheKeys, datastoreKey, idiom, 24*time.Hour)
	if err != nil {
		log.Errorf(ctx, err.Error())
		return err
	}
	// Unfortunately, some previous "getIdiomByImplID(xyz)" might be left uninvalidated.
	// (theoretically)

	if invalidateHTML {
		// Note that cached HTML pages are just evicted, not regenerated here.
		htmlUncacheIdiomAndImpls(ctx, idiom)
	}

	return err
}

func (a *MemcacheDatastoreAccessor) uncacheIdiom(ctx context.Context, idiom *Idiom) error {
	cacheKeys := make([]string, 1+len(idiom.Implementations))
	cacheKeys[0] = fmt.Sprintf("getIdiom(%v)", idiom.Id)
	for i, impl := range idiom.Implementations {
		cacheKeys[1+i] = fmt.Sprintf("getIdiomByImplID(%v)", impl.Id)
	}

	err := memcache.DeleteMulti(ctx, cacheKeys)
	if err != nil {
		log.Errorf(ctx, err.Error())
	}

	// Cached HTML pages.
	htmlUncacheIdiomAndImpls(ctx, idiom)

	return err
}

func (a *MemcacheDatastoreAccessor) getIdiom(ctx context.Context, idiomID int) (*datastore.Key, *Idiom, error) {
	cacheKey := fmt.Sprintf("getIdiom(%v)", idiomID)
	data, cacheerr := a.readCache(ctx, cacheKey)
	if cacheerr != nil {
		log.Errorf(ctx, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.getIdiom(ctx, idiomID)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		key, idiom, err := a.GaeDatastoreAccessor.getIdiom(ctx, idiomID)
		if err == nil {
			err2 := a.recacheIdiom(ctx, key, idiom, false)
			logIf(err2, log.Errorf, ctx, "recaching idiom")
		}
		return key, idiom, err
	}
	// Found in cache :)
	kae := data.(*KeyAndEntity)
	key := kae.Key
	idiom := kae.Entity.(*Idiom)
	return key, idiom, nil
}

func (a *MemcacheDatastoreAccessor) getIdiomByImplID(ctx context.Context, implID int) (*datastore.Key, *Idiom, error) {
	cacheKey := fmt.Sprintf("getIdiomByImplID(%v)", implID)
	data, cacheerr := a.readCache(ctx, cacheKey)
	if cacheerr != nil {
		log.Errorf(ctx, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.getIdiomByImplID(ctx, implID)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		key, idiom, err := a.GaeDatastoreAccessor.getIdiomByImplID(ctx, implID)
		if err == nil {
			err2 := a.cacheKeyValue(ctx, cacheKey, key, idiom, 24*time.Hour)
			logIf(err2, log.Errorf, ctx, "caching idiom")
		}
		return key, idiom, err
	}
	// Found in cache :)
	kae := data.(*KeyAndEntity)
	key := kae.Key
	idiom := kae.Entity.(*Idiom)
	return key, idiom, nil
}

func (a *MemcacheDatastoreAccessor) saveNewIdiom(ctx context.Context, idiom *Idiom) (*datastore.Key, error) {
	key, err := a.GaeDatastoreAccessor.saveNewIdiom(ctx, idiom)
	if err == nil {
		err2 := a.recacheIdiom(ctx, key, idiom, false)
		logIf(err2, log.Errorf, ctx, "saving new idiom")
	}
	_ = memcache.DeleteMulti(ctx, []string{
		"about-block-language-coverage",
		"getAllIdioms(399,-ImplCount)",
	})
	return key, err
}

func (a *MemcacheDatastoreAccessor) saveExistingIdiom(ctx context.Context, key *datastore.Key, idiom *Idiom) error {
	// It is important to invalidate cache with OLD paths, thus before saving
	if _, oldIdiomValue, err := a.getIdiom(ctx, idiom.Id); err == nil {
		htmlUncacheIdiomAndImpls(ctx, oldIdiomValue)
	}

	log.Infof(ctx, "Saving idiom #%v: %v", idiom.Id, idiom.Title)
	err := a.GaeDatastoreAccessor.saveExistingIdiom(ctx, key, idiom)
	if err == nil {
		log.Infof(ctx, "Saved idiom #%v, version %v", idiom.Id, idiom.Version)
		err2 := a.recacheIdiom(ctx, key, idiom, false)
		logIf(err2, log.Errorf, ctx, "saving existing idiom")
	}
	_ = memcache.DeleteMulti(ctx, []string{
		"about-block-language-coverage",
		"getAllIdioms(399,-ImplCount)",
	})
	return err
}

func (a *MemcacheDatastoreAccessor) stealthIncrementIdiomRating(ctx context.Context, idiomID int, delta int) (*datastore.Key, *Idiom, error) {
	key, idiom, err := a.GaeDatastoreAccessor.stealthIncrementIdiomRating(ctx, idiomID, delta)
	err2 := a.recacheIdiom(ctx, key, idiom, true)
	logIf(err2, log.Errorf, ctx, "updating idiom rating")
	return key, idiom, err
}

func (a *MemcacheDatastoreAccessor) stealthIncrementImplRating(ctx context.Context, idiomID, implID int, delta int) (key *datastore.Key, idiom *Idiom, newImplRating int, err error) {
	key, idiom, newImplRating, err = a.GaeDatastoreAccessor.stealthIncrementImplRating(ctx, idiomID, implID, delta)
	err2 := a.recacheIdiom(ctx, key, idiom, true)
	logIf(err2, log.Errorf, ctx, "updating impl rating")
	return
}

func (a *MemcacheDatastoreAccessor) getAllIdioms(ctx context.Context, limit int, order string) ([]*datastore.Key, []*Idiom, error) {
	cacheKey := fmt.Sprintf("getAllIdioms(%v,%v)", limit, order)
	data, cacheerr := a.readZipCache(ctx, cacheKey)
	if cacheerr != nil {
		log.Errorf(ctx, "Reading zip cache for %q: %v", cacheKey, cacheerr)
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.getAllIdioms(ctx, limit, order)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		keys, idioms, err := a.GaeDatastoreAccessor.getAllIdioms(ctx, limit, order)
		if err == nil {
			err2 := a.cacheZipPair(ctx, cacheKey, keys, idioms, 12*time.Hour)
			logIf(err2, log.Errorf, ctx, "caching all idioms")
		}
		return keys, idioms, err
	}
	log.Infof(ctx, "Found %q in cache :)", cacheKey)
	pair := data.(*pair)
	keys := pair.First.([]*datastore.Key)
	idioms := pair.Second.([]*Idiom)
	return keys, idioms, nil
}

func (a *MemcacheDatastoreAccessor) deleteAllIdioms(ctx context.Context) error {
	err := a.GaeDatastoreAccessor.deleteAllIdioms(ctx)
	if err != nil {
		return err
	}
	// Cache : the nuclear option!
	return memcache.Flush(ctx)
}

func (a *MemcacheDatastoreAccessor) unindexAll(ctx context.Context) error {
	return a.GaeDatastoreAccessor.unindexAll(ctx)
}

func (a *MemcacheDatastoreAccessor) unindex(ctx context.Context, idiomID int) error {
	return a.GaeDatastoreAccessor.unindex(ctx, idiomID)
}

func (a *MemcacheDatastoreAccessor) deleteIdiom(ctx context.Context, idiomID int, why string) error {
	// Clear cache entries
	_, idiom, err := a.GaeDatastoreAccessor.getIdiom(ctx, idiomID)
	if err == nil {
		err2 := a.uncacheIdiom(ctx, idiom)
		logIf(err2, log.Errorf, ctx, "deleting idiom")
	} else {
		log.Errorf(ctx, "Failed to load idiom %d to uncache: %v", idiomID, err)
	}

	// Delete in datastore
	return a.GaeDatastoreAccessor.deleteIdiom(ctx, idiomID, why)
}

func (a *MemcacheDatastoreAccessor) deleteImpl(ctx context.Context, idiomID int, implID int, why string) error {
	// Clear cache entries
	_, idiom, err := a.GaeDatastoreAccessor.getIdiom(ctx, idiomID)
	if err == nil {
		err2 := a.uncacheIdiom(ctx, idiom)
		logIf(err2, log.Errorf, ctx, "deleting impl")
	}

	// Delete in datastore
	err = a.GaeDatastoreAccessor.deleteImpl(ctx, idiomID, implID, why)
	return err
}

func (a *MemcacheDatastoreAccessor) searchIdiomsByWordsWithFavorites(ctx context.Context, typedWords, typedLangs []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error) {
	// Personalized searches not cached (yet)
	return a.GaeDatastoreAccessor.searchIdiomsByWordsWithFavorites(ctx, typedWords, typedLangs, favoriteLangs, seeNonFavorite, limit)
}

func (a *MemcacheDatastoreAccessor) searchImplIDs(ctx context.Context, words, langs []string) (map[string]bool, error) {
	// TODO cache this... or not.
	return a.GaeDatastoreAccessor.searchImplIDs(ctx, words, langs)
}

func (a *MemcacheDatastoreAccessor) searchIdiomsByLangs(ctx context.Context, langs []string, limit int) ([]*Idiom, error) {
	cacheKey := fmt.Sprintf("searchIdiomsByLangs(%v,%v)", langs, limit)
	//log.Debugf(ctx, cacheKey)
	data, cacheerr := a.readCache(ctx, cacheKey)
	if cacheerr != nil {
		log.Errorf(ctx, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.searchIdiomsByLangs(ctx, langs, limit)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		idioms, err := a.GaeDatastoreAccessor.searchIdiomsByLangs(ctx, langs, limit)
		if err == nil {
			// Search results will have a 10mn lag after an idiom/impl creation/update.
			err2 := a.cacheValue(ctx, cacheKey, idioms, 10*time.Minute)
			logIf(err2, log.Errorf, ctx, "caching search results by langs")
		}
		return idioms, err
	}
	// Found in cache :)
	idioms := data.([]*Idiom)
	return idioms, nil
}

func (a *MemcacheDatastoreAccessor) recentIdioms(ctx context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	cacheKey := fmt.Sprintf("recentIdioms(%v,%v,%v)", favoriteLangs, showOther, n)
	//log.Debugf(ctx, cacheKey)
	data, cacheerr := a.readCache(ctx, cacheKey)
	if cacheerr != nil {
		log.Errorf(ctx, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.recentIdioms(ctx, favoriteLangs, showOther, n)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		idioms, err := a.GaeDatastoreAccessor.recentIdioms(ctx, favoriteLangs, showOther, n)
		if err == nil {
			// "Popular idioms" will have a 10mn lag after an idiom/impl creation.
			err2 := a.cacheValue(ctx, cacheKey, idioms, 10*time.Minute)
			logIf(err2, log.Errorf, ctx, "caching recent idioms")
		}
		return idioms, err
	}
	// Found in cache :)
	idioms := data.([]*Idiom)
	return idioms, nil
}

func (a *MemcacheDatastoreAccessor) popularIdioms(ctx context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	cacheKey := fmt.Sprintf("popularIdioms(%v,%v,%v)", favoriteLangs, showOther, n)
	//log.Debugf(ctx, cacheKey)
	data, cacheerr := a.readCache(ctx, cacheKey)
	if cacheerr != nil {
		log.Errorf(ctx, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.popularIdioms(ctx, favoriteLangs, showOther, n)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		idioms, err := a.GaeDatastoreAccessor.popularIdioms(ctx, favoriteLangs, showOther, n)
		if err == nil {
			// "Popular idioms" will have a 10mn lag after an idiom/impl creation.
			err2 := a.cacheValue(ctx, cacheKey, idioms, 10*time.Minute)
			logIf(err2, log.Errorf, ctx, "caching popular idioms")
		}
		return idioms, err
	}
	// Found in cache :)
	idioms := data.([]*Idiom)
	return idioms, nil
}

// TODO cache this
func (a *MemcacheDatastoreAccessor) idiomsFilterOrder(ctx context.Context, favoriteLangs []string, limitEachLang int, showOther bool, sortOrder string) ([]*Idiom, error) {
	idioms, err := a.GaeDatastoreAccessor.idiomsFilterOrder(ctx, favoriteLangs, limitEachLang, showOther, sortOrder)
	return idioms, err
}

func (a *MemcacheDatastoreAccessor) getAppConfig(ctx context.Context) (ApplicationConfig, error) {
	cacheKey := "getAppConfig()"
	data, cacheerr := a.readCache(ctx, cacheKey)
	if cacheerr != nil {
		log.Errorf(ctx, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.getAppConfig(ctx)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		appConfig, err := a.GaeDatastoreAccessor.getAppConfig(ctx)
		if err == nil {
			log.Infof(ctx, "Retrieved ApplicationConfig (Toggles) from Datastore")
			err2 := a.cacheValue(ctx, cacheKey, appConfig, 24*time.Hour)
			logIf(err2, log.Errorf, ctx, "caching app config")
		}
		return appConfig, err
	}
	// Found in cache :)
	appConfig := data.(*ApplicationConfig)
	return *appConfig, nil
}

func (a *MemcacheDatastoreAccessor) saveAppConfig(ctx context.Context, appConfig ApplicationConfig) error {
	err := memcache.Flush(ctx)
	if err != nil {
		return err
	}
	return a.GaeDatastoreAccessor.saveAppConfig(ctx, appConfig)
	// TODO force toggles refresh for all instances, after memcache flush
}

func (a *MemcacheDatastoreAccessor) saveAppConfigProperty(ctx context.Context, prop AppConfigProperty) error {
	err := memcache.Flush(ctx)
	if err != nil {
		return err
	}
	return a.GaeDatastoreAccessor.saveAppConfigProperty(ctx, prop)
	// TODO force toggles refresh for all instances, after memcache flush
}

func (a *MemcacheDatastoreAccessor) deleteCache(ctx context.Context) error {
	return memcache.Flush(ctx)
}

func (a *MemcacheDatastoreAccessor) revert(ctx context.Context, idiomID int, version int) (*Idiom, error) {
	idiom, err := a.GaeDatastoreAccessor.revert(ctx, idiomID, version)
	if err != nil {
		return idiom, err
	}
	err2 := a.uncacheIdiom(ctx, idiom)
	logIf(err2, log.Errorf, ctx, "uncaching idiom")
	return idiom, err
}

func (a *MemcacheDatastoreAccessor) historyRestore(ctx context.Context, idiomID int, version int, restoreUser string, why string) (*Idiom, error) {

	idiom, err := a.GaeDatastoreAccessor.historyRestore(ctx, idiomID, version, restoreUser, why)

	if err != nil {
		// Uncaching is useful, even when the restore has failed
		_, idiom, err2 := a.GaeDatastoreAccessor.getIdiom(ctx, idiomID)
		if err2 == nil {
			_ = a.uncacheIdiom(ctx, idiom)
		}
		return idiom, err
	}

	errUCI := a.uncacheIdiom(ctx, idiom)
	logIf(errUCI, log.Errorf, ctx, "uncaching idiom")
	return idiom, err
}

func (a *MemcacheDatastoreAccessor) saveNewMessage(ctx context.Context, msg *MessageForUser) (*datastore.Key, error) {
	key, err := a.GaeDatastoreAccessor.saveNewMessage(ctx, msg)
	if err != nil {
		return key, err
	}

	cacheKey := "getMessagesForUser(" + msg.Username + ")"
	err = memcache.Delete(ctx, cacheKey)
	if err == memcache.ErrCacheMiss {
		// No problem if wasn't in cache anyway
		err = nil
	}
	return key, err
}

func (a *MemcacheDatastoreAccessor) getMessagesForUser(ctx context.Context, username string) ([]*datastore.Key, []*MessageForUser, error) {
	cacheKey := "getMessagesForUser(" + username + ")"

	data, cacheerr := a.readCache(ctx, cacheKey)
	if cacheerr != nil {
		// Ouch.
		return nil, nil, cacheerr
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		keys, messages, err := a.GaeDatastoreAccessor.getMessagesForUser(ctx, username)
		if err == nil {
			err2 := a.cachePair(ctx, cacheKey, keys, messages, 2*time.Hour)
			logIf(err2, log.Errorf, ctx, "caching user messages")
		}
		return keys, messages, err
	}
	pair := data.(*pair)
	keys := pair.First.([]*datastore.Key)
	messages := pair.Second.([]*MessageForUser)
	return keys, messages, nil
}

func (a *MemcacheDatastoreAccessor) dismissMessage(ctx context.Context, key *datastore.Key) (*MessageForUser, error) {
	msg, err := a.GaeDatastoreAccessor.dismissMessage(ctx, key)
	if err != nil {
		return nil, err
	}

	cacheKey := "getMessagesForUser(" + msg.Username + ")"
	err = memcache.Delete(ctx, cacheKey)
	if err == memcache.ErrCacheMiss {
		// No problem if wasn't in cache anyway
		err = nil
	}
	return msg, err
}

// When expected data may be >1MB, but compressible <1MB.
func (a *MemcacheDatastoreAccessor) readZipCache(ctx context.Context, cacheKey string) (interface{}, error) {
	var zipCacheItem *memcache.Item
	var err error
	if zipCacheItem, err = memcache.Get(ctx, cacheKey); err == memcache.ErrCacheMiss {
		// Item not in the cache
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	zipdata := zipCacheItem.Value
	zipbuffer := bytes.NewBuffer(zipdata)
	zipreader, err := gzip.NewReader(zipbuffer)
	if err != nil {
		return nil, fmt.Errorf("Reading zip memcached entry %q: %v", cacheKey, err)
	}

	dec := gob.NewDecoder(zipreader)
	var data interface{}
	err = dec.Decode(&data)
	return data, err
}

// When expected data may be >1MB, but compressible <1MB.
func (a *MemcacheDatastoreAccessor) cacheZipValue(ctx context.Context, cacheKey string, data interface{}, expiration time.Duration) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(&data)
	if err != nil {
		log.Debugf(ctx, "Failed encoding for cache[%v] : %v", cacheKey, err)
		return err
	}

	var zipbuffer bytes.Buffer
	zipwriter := gzip.NewWriter(&zipbuffer)
	_, err = zipwriter.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("Writing zip memcached entry %q: %v", cacheKey, err)
	}
	err = zipwriter.Close()
	if err != nil {
		return fmt.Errorf("Writing (Close) zip memcached entry %q: %v", cacheKey, err)
	}
	const KB = 1024
	const MB = 1024 * KB
	if zipbuffer.Len() > 1*MB {
		return fmt.Errorf("Not caching %q: %dkB (gzipped) is too large", cacheKey, zipbuffer.Len()/KB)
	}
	log.Debugf(ctx, "Writing %d gzip bytes out of %d data bytes for entry %q", zipbuffer.Len(), buffer.Len(), cacheKey)

	cacheItem := &memcache.Item{
		Key:        cacheKey,
		Value:      zipbuffer.Bytes(),
		Expiration: expiration,
	}
	// Set the item, unconditionally
	err = memcache.Set(ctx, cacheItem)
	if err != nil {
		log.Debugf(ctx, "Failed setting cache[%v] : %v", cacheKey, err)
	} else {
		// log.Debugf(ctx, "Successfully set cache[%v]", cacheKey)
	}
	return err
}

// Just a shortcut for caching the pair
func (a *MemcacheDatastoreAccessor) cacheZipPair(ctx context.Context, cacheKey string, first interface{}, second interface{}, expiration time.Duration) error {
	pair := &pair{first, second}
	return a.cacheZipValue(ctx, cacheKey, pair, expiration)
}

// A shortcut for caching the zipped (datastoreKey + value)
func (a *MemcacheDatastoreAccessor) cacheZipKeyValue(ctx context.Context, cacheKey string, datastoreKey *datastore.Key, entity interface{}, expiration time.Duration) error {
	kae := &KeyAndEntity{datastoreKey, entity}
	return a.cacheZipValue(ctx, cacheKey, kae, expiration)
}

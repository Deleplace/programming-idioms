package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
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

func (a *MemcacheDatastoreAccessor) cacheValue(c context.Context, cacheKey string, data interface{}, expiration time.Duration) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(&data)
	if err != nil {
		log.Debugf(c, "Failed encoding for cache[%v] : %v", cacheKey, err)
		return err
	}
	cacheItem := &memcache.Item{
		Key:        cacheKey,
		Value:      buffer.Bytes(),
		Expiration: expiration,
	}
	// Set the item, unconditionally
	err = memcache.Set(c, cacheItem)
	if err != nil {
		log.Debugf(c, "Failed setting cache[%v] : %v", cacheKey, err)
	} else {
		// log.Debugf(c, "Successfully set cache[%v]", cacheKey)
	}
	return err
}

func (a *MemcacheDatastoreAccessor) cacheValues(c context.Context, cacheKeys []string, data []interface{}, expiration time.Duration) error {
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
			log.Debugf(c, "Failed encoding for cache[%v] : %v", cacheKey, err)
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
	err := memcache.SetMulti(c, items)
	if err != nil {
		log.Debugf(c, "Failed setting cache items: %v", cacheKeys, err)
	}
	return err
}

func (a *MemcacheDatastoreAccessor) cacheSameValues(c context.Context, cacheKeys []string, data interface{}, expiration time.Duration) error {
	N := len(cacheKeys)

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(&data)

	if err != nil {
		log.Debugf(c, "Failed encoding for cache keys [%v] : %v", cacheKeys, err)
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
	err = memcache.SetMulti(c, items)
	if err != nil {
		log.Debugf(c, "Failed setting cache items: %v", cacheKeys, err)
	}
	return err
}

// A shortcut for caching the datastoreKey + value
func (a *MemcacheDatastoreAccessor) cacheKeyValue(c context.Context, cacheKey string, datastoreKey *datastore.Key, entity interface{}, expiration time.Duration) error {
	kae := &KeyAndEntity{datastoreKey, entity}
	return a.cacheValue(c, cacheKey, kae, expiration)
}

// A shortcut for caching the datastoreKeys + values
func (a *MemcacheDatastoreAccessor) cacheKeysValues(c context.Context, cacheKeys []string, datastoreKeys []*datastore.Key, entities []interface{}, expiration time.Duration) error {
	if len(cacheKeys) != len(datastoreKeys) || len(cacheKeys) != len(entities) {
		panic(fmt.Errorf("Wrong params length %d, %d, %d", len(cacheKeys), len(datastoreKeys), len(entities)))
	}
	N := len(cacheKeys)

	items := make([]interface{}, N)
	for i, datastoreKey := range datastoreKeys {
		kae := &KeyAndEntity{datastoreKey, entities[i]}
		items[i] = kae
	}
	return a.cacheValues(c, cacheKeys, items, expiration)
}

// A shortcut for caching the same datastore key and same value (to encode once) for each cacheKey
func (a *MemcacheDatastoreAccessor) cacheKeysSameValue(c context.Context, cacheKeys []string, datastoreKey *datastore.Key, entity interface{}, expiration time.Duration) error {
	kae := &KeyAndEntity{datastoreKey, entity}
	return a.cacheSameValues(c, cacheKeys, kae, expiration)
}

// Just a shortcut for caching the pair
func (a *MemcacheDatastoreAccessor) cachePair(c context.Context, cacheKey string, first interface{}, second interface{}, expiration time.Duration) error {
	pair := &pair{first, second}
	return a.cacheValue(c, cacheKey, pair, expiration)
}

func (a *MemcacheDatastoreAccessor) readCache(c context.Context, cacheKey string) (interface{}, error) {
	// Get the item from the memcache
	var cacheItem *memcache.Item
	var err error
	if cacheItem, err = memcache.Get(c, cacheKey); err == memcache.ErrCacheMiss {
		// Item not in the cache
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		buffer := bytes.NewBuffer(cacheItem.Value) // todo avoid bytes copy ??
		dec := gob.NewDecoder(buffer)
		var data interface{}
		err = dec.Decode(&data)
		return data, err
	}
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

func (a *MemcacheDatastoreAccessor) recacheIdiom(c context.Context, datastoreKey *datastore.Key, idiom *Idiom, invalidateHTML bool) error {
	cacheKey := fmt.Sprintf("getIdiom(%v)", idiom.Id)
	err := a.cacheKeyValue(c, cacheKey, datastoreKey, idiom, 24*time.Hour)
	if err != nil {
		log.Errorf(c, err.Error())
		return err
	}

	// Batch memcache set call
	N := len(idiom.Implementations)
	cacheKeys := make([]string, N)
	for i, impl := range idiom.Implementations {
		cacheKeys[i] = fmt.Sprintf("getIdiomByImplID(%v)", impl.Id)
	}
	err = a.cacheKeysSameValue(c, cacheKeys, datastoreKey, idiom, 24*time.Hour)
	if err != nil {
		log.Errorf(c, err.Error())
		return err
	}
	// Unfortunately, some previous "getIdiomByImplID(xyz)" might be left uninvalidated.
	// (theoretically)

	if invalidateHTML {
		// Note that cached HTML pages are just evicted, not regenerated here.
		htmlUncacheIdiomAndImpls(c, idiom)
	}

	return err
}

func (a *MemcacheDatastoreAccessor) uncacheIdiom(c context.Context, idiom *Idiom) error {
	cacheKeys := make([]string, 1+len(idiom.Implementations))
	cacheKeys[0] = fmt.Sprintf("getIdiom(%v)", idiom.Id)
	for i, impl := range idiom.Implementations {
		cacheKeys[1+i] = fmt.Sprintf("getIdiomByImplID(%v)", impl.Id)
	}

	err := memcache.DeleteMulti(c, cacheKeys)
	if err != nil {
		log.Errorf(c, err.Error())
	}

	// Cached HTML pages.
	htmlUncacheIdiomAndImpls(c, idiom)

	return err
}

func (a *MemcacheDatastoreAccessor) getIdiom(c context.Context, idiomID int) (*datastore.Key, *Idiom, error) {
	cacheKey := fmt.Sprintf("getIdiom(%v)", idiomID)
	data, cacheerr := a.readCache(c, cacheKey)
	if cacheerr != nil {
		log.Errorf(c, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.getIdiom(c, idiomID)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		key, idiom, err := a.GaeDatastoreAccessor.getIdiom(c, idiomID)
		if err == nil {
			err2 := a.recacheIdiom(c, key, idiom, false)
			logIf(err2, log.Errorf, c, "recaching idiom")
		}
		return key, idiom, err
	}
	// Found in cache :)
	kae := data.(*KeyAndEntity)
	key := kae.Key
	idiom := kae.Entity.(*Idiom)
	return key, idiom, nil
}

func (a *MemcacheDatastoreAccessor) getIdiomByImplID(c context.Context, implID int) (*datastore.Key, *Idiom, error) {
	cacheKey := fmt.Sprintf("getIdiomByImplID(%v)", implID)
	data, cacheerr := a.readCache(c, cacheKey)
	if cacheerr != nil {
		log.Errorf(c, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.getIdiomByImplID(c, implID)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		key, idiom, err := a.GaeDatastoreAccessor.getIdiomByImplID(c, implID)
		if err == nil {
			err2 := a.cacheKeyValue(c, cacheKey, key, idiom, 24*time.Hour)
			logIf(err2, log.Errorf, c, "caching idiom")
		}
		return key, idiom, err
	}
	// Found in cache :)
	kae := data.(*KeyAndEntity)
	key := kae.Key
	idiom := kae.Entity.(*Idiom)
	return key, idiom, nil
}

func (a *MemcacheDatastoreAccessor) saveNewIdiom(c context.Context, idiom *Idiom) (*datastore.Key, error) {
	key, err := a.GaeDatastoreAccessor.saveNewIdiom(c, idiom)
	if err == nil {
		err2 := a.recacheIdiom(c, key, idiom, false)
		logIf(err2, log.Errorf, c, "saving new idiom")
	}
	htmlCacheEvict(c, "about-block-language-coverage")
	return key, err
}

func (a *MemcacheDatastoreAccessor) saveExistingIdiom(c context.Context, key *datastore.Key, idiom *Idiom) error {
	// It is important to invalidate cache with OLD paths, thus before saving
	if _, oldIdiomValue, err := a.getIdiom(c, idiom.Id); err == nil {
		htmlUncacheIdiomAndImpls(c, oldIdiomValue)
	}

	log.Infof(c, "Saving idiom #%v: %v", idiom.Id, idiom.Title)
	err := a.GaeDatastoreAccessor.saveExistingIdiom(c, key, idiom)
	if err == nil {
		log.Infof(c, "Saved idiom #%v, version %v", idiom.Id, idiom.Version)
		err2 := a.recacheIdiom(c, key, idiom, false)
		logIf(err2, log.Errorf, c, "saving existing idiom")
	}
	htmlCacheEvict(c, "about-block-language-coverage")
	return err
}

func (a *MemcacheDatastoreAccessor) stealthIncrementIdiomRating(c context.Context, idiomID int, delta int) (*datastore.Key, *Idiom, error) {
	key, idiom, err := a.GaeDatastoreAccessor.stealthIncrementIdiomRating(c, idiomID, delta)
	err2 := a.recacheIdiom(c, key, idiom, true)
	logIf(err2, log.Errorf, c, "updating idiom rating")
	return key, idiom, err
}

func (a *MemcacheDatastoreAccessor) stealthIncrementImplRating(c context.Context, idiomID, implID int, delta int) (key *datastore.Key, idiom *Idiom, newImplRating int, err error) {
	key, idiom, newImplRating, err = a.GaeDatastoreAccessor.stealthIncrementImplRating(c, idiomID, implID, delta)
	err2 := a.recacheIdiom(c, key, idiom, true)
	logIf(err2, log.Errorf, c, "updating impl rating")
	return
}

func (a *MemcacheDatastoreAccessor) getAllIdioms(c context.Context, limit int, order string) ([]*datastore.Key, []*Idiom, error) {
	cacheKey := fmt.Sprintf("getAllIdioms(%v,%v)", limit, order)
	data, cacheerr := a.readCache(c, cacheKey)
	if cacheerr != nil {
		log.Errorf(c, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.getAllIdioms(c, limit, order)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		keys, idioms, err := a.GaeDatastoreAccessor.getAllIdioms(c, limit, order)
		if err == nil {
			// Cached "All idioms" will have a 10mn lag after an idiom/impl creation.
			//a.cachePair(c, cacheKey, keys, idioms, 10*time.Minute)
			// For now, it might mange too often
			err2 := a.cachePair(c, cacheKey, keys, idioms, 30*time.Second)
			logIf(err2, log.Errorf, c, "caching all idioms")
		}
		return keys, idioms, err
	}
	// Found in cache :)
	pair := data.(*pair)
	keys := pair.First.([]*datastore.Key)
	idioms := pair.Second.([]*Idiom)
	return keys, idioms, nil
}

func (a *MemcacheDatastoreAccessor) deleteAllIdioms(c context.Context) error {
	err := a.GaeDatastoreAccessor.deleteAllIdioms(c)
	if err != nil {
		return err
	}
	// Cache : the nuclear option!
	return memcache.Flush(c)
}

func (a *MemcacheDatastoreAccessor) unindexAll(c context.Context) error {
	return a.GaeDatastoreAccessor.unindexAll(c)
}

func (a *MemcacheDatastoreAccessor) unindex(c context.Context, idiomID int) error {
	return a.GaeDatastoreAccessor.unindex(c, idiomID)
}

func (a *MemcacheDatastoreAccessor) deleteIdiom(c context.Context, idiomID int, why string) error {
	// Clear cache entries
	_, idiom, err := a.GaeDatastoreAccessor.getIdiom(c, idiomID)
	if err == nil {
		err2 := a.uncacheIdiom(c, idiom)
		logIf(err2, log.Errorf, c, "deleting idiom")
	} else {
		log.Errorf(c, "Failed to load idiom %d to uncache: %v", idiomID, err)
	}

	// Delete in datastore
	return a.GaeDatastoreAccessor.deleteIdiom(c, idiomID, why)
}

func (a *MemcacheDatastoreAccessor) deleteImpl(c context.Context, idiomID int, implID int, why string) error {
	// Clear cache entries
	_, idiom, err := a.GaeDatastoreAccessor.getIdiom(c, idiomID)
	if err == nil {
		err2 := a.uncacheIdiom(c, idiom)
		logIf(err2, log.Errorf, c, "deleting impl")
	}

	// Delete in datastore
	err = a.GaeDatastoreAccessor.deleteImpl(c, idiomID, implID, why)
	return err
}

func (a *MemcacheDatastoreAccessor) searchIdiomsByWordsWithFavorites(c context.Context, typedWords, typedLangs []string, favoriteLangs []string, seeNonFavorite bool, limit int) ([]*Idiom, error) {
	// Personalized searches not cached (yet)
	return a.GaeDatastoreAccessor.searchIdiomsByWordsWithFavorites(c, typedWords, typedLangs, favoriteLangs, seeNonFavorite, limit)
}

func (a *MemcacheDatastoreAccessor) searchImplIDs(c context.Context, words, langs []string) (map[string]bool, error) {
	// TODO cache this... or not.
	return a.GaeDatastoreAccessor.searchImplIDs(c, words, langs)
}

func (a *MemcacheDatastoreAccessor) searchIdiomsByLangs(c context.Context, langs []string, limit int) ([]*Idiom, error) {
	cacheKey := fmt.Sprintf("searchIdiomsByLangs(%v,%v)", langs, limit)
	//log.Debugf(c, cacheKey)
	data, cacheerr := a.readCache(c, cacheKey)
	if cacheerr != nil {
		log.Errorf(c, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.searchIdiomsByLangs(c, langs, limit)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		idioms, err := a.GaeDatastoreAccessor.searchIdiomsByLangs(c, langs, limit)
		if err == nil {
			// Search results will have a 10mn lag after an idiom/impl creation/update.
			err2 := a.cacheValue(c, cacheKey, idioms, 10*time.Minute)
			logIf(err2, log.Errorf, c, "caching search results by langs")
		}
		return idioms, err
	}
	// Found in cache :)
	idioms := data.([]*Idiom)
	return idioms, nil
}

func (a *MemcacheDatastoreAccessor) recentIdioms(c context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	cacheKey := fmt.Sprintf("recentIdioms(%v,%v,%v)", favoriteLangs, showOther, n)
	//log.Debugf(c, cacheKey)
	data, cacheerr := a.readCache(c, cacheKey)
	if cacheerr != nil {
		log.Errorf(c, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.recentIdioms(c, favoriteLangs, showOther, n)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		idioms, err := a.GaeDatastoreAccessor.recentIdioms(c, favoriteLangs, showOther, n)
		if err == nil {
			// "Popular idioms" will have a 10mn lag after an idiom/impl creation.
			err2 := a.cacheValue(c, cacheKey, idioms, 10*time.Minute)
			logIf(err2, log.Errorf, c, "caching recent idioms")
		}
		return idioms, err
	}
	// Found in cache :)
	idioms := data.([]*Idiom)
	return idioms, nil
}

func (a *MemcacheDatastoreAccessor) popularIdioms(c context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	cacheKey := fmt.Sprintf("popularIdioms(%v,%v,%v)", favoriteLangs, showOther, n)
	//log.Debugf(c, cacheKey)
	data, cacheerr := a.readCache(c, cacheKey)
	if cacheerr != nil {
		log.Errorf(c, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.popularIdioms(c, favoriteLangs, showOther, n)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		idioms, err := a.GaeDatastoreAccessor.popularIdioms(c, favoriteLangs, showOther, n)
		if err == nil {
			// "Popular idioms" will have a 10mn lag after an idiom/impl creation.
			err2 := a.cacheValue(c, cacheKey, idioms, 10*time.Minute)
			logIf(err2, log.Errorf, c, "caching popular idioms")
		}
		return idioms, err
	}
	// Found in cache :)
	idioms := data.([]*Idiom)
	return idioms, nil
}

// TODO cache this
func (a *MemcacheDatastoreAccessor) idiomsFilterOrder(c context.Context, favoriteLangs []string, limitEachLang int, showOther bool, sortOrder string) ([]*Idiom, error) {
	idioms, err := a.GaeDatastoreAccessor.idiomsFilterOrder(c, favoriteLangs, limitEachLang, showOther, sortOrder)
	return idioms, err
}

func (a *MemcacheDatastoreAccessor) getAppConfig(c context.Context) (ApplicationConfig, error) {
	cacheKey := "getAppConfig()"
	data, cacheerr := a.readCache(c, cacheKey)
	if cacheerr != nil {
		log.Errorf(c, cacheerr.Error())
		// Ouch. Well, skip the cache if it's broken
		return a.GaeDatastoreAccessor.getAppConfig(c)
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		appConfig, err := a.GaeDatastoreAccessor.getAppConfig(c)
		if err == nil {
			log.Infof(c, "Retrieved ApplicationConfig (Toggles) from Datastore")
			err2 := a.cacheValue(c, cacheKey, appConfig, 24*time.Hour)
			logIf(err2, log.Errorf, c, "caching app config")
		}
		return appConfig, err
	}
	// Found in cache :)
	appConfig := data.(*ApplicationConfig)
	return *appConfig, nil
}

func (a *MemcacheDatastoreAccessor) saveAppConfig(c context.Context, appConfig ApplicationConfig) error {
	err := memcache.Flush(c)
	if err != nil {
		return err
	}
	return a.GaeDatastoreAccessor.saveAppConfig(c, appConfig)
	// TODO force toggles refresh for all instances, after memcache flush
}

func (a *MemcacheDatastoreAccessor) saveAppConfigProperty(c context.Context, prop AppConfigProperty) error {
	err := memcache.Flush(c)
	if err != nil {
		return err
	}
	return a.GaeDatastoreAccessor.saveAppConfigProperty(c, prop)
	// TODO force toggles refresh for all instances, after memcache flush
}

func (a *MemcacheDatastoreAccessor) deleteCache(c context.Context) error {
	return memcache.Flush(c)
}

func (a *MemcacheDatastoreAccessor) revert(c context.Context, idiomID int, version int) (*Idiom, error) {
	idiom, err := a.GaeDatastoreAccessor.revert(c, idiomID, version)
	if err != nil {
		return idiom, err
	}
	err2 := a.uncacheIdiom(c, idiom)
	logIf(err2, log.Errorf, c, "uncaching idiom")
	return idiom, err
}

func (a *MemcacheDatastoreAccessor) historyRestore(c context.Context, idiomID int, version int) (*Idiom, error) {

	idiom, err := a.GaeDatastoreAccessor.historyRestore(c, idiomID, version)
	// Uncaching is useful, even when the restore has failed
	errUCI := a.uncacheIdiom(c, idiom)

	if err != nil {
		return idiom, err
	}
	logIf(errUCI, log.Errorf, c, "uncaching idiom")
	return idiom, err
}

func (a *MemcacheDatastoreAccessor) saveNewMessage(c context.Context, msg *MessageForUser) (*datastore.Key, error) {
	key, err := a.GaeDatastoreAccessor.saveNewMessage(c, msg)
	if err != nil {
		return key, err
	}

	cacheKey := "getMessagesForUser(" + msg.Username + ")"
	err = memcache.Delete(c, cacheKey)
	if err == memcache.ErrCacheMiss {
		// No problem if wasn't in cache anyway
		err = nil
	}
	return key, err
}

func (a *MemcacheDatastoreAccessor) getMessagesForUser(c context.Context, username string) ([]*datastore.Key, []*MessageForUser, error) {
	cacheKey := "getMessagesForUser(" + username + ")"

	data, cacheerr := a.readCache(c, cacheKey)
	if cacheerr != nil {
		// Ouch.
		return nil, nil, cacheerr
	}
	if data == nil {
		// Not in the cache. Then fetch the real datastore data. And cache it.
		keys, messages, err := a.GaeDatastoreAccessor.getMessagesForUser(c, username)
		if err == nil {
			err2 := a.cachePair(c, cacheKey, keys, messages, 2*time.Hour)
			logIf(err2, log.Errorf, c, "caching user messages")
		}
		return keys, messages, err
	}
	pair := data.(*pair)
	keys := pair.First.([]*datastore.Key)
	messages := pair.Second.([]*MessageForUser)
	return keys, messages, nil
}

func (a *MemcacheDatastoreAccessor) dismissMessage(c context.Context, key *datastore.Key) (*MessageForUser, error) {
	msg, err := a.GaeDatastoreAccessor.dismissMessage(c, key)
	if err != nil {
		return nil, err
	}

	cacheKey := "getMessagesForUser(" + msg.Username + ")"
	err = memcache.Delete(c, cacheKey)
	if err == memcache.ErrCacheMiss {
		// No problem if wasn't in cache anyway
		err = nil
	}
	return msg, err
}

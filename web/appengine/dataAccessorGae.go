package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	. "github.com/Deleplace/programming-idioms/idioms"

	"context"

	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// GaeDatastoreAccessor is a dataAccessor that works on the Google App Engine Datastore
type GaeDatastoreAccessor struct {
}

var appConfigPropertyNotFound = fmt.Errorf("Found zero AppConfigProperty in the datastore.")

func newIdiomKey(ctx context.Context, idiomID int) *datastore.Key {
	return datastore.NewKey(ctx, "Idiom", "", int64(idiomID), nil)
}

func (a *GaeDatastoreAccessor) getIdiom(ctx context.Context, idiomID int) (*datastore.Key, *Idiom, error) {
	var idiom Idiom
	key := newIdiomKey(ctx, idiomID)
	err := datastore.Get(ctx, key, &idiom)
	return key, &idiom, err
}

func (a *GaeDatastoreAccessor) getIdiomByImplID(ctx context.Context, implID int) (*datastore.Key, *Idiom, error) {
	q := datastore.NewQuery("Idiom").Filter("Implementations.Id =", implID)
	idioms := make([]*Idiom, 0, 1)
	keys, err := q.GetAll(ctx, &idioms)
	if err != nil {
		return nil, nil, err
	}
	if len(idioms) < 1 {
		err = fmt.Errorf("Idiom with implementation id %d not found.", implID)
		return nil, nil, err
	}
	if len(idioms) > 1 {
		err = fmt.Errorf("Multiple Idioms match implementation id %d !", implID)
		return nil, nil, err
	}
	return keys[0], idioms[0], nil
}

func (a *GaeDatastoreAccessor) getIdiomHistory(ctx context.Context, idiomID int, version int) (*datastore.Key, *IdiomHistory, error) {
	q := datastore.NewQuery("IdiomHistory").
		Filter("Id =", idiomID).
		Filter("Version =", version)
	idioms := make([]*IdiomHistory, 0, 1)
	keys, err := q.GetAll(ctx, &idioms)
	if err != nil {
		return nil, nil, err
	}
	if len(idioms) < 1 {
		err = fmt.Errorf("History idiom %d, %d not found.", idiomID, version)
		return nil, nil, err
	}
	if len(idioms) > 1 {
		err = fmt.Errorf("Multiple history idioms match %d, %d !", idiomID, version)
		return nil, nil, err
	}
	return keys[0], idioms[0], nil
}

func (a *GaeDatastoreAccessor) getIdiomHistoryList(ctx context.Context, idiomID int) ([]*datastore.Key, []*IdiomHistory, error) {
	q := datastore.NewQuery("IdiomHistory").
		Project("Version", "VersionDate", "IdiomOrImplLastEditor", "EditSummary").
		Filter("Id =", idiomID).
		Order("-Version")
	historyList := make([]*IdiomHistory, 0)
	keys, err := q.GetAll(ctx, &historyList)
	return keys, historyList, err
}

// Not projected. Retrieves the whole idiom&impl contents for each version of this idiom.
func (a *GaeDatastoreAccessor) getDenseHistoryList(ctx context.Context, idiomID int) ([]*datastore.Key, []*IdiomHistory, error) {
	q := datastore.NewQuery("IdiomHistory").
		Filter("Id =", idiomID).
		Order("-Version")
	historyList := make([]*IdiomHistory, 0)
	keys, err := q.GetAll(ctx, &historyList)
	return keys, historyList, err
}

func (a *GaeDatastoreAccessor) getGlobalHistoryList(ctx context.Context, n int) ([]*datastore.Key, []*IdiomHistory, error) {
	q := datastore.NewQuery("IdiomHistory").
		Project("Id", "Version", "VersionDate", "IdiomOrImplLastEditor", "Title", "EditSummary").
		Order("-VersionDate").
		Limit(n)
	historyList := make([]*IdiomHistory, 0)
	keys, err := q.GetAll(ctx, &historyList)
	return keys, historyList, err
}

// revert modifies Idiom and deletes IdiomHistory, but not in a transaction (for now)
func (a *GaeDatastoreAccessor) revert(ctx context.Context, idiomID int, version int) (*Idiom, error) {
	q := datastore.NewQuery("IdiomHistory").
		Filter("Id =", idiomID).
		Order("-Version").
		Limit(2)
	histories := make([]*IdiomHistory, 0)
	historyKeys, err := q.GetAll(ctx, &histories)
	if err != nil {
		return nil, err
	}
	if len(histories) == 0 {
		return nil, PiErrorf(http.StatusBadRequest, "No history found for idiom %v", idiomID)
	}
	if len(histories) == 1 {
		return nil, PiErrorf(http.StatusBadRequest, "Can't revert the only version of idiom %v", idiomID)
	}
	if histories[0].Version != version {
		return nil, PiErrorf(http.StatusBadRequest, "Can't revert idiom %v: last version is not %v", idiomID, version)
	}
	log.Infof(ctx, "Reverting idiom %v from version %v to version %v", idiomID, histories[0].Version, histories[1].Version)
	idiomKey := newIdiomKey(ctx, idiomID)
	idiom := &histories[1].Idiom
	_, err = datastore.Put(ctx, idiomKey, idiom)
	if err != nil {
		return nil, err
	}
	return idiom, datastore.Delete(ctx, historyKeys[0])
}

func (a *GaeDatastoreAccessor) historyRestore(ctx context.Context, idiomID int, version int, restoreUser string, why string) (*Idiom, error) {
	q := datastore.NewQuery("IdiomHistory").
		Filter("Id =", idiomID).
		Filter("Version =", version).
		Limit(2)
	histories := make([]*IdiomHistory, 0)
	_, err := q.GetAll(ctx, &histories)
	if err != nil {
		return nil, err
	}
	if len(histories) == 0 {
		return nil, PiErrorf(http.StatusBadRequest, "No history found for idiom %v", idiomID)
	}
	var errTooManyItems error
	historyIdiom := &histories[0].Idiom
	if len(histories) >= 2 {
		// Workaround for unsolved bug when history versions are inconsistent
		// Let's just restore the "most recent" candidate
		errTooManyItems = PiErrorf(http.StatusInternalServerError, "Found many history items for idiom %v, version %v. Restoring most recent candidate.", idiomID, version)
		for i := range histories {
			candidate := &histories[i].Idiom
			if candidate.VersionDate.After(historyIdiom.VersionDate) {
				historyIdiom = candidate
			}
		}
	}

	idiomKey, idiom, err := a.getIdiom(ctx, idiomID)
	if err != nil {
		return nil, err
	}
	if idiom.Version == version {
		return nil, PiErrorf(http.StatusBadRequest, "Won't restore idiom %v, version %v to itself.", idiomID, version)
	}
	currentVersion := idiom.Version
	newVersion := idiom.Version + 1
	log.Infof(ctx, "Restoring idiom %v version %v : overwriting version %v, creating new version %v", idiomID, version, currentVersion, newVersion)

	historyIdiom.Version = currentVersion // will be incremented
	historyIdiom.EditSummary = fmt.Sprintf("Restored version %d: %s", version, why)
	historyIdiom.LastEditor = restoreUser
	err = a.saveExistingIdiom(ctx, idiomKey, historyIdiom)
	if err != nil {
		return nil, err
	}
	if errTooManyItems != nil {
		// Return error *after* the restoring is done
		return historyIdiom, errTooManyItems
	}
	return historyIdiom, nil
}

// Delayers registered at init time

// TODO take real Idiom as parameter, not a Key or a pointer
var historyDelayer = delayFunc("save-history-item", func(ctx context.Context, idiomKey *datastore.Key) error {
	var historyItem IdiomHistory
	// TODO check Memcache first
	err := datastore.Get(ctx, idiomKey, &historyItem.Idiom)
	if err != nil {
		return err
	}
	log.Infof(ctx, "Saving history v%d for idiom %d %q", historyItem.Version, historyItem.Idiom.Id, historyItem.Idiom.Title)

	if historyItem.Version <= 1 {
		// This is supposed to be the very first version for this idiom ID.
		// Unless there is some clutter left from a previously delete idiom?
		obsoleteVersionKeys, obsoleteVersions, err := dao.getIdiomHistoryList(ctx, historyItem.Idiom.Id)
		if err != nil {
			log.Errorf(ctx, "checking for obsolete history versions of idiom %d: %v", historyItem.Idiom.Id, err)
		}
		if len(obsoleteVersionKeys) > 0 {
			log.Errorf(ctx, "Found %d obsolete history versions for idiom %d. Will delete them now.", len(obsoleteVersions), historyItem.Idiom.Id)
			if nk, nv := len(obsoleteVersionKeys), len(obsoleteVersions); nk != nv {
				log.Errorf(ctx, "Expected same numbers of Keys and Versions, got %d and %d", nk, nv)
			}
			for i, k := range obsoleteVersionKeys {
				v := obsoleteVersions[i]
				log.Infof(ctx, "Deleting obsolete history version %d (%v) for idiom %d", v.Version, v.VersionDate, historyItem.Idiom.Id)
				err := datastore.Delete(ctx, k)
				if err != nil {
					log.Errorf(ctx, "deleting obsolete history version %d (%v) for idiom %d: %v", v.Version, v.VersionDate, historyItem.Idiom.Id, err)
					// Keep going though
				}
			}
		}
	}

	historyItem.ComputeIdiomOrImplLastEditor()
	// Saves a new IdiomHistory entity. This causes no contention on the original Idiom entity.
	_, err = datastore.Put(ctx, newHistoryKey(ctx), &historyItem)
	return err
})

var indexDelayer = delayFunc("index-text-idiom", func(ctx context.Context, idiomKey *datastore.Key) error {
	var idiom Idiom
	// TODO check Memcache first
	err := datastore.Get(ctx, idiomKey, &idiom)
	if err != nil {
		return err
	}
	// Full text API causes no contention on the original Idiom entity.
	err = indexIdiomFullText(ctx, &idiom, idiomKey)
	if err != nil {
		return err
	}
	err = indexIdiomCheatsheets(ctx, &idiom)
	return err
})

func (a *GaeDatastoreAccessor) saveNewIdiom(ctx context.Context, idiom *Idiom) (*datastore.Key, error) {
	now := time.Now()
	idiom.CreationDate = now
	idiom.Version = 1
	idiom.VersionDate = now
	idiom.ImplCount = len(idiom.Implementations)
	for i := range idiom.Implementations {
		idiom.Implementations[i].CreationDate = now
		idiom.Implementations[i].Version = 1
		idiom.Implementations[i].VersionDate = now
	}

	key, err := datastore.Put(ctx, newIdiomKey(ctx, idiom.Id), idiom)
	if err != nil {
		return key, err
	}

	// Index full-text : asynchronously
	indexDelayer.Call(ctx, key)

	// Save an history item : asynchronously
	// TODO give real Idiom as parameter, not a Key or a pointer
	historyDelayer.Call(ctx, key)

	return key, err
}

func (a *GaeDatastoreAccessor) saveExistingIdiom(ctx context.Context, key *datastore.Key, idiom *Idiom) error {
	idiom.Version = idiom.Version + 1
	idiom.VersionDate = time.Now()
	idiom.ImplCount = len(idiom.Implementations)
	_, err := datastore.Put(ctx, key, idiom)

	// Index full-text : asynchronously
	indexDelayer.Call(ctx, key)

	// Save an history item : asynchronously
	// TODO give real Idiom as parameter, not a Key or a pointer
	historyDelayer.Call(ctx, key)

	return err
}

// stealthIncrementIdiomRating doesn't update Version and VersionDate
func (a *GaeDatastoreAccessor) stealthIncrementIdiomRating(ctx context.Context, idiomID int, delta int) (*datastore.Key, *Idiom, error) {
	key, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return nil, nil, err
	}

	idiom.Rating += delta

	_, err = datastore.Put(ctx, key, idiom)
	return key, idiom, err
}

// stealthIncrementImplRating doesn't update Version and VersionDate
func (a *GaeDatastoreAccessor) stealthIncrementImplRating(ctx context.Context, idiomID, implID int, delta int) (key *datastore.Key, idiom *Idiom, newImplRating int, err error) {
	key, idiom, err = dao.getIdiom(ctx, idiomID)
	if err != nil {
		return nil, nil, 0, err
	}

	// TODO: more efficient way than iterating?
	_, impl, _ := idiom.FindImplInIdiom(implID)
	impl.Rating += delta

	_, err = datastore.Put(ctx, key, idiom)
	return key, idiom, impl.Rating, err
}

func newHistoryKey(ctx context.Context) *datastore.Key {
	return datastore.NewIncompleteKey(ctx, "IdiomHistory", nil)
}

// implementedLanguages may return duplicates, which is ok
func implementedLanguagesConcat(idiom *Idiom) string {
	langs := ""
	for _, impl := range idiom.Implementations {
		langs += impl.LanguageName + " "
	}
	// TODO add non-canonical synonyms...
	return langs
}

func (a *GaeDatastoreAccessor) getAllIdioms(ctx context.Context, limit int, order string) ([]*datastore.Key, []*Idiom, error) {
	q := datastore.NewQuery("Idiom")
	if order != "" {
		q = q.Order(order)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	idioms := make([]*Idiom, 0, 500)
	keys, err := q.GetAll(ctx, &idioms)
	return keys, idioms, err
}

func (a *GaeDatastoreAccessor) deleteAllIdioms(ctx context.Context) error {
	keys, err := datastore.NewQuery("Idiom").KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return err
	}

	err = a.unindexAll(ctx)
	if err != nil {
		return err
	}

	return datastore.DeleteMulti(ctx, keys)
}

func (a *GaeDatastoreAccessor) deleteIdiom(ctx context.Context, idiomID int, why string) error {
	key, idiom, err := a.getIdiom(ctx, idiomID)
	if err != nil {
		return err
	}
	// Remove from text search indexes
	err = a.unindex(ctx, idiom)
	if err != nil {
		log.Errorf(ctx, "Failed to unindex idiom %d: %v", idiomID, err)
	}
	return datastore.Delete(ctx, key)
	// The why param is ignored for now, because idiom doesn't exist anymore.
}

func (a *GaeDatastoreAccessor) deleteImpl(ctx context.Context, idiomID int, implID int, why string) error {
	key, idiom, err := a.getIdiom(ctx, idiomID)
	if err != nil {
		return err
	}
	idiom.EditSummary = why
	if i, _, found := idiom.FindImplInIdiom(implID); found {
		idiom.Implementations = append(idiom.Implementations[:i], idiom.Implementations[i+1:]...)
		return a.saveExistingIdiom(ctx, key, idiom)
	}
	return fmt.Errorf("Could not find impl %v in idiom %v", idiom.Id, implID)
}

func (a *GaeDatastoreAccessor) processUploadFile(r *http.Request, name string) (string, map[string][]string, error) {
	// See https://developers.google.com/appengine/docs/go/blobstore/#Uploading_a_Blob
	blobs, otherParams, err := blobstore.ParseUpload(r)
	if err != nil {
		return "", nil, err
	}
	file := blobs[name]
	if len(file) == 0 {
		return "", otherParams, nil
	}
	return string(file[0].BlobKey), otherParams, nil
}

func (a *GaeDatastoreAccessor) processUploadFiles(r *http.Request, names []string) ([]string, map[string][]string, error) {
	blobs, otherParams, err := blobstore.ParseUpload(r)
	if err != nil {
		return nil, nil, err
	}
	blobKeys := []string{}
	for _, name := range names {
		if file := blobs[name]; len(file) > 0 {
			blobKeys = append(blobKeys, string(file[0].BlobKey))
		}
	}
	return blobKeys, otherParams, nil
}

func (a *GaeDatastoreAccessor) nextIdiomID(ctx context.Context) (int, error) {
	q := datastore.NewQuery("Idiom").Order("-Id"). /*.Project("Id")*/ Limit(1)
	it := q.Run(ctx)
	var maxIdiom Idiom
	_, err := it.Next(&maxIdiom)
	if err == datastore.Done {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	newID := maxIdiom.Id + 1
	return newID, nil
}

func (a *GaeDatastoreAccessor) nextImplID(ctx context.Context) (int, error) {
	// This is not scalable and may yield the same id twice.
	// TODO cleanup this mess.
	// order by implId desc : is it still ok with multi-valued implId ...?
	q := datastore.NewQuery("Idiom").Order("-Implementations.Id"). /*.Project("Implementations.Id")*/ Limit(1)
	it := q.Run(ctx)
	var maxIdiom Idiom
	_, err := it.Next(&maxIdiom)
	if err == datastore.Done {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	if len(maxIdiom.Implementations) == 0 {
		return 0, fmt.Errorf("Existing idiom %d should not have zero impl", maxIdiom.Id)
	}
	maxImplID := -1
	for j := range maxIdiom.Implementations {
		if maxIdiom.Implementations[j].Id > maxImplID {
			maxImplID = maxIdiom.Implementations[j].Id
		}
	}
	newID := maxImplID + 1

	if _, _, err := a.getIdiomByImplID(ctx, newID); err == nil {
		return 0, fmt.Errorf("Impl %d already exists :(", newID)
	}
	return newID, nil
}

func (a *GaeDatastoreAccessor) recentIdioms(ctx context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	idioms, err := a.idiomsFilterOrder(ctx, favoriteLangs, n, showOther, "-VersionDate")
	if err != nil {
		return idioms, err
	}
	sortIdiomsByVersionDate(idioms)
	if len(idioms) > n {
		idioms = idioms[0:n]
	}
	return idioms, err
}

func (a *GaeDatastoreAccessor) popularIdioms(ctx context.Context, favoriteLangs []string, showOther bool, n int) ([]*Idiom, error) {
	idioms, err := a.idiomsFilterOrder(ctx, favoriteLangs, n, showOther, "-Rating")
	if err != nil {
		return idioms, err
	}
	sortIdiomsByRating(idioms)
	if len(idioms) > n {
		idioms = idioms[0:n]
	}
	return idioms, err
}

// Makes one datastore Query for each favorite language with specified sortOrder, then one Query without a language filter.
// Then concatenates the results (eliminating duplicates).
func (a *GaeDatastoreAccessor) idiomsFilterOrder(ctx context.Context, favoriteLangs []string, limitEachLang int, showOther bool, sortOrder string) ([]*Idiom, error) {
	idiomsResult := make([]*Idiom, 0, limitEachLang*len(favoriteLangs))

	langFilters := make([]string, len(favoriteLangs))
	copy(langFilters, favoriteLangs)
	if showOther {
		langFilters = append(langFilters, "") // 1 extra dummy for "no filter"
	}

	idSet := map[int]struct{}{} // To evict duplicates

	for _, lg := range langFilters {
		q := datastore.NewQuery("Idiom")
		if lg != "" {
			q = q.Filter("Implementations.LanguageName =", lg)
		}
		q = q.Order(sortOrder).Order("-Id").Limit(limitEachLang)
		idioms := make([]*Idiom, 0, limitEachLang)
		if _, err := q.GetAll(ctx, &idioms); err != nil {
			return nil, err
		}
		for _, idiom := range idioms {
			if _, seen := idSet[idiom.Id]; !seen {
				idiomsResult = append(idiomsResult, idiom)
				idSet[idiom.Id] = struct{}{}
			}
		}
	}

	for _, idiom := range idiomsResult {
		seeNonFavorite := true // TODO extract from soft profile!!
		// Inside each Idiom, sort Implementations according to favorites
		implFavoriteLanguagesFirstWithOrder(idiom, favoriteLangs, "", seeNonFavorite)
	}

	return idiomsResult, nil
}

func (a *GaeDatastoreAccessor) randomIdiom(ctx context.Context) (*datastore.Key, *Idiom, error) {
	q := datastore.NewQuery("Idiom")
	//q := q.KeysOnly()
	//keys, err := q.GetAll(ctx, nil)
	count, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	k := rand.Intn(count)
	// This is really slow: ~100ms. TODO find a better way.
	q = q.Offset(k).Limit(1)
	idioms := make([]*Idiom, 0, 1)
	keys, err := q.GetAll(ctx, &idioms)
	if err != nil {
		return nil, nil, err
	}
	return keys[0], idioms[0], err
}

// Similar to randomIdiom, but with a lang filter.
func (a *GaeDatastoreAccessor) randomIdiomHaving(ctx context.Context, havingLang string) (*datastore.Key, *Idiom, error) {
	q := datastore.NewQuery("Idiom")
	q = q.Filter("Implementations.LanguageName =", havingLang)
	count, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	log.Debugf(ctx, "Found %d idioms having %s", count, havingLang)
	if count == 0 {
		return nil, nil, fmt.Errorf("No implementations found in language [%s]", havingLang)
	}
	k := rand.Intn(count)
	q = q.Offset(k).Limit(1)
	idioms := make([]*Idiom, 0, 1)
	keys, err := q.GetAll(ctx, &idioms)
	if err != nil {
		return nil, nil, err
	}
	if len(keys) == 0 {
		return nil, nil, fmt.Errorf("No idiom found for lang %s :|", havingLang)
	}
	log.Debugf(ctx, "Fetched 1 idiom having %s", havingLang)
	return keys[0], idioms[0], err
}

// randomIdiomNotHaving uses big lists of keys, because the Datastore
// doesn't handle natively the query "All idioms not having this language".
func (a *GaeDatastoreAccessor) randomIdiomNotHaving(ctx context.Context, notHavingLang string) (*datastore.Key, *Idiom, error) {
	// All keys
	q1 := datastore.NewQuery("Idiom")
	q1 = q1.KeysOnly()
	keys1, err := q1.GetAll(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Keys of idioms having this lang
	q2 := datastore.NewQuery("Idiom")
	q2 = q2.Filter("Implementations.LanguageName =", notHavingLang)
	q2 = q2.KeysOnly()
	keys2, err := q2.GetAll(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	keySet2 := make(map[datastore.Key]bool, len(keys2))
	for _, key2 := range keys2 {
		keySet2[*key2] = true
	}

	// Difference = keys of idioms not having this lang
	keys3 := make([]*datastore.Key, 0, len(keys1)-len(keys2))
	for _, key1 := range keys1 {
		if !keySet2[*key1] {
			keys3 = append(keys3, key1)
		}
	}
	count := len(keys3)

	if count == 0 {
		return nil, nil, PiErrorf(http.StatusInternalServerError, "%v contributors are so effective, that no unimplemented idiom could be found :|", notHavingLang)
	}

	k := rand.Intn(count)
	key := keys3[k]

	var idiom Idiom
	err = datastore.Get(ctx, key, &idiom)
	return key, &idiom, err
}

// AppConfigProperty is a (global) application property
type AppConfigProperty struct {
	AppConfigId int
	Name        string
	Value       bool
}

func (a *GaeDatastoreAccessor) getAppConfig(ctx context.Context) (ApplicationConfig, error) {
	q := datastore.NewQuery("AppConfigProperty") // TODO .Filter("AppConfigId =", appConfigId)
	properties := make([]*AppConfigProperty, 0, 100)
	_, err := q.GetAll(ctx, &properties)
	if err != nil {
		return ApplicationConfig{}, err
	}
	if len(properties) == 0 {
		return ApplicationConfig{}, appConfigPropertyNotFound
	}

	appConfig := ApplicationConfig{
		Id:      0, // TODO meaningful appConfigId
		Toggles: make(Toggles, len(properties)),
	}
	for _, prop := range properties {
		appConfig.Toggles[prop.Name] = prop.Value
	}
	return appConfig, nil
}

func (a *GaeDatastoreAccessor) saveAppConfig(ctx context.Context, appConfig ApplicationConfig) error {
	keys := make([]*datastore.Key, len(appConfig.Toggles))
	properties := make([]*AppConfigProperty, len(appConfig.Toggles))
	i := 0
	for name, value := range appConfig.Toggles {
		prop := AppConfigProperty{
			AppConfigId: 0, // TODO: meaningful appConfigId
			Name:        name,
			Value:       value,
		}
		keystr := fmt.Sprintf("%d_%s", prop.AppConfigId, prop.Name)
		keys[i] = datastore.NewKey(ctx, "AppConfigProperty", keystr, 0, nil)
		properties[i] = &prop
		i++
	}
	_, err := datastore.PutMulti(ctx, keys, properties)
	return err
}

func (a *GaeDatastoreAccessor) saveAppConfigProperty(ctx context.Context, prop AppConfigProperty) error {
	keystr := fmt.Sprintf("%d_%s", prop.AppConfigId, prop.Name)
	key := datastore.NewKey(ctx, "AppConfigProperty", keystr, 0, nil)
	_, err := datastore.Put(ctx, key, &prop)
	return err
}

func (a *GaeDatastoreAccessor) saveNewMessage(ctx context.Context, message *MessageForUser) (*datastore.Key, error) {
	return datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "MessageForUser", nil), message)
}

func (a *GaeDatastoreAccessor) getMessagesForUser(ctx context.Context, username string) ([]*datastore.Key, []*MessageForUser, error) {
	var dateZero time.Time
	q := datastore.NewQuery("MessageForUser").
		Filter("Username =", username).
		Filter("DismissalDate =", dateZero)
	messages := make([]*MessageForUser, 0)
	keys, err := q.GetAll(ctx, &messages)

	// Mark as seen
	now := time.Now()
	for _, msg := range messages {
		msg.LastViewDate = now
		if msg.FirstViewDate == dateZero {
			msg.FirstViewDate = now
		}
	}
	_, err = datastore.PutMulti(ctx, keys, messages)
	if err != nil {
		log.Warningf(ctx, "Could not save messages view dates: %v", err)
	}

	return keys, messages, err
}

func (a *GaeDatastoreAccessor) dismissMessage(ctx context.Context, key *datastore.Key) (*MessageForUser, error) {
	var userMessage MessageForUser
	err := datastore.Get(ctx, key, &userMessage)
	if err != nil {
		return nil, err
	}
	userMessage.DismissalDate = time.Now()
	_, err = datastore.Put(ctx, key, &userMessage)
	return &userMessage, err
}

func (a *GaeDatastoreAccessor) getAllIdiomTitles(ctx context.Context) ([]*Idiom, error) {
	q := datastore.NewQuery("Idiom").Project("Id", "Title")
	idioms := make([]*Idiom, 0, 10)
	_, err := q.GetAll(ctx, &idioms)
	return idioms, err
}

func (a *GaeDatastoreAccessor) getMaxIdiomIDTitle(ctx context.Context) (*Idiom, error) {
	// Used for #192 Keyboard shortcut 'p'
	q := datastore.NewQuery("Idiom").Order("-Id").Project("Id", "Title").Limit(1)
	idioms := make([]*Idiom, 0, 1)
	_, err := q.GetAll(ctx, &idioms)
	if err != nil {
		return nil, err
	}
	return idioms[0], err
}

package main

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/idioms"

	"context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

// Low-level Datastore entities manipulation, outside
// the scope of a normal request.
// Useful for patches or migration.

func adminResaveEntities(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var err error
	switch r.FormValue("kind") {
	case "IdiomHistory":
		err = resaveAllIdiomHistory(ctx)
	default:
		return PiErrorf(http.StatusBadRequest, "Wrong kind [%s]", r.FormValue("kind"))
	}

	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Done.")
	return nil
}

// 2015-11-06 to force field EditSummary (even if empty) on every IdiomHistory persisted entity.
func resaveAllIdiomHistory(ctx context.Context) error {
	defer memcache.Flush(ctx)
	saved := 0
	keys, err := datastore.NewQuery("IdiomHistory").KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return err
	}
	nbEntities := len(keys)

	defer func() {
		log.Infof(ctx, "Resaved %d IdiomHistory entities out of %d.", saved, nbEntities)
	}()

	for len(keys) > 0 {
		bunch := 100
		if len(keys) < bunch {
			bunch = len(keys)
		}
		histories := make([]*IdiomHistory, bunch)
		err := datastore.GetMulti(ctx, keys[:bunch], histories)
		if err != nil {
			return err
		}
		_, err = datastore.PutMulti(ctx, keys[:bunch], histories)
		if err != nil {
			return err
		}
		saved += bunch

		// Remove processed keys
		keys = keys[bunch:]
	}
	return nil
}

func adminRepairHistoryVersions(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	defer memcache.Flush(ctx)

	idiomIDStr := r.FormValue("idiomId")
	if idiomIDStr == "" {
		return PiErrorf(http.StatusBadRequest, "Mandatory param: idiomId")
	}
	idiomID := String2Int(idiomIDStr)

	// Warning: fetching the whole history of 1 idiom
	// may have quite a big memory footprint
	log.Infof(ctx, "Repairing versions for idiom: %v", idiomID)

	q := datastore.NewQuery("IdiomHistory").
		Filter("Id =", idiomID).
		Order("VersionDate")
	histories := make([]*IdiomHistory, 0)
	historyKeys, err := q.GetAll(ctx, &histories)
	if err != nil {
		return err
	}
	for i := range histories[1:] {
		if histories[i].VersionDate.After(histories[i+1].VersionDate) {
			return PiErrorf(http.StatusInternalServerError, "History items not well sorted")
		}
	}

	for i := range histories {
		histories[i].Version = 1 + i
	}
	lastVersion := len(histories)
	log.Infof(ctx, "\tSaving %v history entities.", len(histories))
	for len(historyKeys) > 0 {
		bunch := 10
		if len(historyKeys) < 10 {
			bunch = len(historyKeys)
		}
		_, err = datastore.PutMulti(ctx, historyKeys[:bunch], histories[:bunch])
		if err != nil {
			return err
		}
		// Remove processed items
		historyKeys = historyKeys[bunch:]
		histories = histories[bunch:]
	}

	var idiom Idiom
	idiomKey := newIdiomKey(ctx, idiomID)
	err = datastore.Get(ctx, idiomKey, &idiom)
	if err != nil {
		return err
	}
	if idiom.Version == lastVersion {
		log.Infof(ctx, "\tIdiom version %v already clean", idiom.Version)
	} else {
		log.Infof(ctx, "\tFixing idiom version %v -> %v", idiom.Version, lastVersion)
		idiom.Version = lastVersion
		_, err = datastore.Put(ctx, idiomKey, &idiom)
		if err != nil {
			return err
		}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, Response{"success": true, "message": "History repaired for idiom " + idiomIDStr})
	return nil
}

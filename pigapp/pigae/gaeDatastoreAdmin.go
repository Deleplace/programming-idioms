package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

// Low-level Datastore entities manipulation, outside
// the scope of a normal request.
// Useful for patches or migration.

func adminResaveEntities(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	var err error
	switch r.FormValue("kind") {
	case "IdiomHistory":
		err = resaveAllIdiomHistory(c)
	default:
		return PiError{ErrorText: "Wrong kind [" + r.FormValue("kind") + "]"}
	}

	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Done.")
	return nil
}

// 2015-11-06 to force field EditSummary (even if empty) on every IdiomHistory persisted entity.
func resaveAllIdiomHistory(c appengine.Context) error {
	defer memcache.Flush(c)
	saved := 0
	keys, err := datastore.NewQuery("IdiomHistory").KeysOnly().GetAll(c, nil)
	if err != nil {
		return err
	}
	nbEntities := len(keys)

	defer func() {
		c.Infof("Resaved %d IdiomHistory entities out of %d.", saved, nbEntities)
	}()

	for len(keys) > 0 {
		bunch := 100
		if len(keys) < bunch {
			bunch = len(keys)
		}
		histories := make([]*IdiomHistory, bunch)
		err := datastore.GetMulti(c, keys[:bunch], histories)
		if err != nil {
			return err
		}
		_, err = datastore.PutMulti(c, keys[:bunch], histories)
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
	c := appengine.NewContext(r)
	defer memcache.Flush(c)

	// TODO turn this into small delayed tasks!
	idiomKeys, err := datastore.NewQuery("Idiom").KeysOnly().GetAll(c, nil)
	if err != nil {
		return err
	}
	for _, idiomKey := range idiomKeys {
		// Warning: fetching the whole history of 1 idiom
		// may have quite a big memory footprint
		idiomID := idiomKey.IntID()
		c.Infof("Repairing versions for idiom: %v", idiomID)

		q := datastore.NewQuery("IdiomHistory").
			Filter("Id =", idiomID).
			Order("VersionDate")
		histories := make([]*IdiomHistory, 0)
		historyKeys, err := q.GetAll(c, &histories)
		if err != nil {
			return err
		}
		for i := range histories[1:] {
			if histories[i].VersionDate.After(histories[i+1].VersionDate) {
				return PiError{ErrorText: "History items not well sorted", Code: 500}
			}
		}

		for i := range histories {
			histories[i].Version = 1 + i
		}
		lastVersion := len(histories)
		c.Infof("\tSaving %v history entities.", len(histories))
		for len(historyKeys) > 0 {
			bunch := 10
			if len(historyKeys) < 10 {
				bunch = len(historyKeys)
			}
			_, err = datastore.PutMulti(c, historyKeys[:bunch], histories[:bunch])
			if err != nil {
				return err
			}
			// Remove processed items
			historyKeys = historyKeys[bunch:]
			histories = histories[bunch:]
		}

		var idiom Idiom
		err = datastore.Get(c, idiomKey, &idiom)
		if err != nil {
			return err
		}
		if idiom.Version == lastVersion {
			c.Infof("\tIdiom version %v already clean", idiom.Version)
		} else {
			c.Infof("\tFixing idiom version %v -> %v", idiom.Version, lastVersion)
			idiom.Version = lastVersion
			_, err = datastore.Put(c, idiomKey, &idiom)
			if err != nil {
				return err
			}
		}
	}

	fmt.Fprintln(w, "Done.")
	return nil
}

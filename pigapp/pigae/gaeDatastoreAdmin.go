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
		if len(keys) < 100 {
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

package pigae

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func adminReindexAjax(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	err := dao.deleteCache(c)
	if err != nil {
		log.Warningf(c, "Problem deleting cache: %v", err.Error())
	}
	err = dao.unindexAll(c)
	if err != nil {
		log.Warningf(c, "Problem deleting cache: %v", err.Error())
	}
	limit := 10000 // TODO chunk that??
	keys, idioms, err := dao.getAllIdioms(c, limit, "")
	if err != nil {
		return err
	}
	log.Infof(c, "Reindexing %d idioms...", len(keys))
	indexed := 0
	for i := range keys {
		key := keys[i]
		idiom := idioms[i]
		err := indexIdiomFullText(c, idiom, key)
		if err != nil {
			return nil
		}
		indexed++
	}
	log.Infof(c, "Reindexed %d idioms.", indexed)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"indexed": indexed})
	return nil
}

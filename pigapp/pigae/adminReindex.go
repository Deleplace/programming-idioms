package pigae

import (
	"fmt"
	"net/http"

	"appengine"
)

func adminReindexAjax(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	dao.deleteCache(c)
	dao.unindexAll(c)
	limit := 10000 // TODO chunk that??
	keys, idioms, err := dao.getAllIdioms(c, limit, "")
	if err != nil {
		return err
	}
	c.Infof("Reindexing %d idioms...", len(keys))
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
	c.Infof("Reindexed %d idioms.", indexed)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"indexed": indexed})
	return nil
}

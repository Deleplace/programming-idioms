package pigae

import (
	"fmt"
	"net/http"

	"appengine"
)

func adminReindexAjax(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	dao.deleteCache(c)
	limit := 10000 // TODO chunk that??
	keys, idioms, err := dao.getAllIdioms(c, limit, "")
	if err != nil {
		return err
	}
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

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"indexed": indexed})
	return nil
}

package pigae

import (
	"fmt"
	"net/http"
	"sync"

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

	// Get rid of deprecated Idiom.Words and Idiom.WordsTitle
	var wg sync.WaitGroup
	for i := range keys {
		key := keys[i]
		idiom := idioms[i]
		// 1 at a time = not very efficient, but who cares
		if len(idiom.Words) > 0 || len(idiom.Words) > 0 {
			idiom.Words = nil
			idiom.WordsTitle = nil
			wg.Add(1)
			go func() {
				dao.saveExistingIdiom(c, key, idiom)
				wg.Done()
			}()
		}
	}
	wg.Wait()
	return nil
}

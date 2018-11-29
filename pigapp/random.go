package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"
	"github.com/gorilla/mux"
)

func randomIdiom(w http.ResponseWriter, r *http.Request) error {
	c := r.Context()

	var urls []string
	cachedUrls, err := dao.readCache(c, "all-idioms-urls")
	if err == nil && cachedUrls != nil {
		urls = cachedUrls.([]string)
	} else {
		// Not found (or cache failed)
		// Let's fetch in the Datastore
		infof(c, "Fetching all idiom titles from Datastore")
		idiomHeads, err := dao.getAllIdiomTitles(c)
		if err != nil {
			return err
		}
		urls = make([]string, len(idiomHeads))
		for i, head := range idiomHeads {
			urls[i] = NiceIdiomRelativeURL(head)
		}
		err = dao.cacheValue(c, "all-idioms-urls", urls, 24*time.Hour)
		if err != nil {
			errorf(c, "Failed idioms URLs list in cache: %v", err)
		}
	}
	k := rand.Intn(len(urls))
	url := urls[k]
	// Note that we're redirecting to a *relative* URL
	infof(c, "Picked idiom url %s (out of %d)", url, len(urls))

	// TODO when AppEngine has Go1.8:
	// if pusher, ok := w.(http.Pusher); ok {
	// 	if err := pusher.Push(url, nil); err != nil {
	// 		errorf("Failed to push %s: %v", url, err)
	// 	}
	// }

	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

// Among idioms having an impl in this language
func randomIdiomHaving(w http.ResponseWriter, r *http.Request) error {
	c := r.Context()
	vars := mux.Vars(r)

	havingLang := vars["havingLang"]
	havingLang = NormLang(havingLang)
	if havingLang == "" {
		return fmt.Errorf("Invalid language [%s]", vars["havingLang"])
	}

	var idiom *Idiom
	var url string
	var err error

	infof(c, "Going to a random idiom having lang %v", havingLang)
	_, idiom, err = dao.randomIdiomHaving(c, havingLang)
	if err != nil {
		return err
	}
	url = NiceIdiomURL(idiom)
	for _, impl := range idiom.Implementations {
		if impl.LanguageName == havingLang {
			url = NiceImplURL(idiom, impl.Id, havingLang)
		}
	}

	infof(c, "Picked idiom #%v: %v", idiom.Id, idiom.Title)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

// Among idioms not having an impl in this language
func randomIdiomNotHaving(w http.ResponseWriter, r *http.Request) error {
	c := r.Context()
	vars := mux.Vars(r)

	notHavingLang := vars["notHavingLang"]
	notHavingLang = NormLang(notHavingLang)
	if notHavingLang == "" {
		return fmt.Errorf("Invalid language [%s]", vars["notHavingLang"])
	}

	var idiom *Idiom
	var url string
	var err error

	infof(c, "Going to a random idiom having lang %v", notHavingLang)
	_, idiom, err = dao.randomIdiomNotHaving(c, notHavingLang)
	if err != nil {
		return err
	}
	url = NiceIdiomURL(idiom)

	infof(c, "Picked idiom #%v: %v", idiom.Id, idiom.Title)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

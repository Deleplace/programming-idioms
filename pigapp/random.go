package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"
	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func randomIdiom(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	var urls []string
	cachedUrls, err := dao.readCache(c, "all-idioms-urls")
	if err == nil && cachedUrls != nil {
		urls = cachedUrls.([]string)
	} else {
		// Not found (or Memcache failed)
		// Let's fetch in the Datastore
		log.Infof(c, "Fetching all idiom titles from Datastore")
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
			log.Errorf(c, "Failed idioms URLs list in Memcache: %v", err)
		}
	}
	k := rand.Intn(len(urls))
	url := urls[k]
	// Note that we're redirecting to a *relative* URL
	log.Infof(c, "Picked idiom url %s (out of %d)", url, len(urls))

	// TODO when AppEngine has Go1.8:
	// if pusher, ok := w.(http.Pusher); ok {
	// 	if err := pusher.Push(url, nil); err != nil {
	// 		log.Errorf("Failed to push %s: %v", url, err)
	// 	}
	// }

	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

// Among idioms having an impl in this language
func randomIdiomHaving(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)

	havingLang := vars["havingLang"]
	havingLang = NormLang(havingLang)
	if havingLang == "" {
		return fmt.Errorf("Invalid language [%s]", vars["havingLang"])
	}

	var idiom *Idiom
	var url string
	var err error

	log.Infof(c, "Going to a random idiom having lang %v", havingLang)
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

	log.Infof(c, "Picked idiom #%v: %v", idiom.Id, idiom.Title)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

// Among idioms not having an impl in this language
func randomIdiomNotHaving(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)

	notHavingLang := vars["notHavingLang"]
	notHavingLang = NormLang(notHavingLang)
	if notHavingLang == "" {
		return fmt.Errorf("Invalid language [%s]", vars["notHavingLang"])
	}

	var idiom *Idiom
	var url string
	var err error

	log.Infof(c, "Going to a random idiom having lang %v", notHavingLang)
	_, idiom, err = dao.randomIdiomNotHaving(c, notHavingLang)
	if err != nil {
		return err
	}
	url = NiceIdiomURL(idiom)

	log.Infof(c, "Picked idiom #%v: %v", idiom.Id, idiom.Title)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

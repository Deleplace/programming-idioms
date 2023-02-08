package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	. "github.com/Deleplace/programming-idioms/idioms"
	"github.com/gorilla/mux"

	"google.golang.org/appengine/v2/log"
)

func randomIdiom(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var urls []string
	cachedUrls, err := dao.readCache(ctx, "all-idioms-urls")
	if err == nil && cachedUrls != nil {
		urls = cachedUrls.([]string)
	} else {
		// Not found (or Memcache failed)
		// Let's fetch in the Datastore
		log.Infof(ctx, "Fetching all idiom titles from Datastore")
		idiomHeads, err := dao.getAllIdiomTitles(ctx)
		if err != nil {
			return err
		}
		urls = make([]string, len(idiomHeads))
		for i, head := range idiomHeads {
			urls[i] = NiceIdiomRelativeURL(head)
		}
		err = dao.cacheValue(ctx, "all-idioms-urls", urls, 24*time.Hour)
		if err != nil {
			log.Errorf(ctx, "Failed idioms URLs list in Memcache: %v", err)
		}
	}
	if len(urls) == 0 {
		return errors.New("There are no idioms in the database, yet")
	}
	k := rand.Intn(len(urls))
	url := urls[k]
	// Note that we're redirecting to a *relative* URL
	log.Infof(ctx, "Picked idiom url %s (out of %d)", url, len(urls))

	// 2018-09 w doesn't seem to implement Pusher :(
	// if pusher, ok := w.(http.Pusher); ok {
	// 	if err := pusher.Push(url, nil); err != nil {
	// 		log.Errorf("Failed to push %s: %v", url, err)
	// 	}
	// }
	// Automagic "Link" header seems fine, thanks to GFE
	w.Header().Set("Link", "<"+url+">; rel=preload; as=document")

	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

// Among idioms having an impl in this language
func randomIdiomHaving(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)

	havingLang := vars["havingLang"]
	havingLang = NormLang(havingLang)
	if havingLang == "" {
		return fmt.Errorf("Invalid language [%s]", vars["havingLang"])
	}

	var idiom *Idiom
	var url string
	var err error

	log.Infof(ctx, "Going to a random idiom having lang %v", havingLang)
	_, idiom, err = dao.randomIdiomHaving(ctx, havingLang)
	if err != nil {
		return err
	}
	url = NiceIdiomURL(idiom)
	for _, impl := range idiom.Implementations {
		if impl.LanguageName == havingLang {
			url = NiceImplURL(idiom, impl.Id, havingLang)
		}
	}

	log.Infof(ctx, "Picked idiom #%v: %v", idiom.Id, idiom.Title)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

// Among idioms not having an impl in this language
func randomIdiomNotHaving(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)

	notHavingLang := vars["notHavingLang"]
	notHavingLang = NormLang(notHavingLang)
	if notHavingLang == "" {
		return fmt.Errorf("Invalid language [%s]", vars["notHavingLang"])
	}

	var idiom *Idiom
	var url string
	var err error

	log.Infof(ctx, "Going to a random idiom having lang %v", notHavingLang)
	_, idiom, err = dao.randomIdiomNotHaving(ctx, notHavingLang)
	if err != nil {
		return err
	}
	url = NiceIdiomURL(idiom)

	log.Infof(ctx, "Picked idiom #%v: %v", idiom.Id, idiom.Title)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

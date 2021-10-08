package main

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"
	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"
)

// Language Backlog: community contribution nudge.
// Shows a few stats about the Coverage for this lang,
// and links to idioms/impls to be improved.

// 4 sections will show sampleSize lines each: Missing impl, Missing doc, Missing demo, Curation.
const sampleSize = 3

func backlogForLanguage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	rawLang := vars["lang"]
	lang := NormLang(rawLang)
	log.Infof(ctx, "Computing backlog for %s", lang)

	data := &BacklogLanguageFacade{
		PageMeta: PageMeta{
			PageTitle: "Idioms missing data for the " + lang + " implementation",
			Toggles:   toggles,
			ExtraCss:  []string{hostPrefix() + themeDirectory() + "/css/pages/backlog.css"},
			ExtraJs:   []string{hostPrefix() + themeDirectory() + "/js/pages/backlog.js"},
		},
		UserProfile:         readUserProfile(r),
		Lang:                lang,
		RecommendedDemoSite: recommendedDemoSite(rawLang),
	}

	log.Infof(ctx, "searchRandomImplsForLang(%q)...", rawLang)
	tip := time.Now()
	var err error
	data.CurationSuggestions, err = searchRandomImplsForLang(ctx, rawLang, sampleSize)
	log.Infof(ctx, "got %d curation suggestions for %s in %dms", len(data.CurationSuggestions), rawLang, time.Since(tip)/time.Millisecond)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}
	/*
		log.Infof(ctx, "Having/NotHaving sampling...")
		for i := 0; i < sampleSize; i++ {
			// TODO make a single Datastore call for 3 Missing impls
			// TODO make a single Datastore call for 3 Curation suggestions
			// TODO also, why are they so damn slow?
			// TODO last resort, consider fetching those concurrently
			// TODO retrieve the Having data from full text indexes instead?
			log.Infof(ctx, "1 NotHaving...")
			{
				_, idiom, err := dao.randomIdiomNotHaving(ctx, rawLang)
				if err == nil {
					idiom.Implementations = nil
					data.MissingImpl = append(data.MissingImpl, idiom)
				} else {
					log.Errorf(ctx, "getting random idiom not having %s: %v", rawLang, err)
				}
			}
			log.Infof(ctx, "1 Having...")
			{
				_, idiom, err := dao.randomIdiomHaving(ctx, rawLang)
				if err == nil {
					idiom.Implementations = keepSingleImplForLanguage(idiom.Implementations, rawLang)
					data.CurationSuggestions = append(data.CurationSuggestions, idiom)
				} else {
					log.Errorf(ctx, "getting random idiom not having %s: %v", rawLang, err)
				}
			}
		}
	*/
	tip = time.Now()
	data.MissingDocDemo, err = searchMissingDocDemoForLang(ctx, rawLang, sampleSize)
	log.Infof(ctx, "got %d missingDoc and %d missingDemo for %s in %dms", len(data.MissingDocDemo.MissingDoc), len(data.MissingDocDemo.MissingDemo), rawLang, time.Since(tip)/time.Millisecond)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}

	tip = time.Now()
	data.MissingImpl, err = searchMissingImplForLang(ctx, lang, sampleSize)
	log.Infof(ctx, "got %d missingImpl idioms %s in %dms", len(data.MissingImpl.Stubs), rawLang, time.Since(tip)/time.Millisecond)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}

	log.Infof(ctx, "Done computing backlog for %s", lang)
	log.Infof(ctx, "Executing backlog template for %s", lang)
	tip = time.Now()
	err = templates.ExecuteTemplate(w, "page-backlog-language", data)
	log.Infof(ctx, "Done executing backlog template for %s in %dms", lang, time.Since(tip)/time.Millisecond)
	return err
}

// SearchResultsFacade is the facade for the Search Results page.
type BacklogLanguageFacade struct {
	PageMeta            PageMeta
	UserProfile         UserProfile
	Lang                string
	RecommendedDemoSite DemoSite
	CurationSuggestions []*IdiomSingleton
	MissingDocDemo      backlogMissingDocDemo
	MissingImpl         backlogMissingImpl
}

// IdiomSingleton is an idiom that contains a single Implementation.
// It is intended for display logic only.
// It is NOT intended to be saved to the database.
type IdiomSingleton = Idiom

// IdiomStub is an idiom that contains zero Implementation.
// It is intended for display logic only.
// It is NOT intended to be saved to the database.
type IdiomStub = Idiom

type DemoSite struct {
	Name string
	URL  string
}

func recommendedDemoSite(lang string) DemoSite {
	var ds DemoSite
	switch strings.TrimSpace(strings.ToLower(lang)) {
	case "go":
		ds.Name = "the Go Playground"
		ds.URL = "https://play.golang.org/"
	case "csharp", "cs", "fsharp", "fs":
		ds.Name = "SharpLab"
		ds.URL = "https://sharplab.io/"
	default:
		// No recommended demo site for this lang
	}
	return ds
}

func keepForLanguage(impls []Impl, lang string) []Impl {
	var keep []Impl
	for _, impl := range impls {
		// TODO normalize lowercase, etc?
		if impl.LanguageName == lang {
			keep = append(keep, impl)
		}
	}
	return keep
}

func keepSingleImplForLanguage(impls []Impl, lang string) []Impl {
	impls = keepForLanguage(impls, lang)
	if len(impls) == 0 {
		return impls
	}
	impl := impls[rand.Intn(len(impls))]
	return []Impl{
		impl,
	}
}

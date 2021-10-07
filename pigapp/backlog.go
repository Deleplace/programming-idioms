package main

import (
	"math/rand"
	"net/http"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"
	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"
)

// Language Backlog: community contribution nudge.
// Shows a few stats about the Coverage for this lang,
// and links to idioms/impls to be improved.

// 4 sections will show sampleSize lines each: Missing impl, Missing doc, Missing demo, Curation.
const sampleSize = 2

func backlogForLanguage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	rawLang := vars["lang"]
	lang := NormLang(rawLang)
	langs := []string{lang}
	log.Infof(ctx, "Computing backlog for %s", lang)

	data := &BacklogLanguageFacade{
		PageMeta: PageMeta{
			PageTitle: "Idioms missing data for the " + lang + " implementation",
			Toggles:   toggles,
			ExtraCss:  []string{hostPrefix() + themeDirectory() + "/css/pages/backlog.css"},
		},
		UserProfile:         readUserProfile(r),
		Lang:                lang,
		RecommendedDemoSite: recommendedDemoSite(rawLang),
	}

	log.Infof(ctx, "searchRandomImplsForLang(%q)...", rawLang)
	var err error
	data.CurationSuggestions, err = searchRandomImplsForLang(ctx, rawLang, sampleSize)
	log.Infof(ctx, "got %d results", len(data.CurationSuggestions))
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
	data.MissingDoc, data.MissingDemo, err = searchMissingDocDemoForLang(ctx, rawLang, sampleSize)
	log.Infof(ctx, "got %d missingDoc and %d missingDemo for %s", len(data.MissingDoc), len(data.MissingDemo), rawLang)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}

	/*
		log.Infof(ctx, "Obsolete computations...")
		// Warning: manipulating a lot of idioms in memory can get expensive.
		// TODO: get IDs only, then substract IDs of those having DemoURL + DocumentationURL,
		// then get the idioms by IDs.
		maxFetch := 200
		hits, err := dao.searchIdiomsByLangs(ctx, langs, maxFetch)
		if err != nil {
			return err
		}

		// Better if the portion shown varies
		shuffleIdioms(hits)

		maxShow := 25
		results := make([]*IdiomSingleton, 0, maxShow)
		for _, idiom := range hits {
			keep := false
			for _, impl := range implementationsFor(idiom, lang) {
				if isBlank(impl.DemoURL) || isBlank(impl.DocumentationURL) {
					keep = true
					break
				}
			}
			if keep {
				results = append(results, (*IdiomSingleton)(idiom))
				if len(results) >= maxShow {
					break
				}
			}
		}
		//data.CurationSuggestions = results
	*/
	_ = langs

	log.Infof(ctx, "Done computing backlog for %s", lang)
	log.Infof(ctx, "Executing backlog template for %s", lang)
	err = templates.ExecuteTemplate(w, "page-backlog-language", data)
	log.Infof(ctx, "Done executing backlog template for %s", lang)
	return err
}

// SearchResultsFacade is the facade for the Search Results page.
type BacklogLanguageFacade struct {
	PageMeta            PageMeta
	UserProfile         UserProfile
	Lang                string
	RecommendedDemoSite DemoSite
	MissingImpl         []*IdiomStub
	CurationSuggestions []*IdiomSingleton
	MissingDoc          []*IdiomSingleton
	MissingDemo         []*IdiomSingleton
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

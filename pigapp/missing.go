package main

import (
	"math/rand"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"
)

// This screen shows, for a given language, which implementations
// don't have a DemoURL and/or a DocumentationURL.

func missingList(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	lang := vars["lang"]
	lang = NormLang(lang)
	langs := []string{lang}

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
	results := make([]*Idiom, 0, maxShow)
	for _, idiom := range hits {
		keep := false
		for _, impl := range implementationsFor(idiom, lang) {
			if isBlank(impl.DemoURL) || isBlank(impl.DocumentationURL) {
				keep = true
				break
			}
		}
		if keep {
			results = append(results, idiom)
			if len(results) >= maxShow {
				break
			}
		}
	}

	data := &MissingFieldsFacade{
		PageMeta: PageMeta{
			PageTitle: "Idioms missing data for the " + lang + " implementation",
			Toggles:   toggles,
		},
		UserProfile: readUserProfile(r),
		Lang:        lang,
		Results:     results,
	}

	return templates.ExecuteTemplate(w, "page-missing-list", data)
}

// SearchResultsFacade is the facade for the Search Results page.
type MissingFieldsFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Lang        string
	Results     []*Idiom
}

// Returns the list of implementations in this idiom, for this languages.
func implementationsFor(idiom *Idiom, lang string) []*Impl {
	impls := make([]*Impl, 0, 2)
	for i := range idiom.Implementations {
		impl := &idiom.Implementations[i]
		if impl.LanguageName == lang {
			impls = append(impls, impl)
		}
	}
	return impls
}

func shuffleIdioms(idioms []*Idiom) {
	for i := range idioms {
		j := rand.Intn(i + 1)
		idioms[i], idioms[j] = idioms[j], idioms[i]
	}
}

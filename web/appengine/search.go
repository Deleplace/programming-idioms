package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	. "github.com/Deleplace/programming-idioms/idioms"

	"github.com/gorilla/mux"

	"context"

	"google.golang.org/appengine/log"
)

//
// This file is about full text search of idioms and implementations.
//

// SearchResultsFacade is the facade for the Search Results page.
type SearchResultsFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Q           string
	Results     []*Idiom
}

// User types terms, but it's important to recognize when a term is a language name.
func separateLangKeywords(terms []string) (words, langs []string) {
	words = make([]string, 0, len(terms))
	for _, term := range terms {
		lang := NormLang(term)
		if lang == "" {
			words = append(words, term)
		} else {
			langs = append(langs, lang)
		}
	}
	return words, langs
}

// This is a "word by word" search, not a rdbms "like" filter
func search(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	q := vars["q"]
	//q := url.QueryUnescape(q)  Not needed, so it seems.
	hits, normalizedQ, err := findResults(r, q)
	if err != nil {
		if err == errEmptyQ {
			redirURL := hostPrefix() + "/about#about-block-all-idioms"
			http.Redirect(w, r, redirURL, http.StatusFound)
			return nil
		}
		return err
	}

	return listResults(w, r, normalizedQ, hits)
}

var errEmptyQ = fmt.Errorf("Empty search query")

func findResults(r *http.Request, q string) (results []*Idiom, normalizedQ string, err error) {
	ctx := r.Context()

	// Maybe someday we find a graceful way to handle "c++", "c#", etc. but...
	q = strings.Replace(q, "C++", "Cpp", -1)
	q = strings.Replace(q, "c++", "cpp", -1)
	q = strings.Replace(q, "C#", "Csharp", -1)
	q = strings.Replace(q, "c#", "csharp", -1)
	q = strings.Replace(q, "C♯", "Csharp", -1)
	q = strings.Replace(q, "c♯", "csharp", -1)

	terms := SplitForSearching(q, true)

	words, typedLangs := separateLangKeywords(terms)

	words = FilterStrings(words, func(chunk string) bool {
		// Small terms (1 or 2 chars) must be discarded,
		// because they weren't indexed in the first place.
		// (Words only, not language names.)
		return len(chunk) >= 3 || RegexpDigitsOnly.MatchString(chunk)
	})

	if len(words)+len(typedLangs) == 0 {
		// Search query is empty or illegible...
		return nil, "", errEmptyQ
	}

	if len(words) == 0 {
		words, typedLangs = typedLangs, nil
	}

	typedLangsSet := make(map[string]bool, len(typedLangs))
	for _, lang := range typedLangs {
		typedLangsSet[lang] = true
	}

	matchingPromise := matchingImplPromise(ctx, words, typedLangs)

	numberMaxResults := 20
	// Note that this currently depends on userProfile.FavoriteLanguages
	// (not the best for caching and for SAP)
	userProfile := readUserProfile(r)
	hits, err := dao.searchIdiomsByWordsWithFavorites(ctx, words, typedLangs, userProfile.FavoriteLanguages, userProfile.SeeNonFavorite, numberMaxResults)
	if err != nil {
		return nil, "", err
	}

	matchingImplIDs := <-matchingPromise

	for _, idiom := range hits {
		implFavoriteLanguagesFirstWithOrder(idiom, userProfile.FavoriteLanguages, "", userProfile.SeeNonFavorite)
		for i := range idiom.Implementations {
			impl := &idiom.Implementations[i]
			implIDStr := fmt.Sprintf("%d", impl.Id)
			if matchingImplIDs[implIDStr] {
				impl.Deco.Matching = true
			}
			if typedLangsSet[impl.LanguageName] {
				impl.Deco.SearchedLang = true
			}
		}
	}
	return hits, strings.Join(terms, " "), nil
}

func matchingImplPromise(ctx context.Context, words, typedLangs []string) chan map[string]bool {
	ch := make(chan map[string]bool)
	go func() {
		// Highlight matching impls :)
		matchingImplIDs, err := dao.searchImplIDs(ctx, words, typedLangs)
		if err == nil {
			ch <- matchingImplIDs
		} else {
			log.Errorf(ctx, "problem fetching impl highlights: %v", err)
			ch <- map[string]bool{}
		}
		close(ch)
	}()
	return ch
}

func listResults(w http.ResponseWriter, r *http.Request, q string, idioms []*Idiom) error {
	data := &SearchResultsFacade{
		PageMeta: PageMeta{
			PageTitle:   "Idioms for \"" + q + "\"",
			Toggles:     toggles,
			SearchQuery: q,
		},
		UserProfile: readUserProfile(r),
		Q:           q,
		Results:     idioms,
	}

	return templates.ExecuteTemplate(w, "page-list-results", data)
	// return templates.ExecuteTemplate(w, "page-list-results-minimal", data)
}

func searchRedirect(w http.ResponseWriter, r *http.Request) error {
	q := r.FormValue("q")
	if q == "" {
		// Small hack for own convenience: empty search -> homepage
		http.Redirect(w, r, "/", 301)
		return nil
	}

	// #102 Slash in query breaks the search URL.
	// Mitigation: remove them.
	q = strings.Replace(q, "/", " ", -1)

	safeQ := url.QueryEscape(q)

	http.Redirect(w, r, hostPrefix()+"/search/"+safeQ, http.StatusMovedPermanently)
	return nil
}

func listByLanguage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	userProfile := readUserProfile(r)

	ctx := r.Context()
	langsStr := vars["langs"]
	langs := strings.Split(langsStr, "_")
	langs = MapStrings(langs, NormLang)
	langs = RemoveEmptyStrings(langs)

	numberMaxResults := 20
	hits, err := dao.searchIdiomsByLangs(ctx, langs, numberMaxResults)
	if err != nil {
		return err
	}

	for _, idiom := range hits {
		implFavoriteLanguagesFirstWithOrder(idiom, userProfile.FavoriteLanguages, "", userProfile.SeeNonFavorite)
	}

	niceLangs := MapStrings(langs, PrintNiceLang)
	return listResults(w, r, fmt.Sprintf("Language=%v", niceLangs), hits)
}

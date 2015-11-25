package pigae

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"appengine"
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
		if normLang(term) == "" {
			words = append(words, term)
		} else {
			langs = append(langs, term)
		}
	}
	return words, langs
}

// This is a "word by word" search, not a rdbms "like" filter
func search(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	userProfile := readUserProfile(r)

	c := appengine.NewContext(r)
	q := vars["q"]
	terms := SplitForSearching(q, true)

	words, typedLangs := separateLangKeywords(terms)
	if len(words) == 0 {
		words, typedLangs = typedLangs, nil
	}

	numberMaxResults := 20
	hits, err := dao.searchIdiomsByWordsWithFavorites(c, words, typedLangs, userProfile.FavoriteLanguages, userProfile.SeeNonFavorite, numberMaxResults)
	if err != nil {
		return err
	}

	// Highlight matching impls :)
	matchingImplIDs, err := dao.searchImplIDs(c, words)
	if err != nil {
		return err
	}
	for _, idiom := range hits {
		implFavoriteLanguagesFirstWithOrder(idiom, userProfile.FavoriteLanguages, "", userProfile.SeeNonFavorite)
		for i := range idiom.Implementations {
			impl := &idiom.Implementations[i]
			implIDStr := fmt.Sprintf("%d", impl.Id)
			if matchingImplIDs[implIDStr] {
				impl.Deco.Matching = true
			}
		}
	}

	normalizedQ := strings.Join(words, " ")
	return listResults(w, r, normalizedQ, hits)
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
}

func searchRedirect(w http.ResponseWriter, r *http.Request) error {
	q := r.FormValue("q")
	if q == "" {
		// Small hack for own convenience: empty search -> homepage
		http.Redirect(w, r, "/", 301)
		return nil
	}
	q = strings.Replace(q, " ", "+", -1)

	http.Redirect(w, r, hostPrefix()+"/search/"+q, http.StatusMovedPermanently)
	return nil
}

func listByLanguage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	userProfile := readUserProfile(r)

	c := appengine.NewContext(r)
	langsStr := vars["langs"]
	langs := strings.Split(langsStr, "_")
	langs = MapStrings(langs, normLang)
	langs = RemoveEmptyStrings(langs)

	numberMaxResults := 20
	hits, err := dao.searchIdiomsByLangs(c, langs, numberMaxResults)
	if err != nil {
		return err
	}

	for _, idiom := range hits {
		implFavoriteLanguagesFirstWithOrder(idiom, userProfile.FavoriteLanguages, "", userProfile.SeeNonFavorite)
	}

	niceLangs := MapStrings(langs, printNiceLang)
	return listResults(w, r, fmt.Sprintf("Language=%v", niceLangs), hits)
}

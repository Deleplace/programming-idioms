package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"appengine"
)

// This screen shows, for a given language, which implementations
// don't have a DemoURL and/or a DocumentationURL.

func missingList(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	lang := vars["lang"]
	lang = normLang(lang)
	langs := []string{lang}

	numberMaxResults := 20
	hits, err := dao.searchIdiomsByLangs(c, langs, numberMaxResults)
	if err != nil {
		return err
	}

	data := &MissingFieldsFacade{
		PageMeta: PageMeta{
			PageTitle: "Idioms missing data for the " + lang + " implementation",
			Toggles:   toggles,
		},
		UserProfile: readUserProfile(r),
		Results:     hits,
	}

	return templates.ExecuteTemplate(w, "page-missing-list", data)
}

// SearchResultsFacade is the facade for the Search Results page.
type MissingFieldsFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
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

package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"
	"github.com/gorilla/mux"
)

//
// A Cheat Sheet is a single page containing the implementations of all idioms, for a
// single language.
//

// AllIdiomsFacade is the Facade for the Cheat Sheets.
type CheatSheetFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Lang        string
	Idioms      []*Idiom
}

func cheatsheet(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	lang := vars["lang"]

	idioms, err := retrieveAllIdioms(r) // TODO filter this very much
	if err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}

	filteredIdioms := make([]*Idiom, 0, len(idioms))
	// This loop won't be necessary, when we have a Datastore query
	// which fetches only the relevant implementations.
	for _, idiom := range idioms {
		for _, impl := range idiom.Implementations {
			if impl.LanguageName == lang {
				newIdiom := *idiom
				newIdiom.Implementations = []Impl{impl}
				filteredIdioms = append(filteredIdioms, &newIdiom)
			}
		}
	}

	data := CheatSheetFacade{
		PageMeta: PageMeta{
			PageTitle: printNiceLang(lang) + " cheat sheet",
			Toggles:   toggles,
		},
		UserProfile: readUserProfile(r),
		Lang:        lang,
		Idioms:      filteredIdioms,
	}

	if err := templates.ExecuteTemplate(w, "page-cheatsheet", data); err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}
	return nil
}

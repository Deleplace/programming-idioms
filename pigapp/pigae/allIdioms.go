package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"
	"golang.org/x/net/context"
)

// AllIdiomsFacade is the Facade for block All Idioms of the About page.
type AllIdiomsFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	AllIdioms   []*Idiom
}

func allIdioms(c context.Context, w http.ResponseWriter, r *http.Request) error {

	idioms, err := retrieveAllIdioms(c, r)
	if err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}

	data := AllIdiomsFacade{
		PageMeta: PageMeta{
			PageTitle: "All idioms",
			Toggles:   toggles,
		},
		UserProfile: readUserProfile(c, r),
		AllIdioms:   idioms,
	}

	if err := templates.ExecuteTemplate(w, "page-all-idioms", data); err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}
	return nil
}

func retrieveAllIdioms(c context.Context, r *http.Request) ([]*Idiom, error) {
	// TODO sort by popularity desc
	// TODO limit to 50, + button [See more...]  or pagination

	_, idioms, err := dao.getAllIdioms(c, 0, "Id")
	if err != nil {
		return nil, err
	}

	favlangs := lookForFavoriteLanguages(r)
	includeNonFav := seeNonFavorite(r)
	for _, idiom := range idioms {
		implFavoriteLanguagesFirstWithOrder(idiom, favlangs, "", includeNonFav)
	}

	return idioms, nil
}

package main

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"
)

// AllIdiomsFacade is the Facade for block All Idioms of the About page.
type AllIdiomsFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	AllIdioms   []*Idiom
}

func allIdioms(w http.ResponseWriter, r *http.Request) error {

	idioms, err := retrieveAllIdioms(r, true)
	if err != nil {
		return PiErrorf(http.StatusInternalServerError, "%v", err)
	}

	data := AllIdiomsFacade{
		PageMeta: PageMeta{
			PageTitle: "All idioms",
			Toggles:   toggles,
		},
		UserProfile: readUserProfile(r),
		AllIdioms:   idioms,
	}

	if err := templates.ExecuteTemplate(w, "page-all-idioms", data); err != nil {
		return PiErrorf(http.StatusInternalServerError, "%v", err)
	}
	return nil
}

func retrieveAllIdioms(r *http.Request, orderByFav bool) ([]*Idiom, error) {
	ctx := r.Context()
	// TODO sort by popularity desc
	// TODO limit to 50, + button [See more...]  or pagination

	_, idioms, err := dao.getAllIdioms(ctx, 0, "Id")
	if err != nil {
		return nil, err
	}

	if orderByFav {
		favlangs := lookForFavoriteLanguages(r)
		includeNonFav := seeNonFavorite(r)
		for _, idiom := range idioms {
			implFavoriteLanguagesFirstWithOrder(idiom, favlangs, "", includeNonFav)
		}
	}

	return idioms, nil
}

package main

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/idioms"
)

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

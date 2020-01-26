package main

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"context"
)

// HomeFacade is the Facade for the homepage.
type HomeFacade struct {
	PageMeta          PageMeta
	UserProfile       UserProfile
	LastUpdatedIdioms []*Idiom
	PopularIdioms     []*Idiom
}

func home(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userProfile := readUserProfile(r)
	return homeView(w, ctx, userProfile)
}

// Possible controllers include : home(), bookmarkableUserURL()
func homeView(w http.ResponseWriter, ctx context.Context, userProfile UserProfile) error {

	homeToggles := copyToggles(toggles)

	data := &HomeFacade{
		PageMeta: PageMeta{
			PageTitle: "Programming Idioms",
			Toggles:   homeToggles,
		},
		UserProfile:       userProfile,
		LastUpdatedIdioms: nil,
		PopularIdioms:     nil,
	}

	var err error
	if homeToggles["homeBlockLastUpdated"] {
		data.LastUpdatedIdioms, err = dao.recentIdioms(ctx, userProfile.FavoriteLanguages, userProfile.SeeNonFavorite, 5)
		if err != nil {
			return err
		}
	}

	if homeToggles["homeBlockPopular"] {
		data.PopularIdioms, err = dao.popularIdioms(ctx, userProfile.FavoriteLanguages, userProfile.SeeNonFavorite, 3)
		if err != nil {
			return err
		}
	}

	// return templates.ExecuteTemplate(w, "page-home", data)
	return templates.ExecuteTemplate(w, "page-home-minimal", data)
}

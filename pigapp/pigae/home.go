package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
)

// HomeFacade is the Facade for the homepage.
type HomeFacade struct {
	PageMeta          PageMeta
	UserProfile       UserProfile
	LastUpdatedIdioms []*Idiom
	PopularIdioms     []*Idiom
}

func home(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	userProfile := readUserProfile(r)
	return homeView(w, c, userProfile)
}

// Possible controllers include : home(), bookmarkableUserURL()
func homeView(w http.ResponseWriter, c appengine.Context, userProfile UserProfile) error {
	idiomsRecent, err := dao.recentIdioms(c, userProfile.FavoriteLanguages, userProfile.SeeNonFavorite, 5)
	if err != nil {
		return err
	}

	idiomsPopular, err := dao.popularIdioms(c, userProfile.FavoriteLanguages, userProfile.SeeNonFavorite, 3)
	if err != nil {
		return err
	}

	homeToggles := copyToggles(toggles)

	data := &HomeFacade{
		PageMeta: PageMeta{
			PageTitle: "Programming Idioms",
			Toggles:   homeToggles,
		},
		UserProfile:       userProfile,
		LastUpdatedIdioms: idiomsRecent,
		PopularIdioms:     idiomsPopular,
	}

	return templates.ExecuteTemplate(w, "page-home", data)
}

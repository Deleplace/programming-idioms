package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"appengine"
)

// IdiomHistoryFacade is the Facade for the Idiom History page.
type IdiomHistoryFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
}

func idiomHistory(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	c := appengine.NewContext(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	// TODO
	_, _ = c, idiomID

	userProfile := readUserProfile(r)
	myToggles := copyToggles(toggles)
	myToggles["actionEditIdiom"] = false
	myToggles["actionAddImpl"] = false
	data := &IdiomHistoryFacade{
		PageMeta: PageMeta{
			PageTitle: "Idiom " + idiomIDStr + " history",
			Toggles:   myToggles,
		},
		UserProfile: userProfile,
	}
	return templates.ExecuteTemplate(w, "page-history", data)
}

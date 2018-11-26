package main

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
)

// IdiomEditFacade is the Facade for the Add Idiom Picture page.
type IdiomEditFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Idiom       *Idiom
}

func idiomEdit(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	c := appengine.NewContext(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Idiom " + idiomIDStr + " not found : " + err.Error(), http.StatusNotFound}
	}

	userProfile := readUserProfile(r)
	myToggles := copyToggles(toggles)
	myToggles["editing"] = true

	data := &IdiomEditFacade{
		PageMeta: PageMeta{
			PageTitle:             fmt.Sprintf("Editing Idiom %d : %s", idiom.Id, idiom.Title),
			Toggles:               myToggles,
			PreventIndexingRobots: true,
		},
		UserProfile: userProfile,
		Idiom:       idiom,
	}

	return templates.ExecuteTemplate(w, "page-idiom-edit", data)
}

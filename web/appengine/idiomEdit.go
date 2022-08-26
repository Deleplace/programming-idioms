package main

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/idioms"

	"github.com/gorilla/mux"
)

// IdiomEditFacade is the Facade for the Add Idiom Picture page.
type IdiomEditFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Idiom       *Idiom
}

func idiomEdit(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ctx := r.Context()

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Idiom %q not found : %v", idiomIDStr, err)
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

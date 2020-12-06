package main

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"
)

//
// 2016-09 strategy: no file upload.
// Admin only.
// Admin types URL, preferably to dedicated GCS bucket programming-idioms-pictures.
//

// IdiomAddPictureFacade is the Facade for the Add Idiom Picture page.
type IdiomAddPictureFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Idiom       *Idiom
}

func idiomAddPicture(w http.ResponseWriter, r *http.Request) error {
	if !IsAdmin(r) {
		return fmt.Errorf("For now, only the Admin may add an idiom picture.")
	}

	vars := mux.Vars(r)
	ctx := r.Context()

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)
	if idiomID == -1 {
		return PiErrorf(http.StatusBadRequest, "%q is not a valid idiom id.", idiomIDStr)
	}

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", idiomIDStr)
	}

	myToggles := copyToggles(toggles)
	myToggles["editing"] = true

	data := &IdiomAddPictureFacade{
		PageMeta: PageMeta{
			PageTitle: fmt.Sprintf("Adding picture to Idiom %d : %s", idiom.Id, idiom.Title),
			Toggles:   myToggles,
		},
		UserProfile: readUserProfile(r),
		Idiom:       idiom,
	}

	return templates.ExecuteTemplate(w, "page-idiom-add-picture", data)
}

func idiomSavePicture(w http.ResponseWriter, r *http.Request) error {
	if !IsAdmin(r) {
		return fmt.Errorf("For now, only the Admin may add an idiom picture.")
	}

	ctx := r.Context()
	userProfile := readUserProfile(r)

	idiomIDStr := r.FormValue("idiom_id")
	pictureURL := r.FormValue("picture_url")

	idiomID := String2Int(idiomIDStr)
	if idiomID == -1 {
		return PiErrorf(http.StatusBadRequest, "%q is not a valid idiom id.", idiomIDStr)
	}

	key, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", idiomIDStr)
	}

	idiom.ImageURL = pictureURL
	idiom.EditSummary = "Updated picture URL by user [" + userProfile.Nickname + "]"
	idiom.LastEditor = userProfile.Nickname

	err = dao.saveExistingIdiom(ctx, key, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceIdiomURL(idiom), http.StatusFound)
	return nil
}

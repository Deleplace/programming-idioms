package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
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
	c := appengine.NewContext(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)
	if idiomID == -1 {
		return PiError{idiomIDStr + " is not a valid idiom id.", http.StatusBadRequest}
	}

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
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

	c := appengine.NewContext(r)
	userProfile := readUserProfile(r)

	idiomIDStr := r.FormValue("idiom_id")
	pictureURL := r.FormValue("picture_url")

	idiomID := String2Int(idiomIDStr)
	if idiomID == -1 {
		return PiError{idiomIDStr + " is not a valid idiom id.", http.StatusBadRequest}
	}

	key, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
	}

	idiom.ImageURL = pictureURL
	idiom.EditSummary = "Updated picture URL by user [" + userProfile.Nickname + "]"

	err = dao.saveExistingIdiom(c, key, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceIdiomURL(idiom), http.StatusFound)
	return nil
}

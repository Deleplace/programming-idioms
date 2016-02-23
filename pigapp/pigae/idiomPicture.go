package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
)

// IdiomAddPictureFacade is the Facade for the Add Idiom Picture page.
type IdiomAddPictureFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Idiom       *Idiom
	UploadURL   string
}

func idiomAddPicture(w http.ResponseWriter, r *http.Request) error {
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

	uploadURL, err := blobstore.UploadURL(c, "/picture-upload", nil) /*,
	&blobstore.UploadURLOptions{
		MaxUploadBytes:        105 * 1024,
		MaxUploadBytesPerBlob: 105 * 1024,
		StorageBucket:         "programming-idioms-pictures/idiom/" + idiomIDStr + "/lead/",
	})*/
	if err != nil {
		return err
	}

	data := &IdiomAddPictureFacade{
		PageMeta: PageMeta{
			PageTitle: fmt.Sprintf("Adding picture to Idiom %d : %s", idiom.Id, idiom.Title),
			Toggles:   myToggles,
		},
		UserProfile: readUserProfile(r),
		Idiom:       idiom,
		UploadURL:   uploadURL.String(),
	}

	return templates.ExecuteTemplate(w, "page-idiom-add-picture", data)
}

func idiomSavePicture(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	newBlobKey, otherParams, err := dao.processUploadFile(r, "idiom_picture")
	if err != nil {
		return err
	}
	// TODO check picture weight, format and dimensions

	idiomIDStr := otherParams["idiom_id"][0]
	idiomID := String2Int(idiomIDStr)
	if idiomID == -1 {
		return PiError{idiomIDStr + " is not a valid idiom id.", http.StatusBadRequest}
	}

	key, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
	}

	idiom.Picture = newBlobKey

	err = dao.saveExistingIdiom(c, key, idiom)
	if err != nil {
		return err
	}

	http.Redirect(w, r, NiceIdiomURL(idiom), http.StatusFound)
	return nil
}

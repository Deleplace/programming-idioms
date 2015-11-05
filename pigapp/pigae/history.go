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
	IdiomID     int
	// HistoryList contains incomplete IdiomHistory objects: just a few fields.
	HistoryList []*IdiomHistory
}

func idiomHistory(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	c := appengine.NewContext(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	_, list, err := dao.getIdiomHistoryList(c, idiomID)
	if err != nil {
		return err
	}

	revertedVersionStr := r.FormValue("reverted")
	revertedVersion := String2Int(revertedVersionStr)
	if revertedVersion == list[0].Version {
		// Most often, the deletion has not been taken into
		// account by the datastore yet, so we remove manually.
		list = list[1:]
	}

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
		IdiomID:     idiomID,
		HistoryList: list,
	}
	return templates.ExecuteTemplate(w, "page-history", data)
}

func revertIdiomVersion(w http.ResponseWriter, r *http.Request) error {
	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)
	versionStr := r.FormValue("version")
	version := String2Int(versionStr)
	c := appengine.NewContext(r)

	_, err := dao.revert(c, idiomID, version)
	if err != nil {
		return err
	}
	redirUrl := hostPrefix() + "/history/" + idiomIDStr + "?reverted=" + versionStr
	http.Redirect(w, r, redirUrl, http.StatusFound)
	return nil
	// Unfortunately, the redirect page doesn't see the history deletion, yet.
}

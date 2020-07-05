package main

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"
)

// IdiomHistoryFacade is the Facade for the Idiom History page.
type IdiomHistoryFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Idiom       *Idiom
	// HistoryList contains incomplete IdiomHistory objects: just a few fields.
	HistoryList []*IdiomHistory
}

func idiomHistory(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	ctx := r.Context()

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	_, list, err := dao.getIdiomHistoryList(ctx, idiomID)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return PiError{ErrorText: "No history entries found for idiom " + idiomIDStr, Code: 500}
	}

	revertedVersionStr := r.FormValue("reverted")
	revertedVersion := String2Int(revertedVersionStr)
	if revertedVersion == list[0].Version {
		// Most often, the deletion has not been taken into
		// account by the datastore yet, so we remove manually.
		list = list[1:]
	}

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return err
	}

	userProfile := readUserProfile(r)
	myToggles := copyToggles(toggles)
	myToggles["actionEditIdiom"] = false
	myToggles["actionIdiomHistory"] = false
	myToggles["actionAddImpl"] = false
	data := &IdiomHistoryFacade{
		PageMeta: PageMeta{
			PageTitle:             "Idiom " + idiomIDStr + " history",
			Toggles:               myToggles,
			PreventIndexingRobots: true,
		},
		UserProfile: userProfile,
		Idiom:       idiom,
		HistoryList: list,
	}
	return templates.ExecuteTemplate(w, "page-history", data)
}

func revertIdiomVersion(w http.ResponseWriter, r *http.Request) error {
	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)
	versionStr := r.FormValue("version")
	version := String2Int(versionStr)
	ctx := r.Context()

	_, err := dao.revert(ctx, idiomID, version)
	if err != nil {
		return err
	}
	redirUrl := hostPrefix() + "/history/" + idiomIDStr + "?reverted=" + versionStr
	http.Redirect(w, r, redirUrl, http.StatusFound)
	return nil
	// Unfortunately, the redirect page doesn't see the history deletion, yet.
}

func restoreIdiomVersion(w http.ResponseWriter, r *http.Request) error {
	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)
	versionStr := r.FormValue("version")
	version := String2Int(versionStr)
	ctx := r.Context()
	restoreUser := lookForNickname(r)

	idiom, err := dao.historyRestore(ctx, idiomID, version, restoreUser)
	if err != nil {
		return err
	}
	redirUrl := NiceIdiomURL(idiom)
	http.Redirect(w, r, redirUrl, http.StatusFound)
	return nil
	// Unfortunately, the redirect page doesn't see the history deletion, yet.
}

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
	return templates.ExecuteTemplate(w, "page-idiom-history", data)
}

// ImplHistoryFacade is the Facade for the Impl History page.
type ImplHistoryFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Idiom       *Idiom
	HistoryList []*IdiomHistory
	// ImplID because we're interested only in a specific implementation
	ImplID int
}

func implHistory(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	ctx := r.Context()

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)
	implIDStr := vars["implId"]
	implID := String2Int(implIDStr)

	_, list, err := dao.getDenseHistoryList(ctx, idiomID)
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

	sublist := make([]*IdiomHistory, 0, len(list))
	for i := 0; i < len(list)-1; i++ {
		h1, h2 := list[i], list[i+1]
		_, impl1, ok1 := h1.FindImplInIdiom(implID)
		_, impl2, ok2 := h2.FindImplInIdiom(implID)
		if (ok1 != ok2) || (ok1 && ok2 && impl1.Version != impl2.Version) {
			sublist = append(sublist, h1)
		}
	}
	if len(list) > 0 {
		h := list[len(list)-1]
		if _, _, ok := h.FindImplInIdiom(implID); ok {
			sublist = append(sublist, h)
		}
	}
	list = sublist

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return err
	}

	userProfile := readUserProfile(r)
	myToggles := copyToggles(toggles)
	myToggles["actionEditIdiom"] = false
	myToggles["actionIdiomHistory"] = false
	myToggles["actionAddImpl"] = false
	data := &ImplHistoryFacade{
		PageMeta: PageMeta{
			PageTitle:             "Idiom " + idiomIDStr + " history",
			Toggles:               myToggles,
			PreventIndexingRobots: true,
		},
		UserProfile: userProfile,
		Idiom:       idiom,
		HistoryList: list,
		ImplID:      implID,
	}
	return templates.ExecuteTemplate(w, "page-impl-history", data)
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
	why := r.FormValue("why")
	ctx := r.Context()
	restoreUser := lookForNickname(r)

	idiom, err := dao.historyRestore(ctx, idiomID, version, restoreUser, why)
	if err != nil {
		return err
	}
	redirUrl := NiceIdiomURL(idiom)
	http.Redirect(w, r, redirUrl, http.StatusFound)
	return nil
	// Unfortunately, the redirect page doesn't see the history deletion, yet.
}

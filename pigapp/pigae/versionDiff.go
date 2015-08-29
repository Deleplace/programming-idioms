package pigae

import (
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"appengine"
)

// VersionDiffFacade is the Facade for the Diff page.
type VersionDiffFacade struct {
	PageMeta              PageMeta
	UserProfile           UserProfile
	IdiomLeft, IdiomRight *IdiomHistory
}

func versionDiff(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	c := appengine.NewContext(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)
	v1Str := vars["v1"]
	v1 := String2Int(v1Str)
	v2Str := vars["v2"]
	v2 := String2Int(v2Str)

	_, left, err := dao.getIdiomHistory(c, idiomID, v1)
	if err != nil {
		return PiError{err.Error(), http.StatusNotFound}
	}
	_, right, err := dao.getIdiomHistory(c, idiomID, v2)
	if err != nil {
		return PiError{err.Error(), http.StatusNotFound}
	}

	userProfile := readUserProfile(r)
	myToggles := copyToggles(toggles)
	myToggles["actionEditIdiom"] = false
	myToggles["actionAddImpl"] = false
	data := &VersionDiffFacade{
		PageMeta: PageMeta{
			PageTitle: right.Title,
			Toggles:   myToggles,
		},
		UserProfile: userProfile,
		IdiomLeft:   left,
		IdiomRight:  right,
	}
	return templates.ExecuteTemplate(w, "page-idiom-version-diff", data)
}

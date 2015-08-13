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

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
	}
	left := &IdiomHistory{*idiom}
	right := &IdiomHistory{*idiom}

	// TODO fetch real versions...
	left.Title = "The original Title"
	right.Title = "The very modified Title"

	userProfile := readUserProfile(r)
	myToggles := copyToggles(toggles)
	//	myToggles["actionEditIdiom"] = true
	//	myToggles["actionAddImpl"] = true
	data := &VersionDiffFacade{
		PageMeta: PageMeta{
			PageTitle: idiom.Title,
			Toggles:   myToggles,
		},
		UserProfile: userProfile,
		IdiomLeft:   left,
		IdiomRight:  right,
	}
	return templates.ExecuteTemplate(w, "page-idiom-version-diff", data)
}

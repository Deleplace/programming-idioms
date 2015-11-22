package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"appengine"
)

// ImplEditFacade is the Facade for the Implementation Edit page
type ImplEditFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Idiom       *Idiom
	Impl        *Impl
}

func implEdit(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	c := appengine.NewContext(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	implIDStr := vars["implId"]
	implID := String2Int(implIDStr)

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
	}

	_, impl, exists := idiom.FindImplInIdiom(implID)
	if !exists {
		return PiError{"Could not find implementation " + implIDStr + " in idiom " + idiomIDStr, http.StatusNotFound}
	}

	myToggles := copyToggles(toggles)
	myToggles["editing"] = true

	data := &ImplEditFacade{
		PageMeta: PageMeta{
			PageTitle: fmt.Sprintf("Editing Idiom %d : %s", idiom.Id, idiom.Title),
			Toggles:   myToggles,
		},
		UserProfile: readUserProfile(r),
		Idiom:       idiom,
		Impl:        impl,
	}

	return templates.ExecuteTemplate(w, "page-impl-edit", data)
}

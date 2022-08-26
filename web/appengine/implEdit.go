package main

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/idioms"

	"github.com/gorilla/mux"
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

	ctx := r.Context()

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	implIDStr := vars["implId"]
	implID := String2Int(implIDStr)

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", idiomIDStr)
	}

	_, impl, exists := idiom.FindImplInIdiom(implID)
	if !exists {
		return PiErrorf(http.StatusNotFound, "Could not find implementation %q in idiom %q", implIDStr, idiomIDStr)
	}
	implCopy := *impl

	myToggles := copyToggles(toggles)
	myToggles["editing"] = true

	// Alter the idiom content, in the Facade only, to skip current impl in the
	// [Other implementations] block.
	// Warning: the facade idiom is now 1 impl shorter than reality.
	excludeImpl(idiom, implID)

	data := &ImplEditFacade{
		PageMeta: PageMeta{
			PageTitle:             fmt.Sprintf("Editing Idiom %d : %s", idiom.Id, idiom.Title),
			Toggles:               myToggles,
			PreventIndexingRobots: true,
			ExtraCss: []string{
				hostPrefix() + themeDirectory() + "/css/edit.css",
			},
		},
		UserProfile: readUserProfile(r),
		Idiom:       idiom,
		Impl:        &implCopy,
	}

	return templates.ExecuteTemplate(w, "page-impl-edit", data)
}

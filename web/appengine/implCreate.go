package main

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/idioms"

	"github.com/gorilla/mux"
)

// ImplCreateFacade is the Facade for the New Implementation page.
type ImplCreateFacade struct {
	PageMeta               PageMeta
	UserProfile            UserProfile
	Idiom                  *Idiom
	LanguageSingleSelector LanguageSingleSelector
}

func implCreate(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	ctx := r.Context()
	userProfile := readUserProfile(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	preSelectedLanguage := NormLang(vars["lang"])

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", idiomIDStr)
	}

	// This alters the idiom content in the Facade only
	implFavoriteLanguagesFirstWithOrder(idiom, userProfile.FavoriteLanguages, "", userProfile.SeeNonFavorite)

	myToggles := copyToggles(toggles)
	myToggles["editing"] = true

	data := &ImplCreateFacade{
		PageMeta: PageMeta{
			PageTitle:             fmt.Sprintf("Creating implementation for idiom %d : %s", idiom.Id, idiom.Title),
			Toggles:               myToggles,
			PreventIndexingRobots: true,
			ExtraCss: []string{
				hostPrefix() + themeDirectory() + "/css/edit.css",
			},
		},
		UserProfile: readUserProfile(r),
		Idiom:       idiom,
		LanguageSingleSelector: LanguageSingleSelector{
			FieldName: "impl_language",
			Selected:  preSelectedLanguage,
		},
	}

	return templates.ExecuteTemplate(w, "page-impl-create", data)
}

func excludeImpl(idiom *Idiom, excludedImplID int) {
	// This alters the idiom content in the Facade only
	for i, impl := range idiom.Implementations {
		if impl.Id == excludedImplID {
			copy(idiom.Implementations[i:], idiom.Implementations[i+1:])
			idiom.Implementations = idiom.Implementations[:len(idiom.Implementations)-1]
			break
		}
	}
}

// The block [Other implementations], read-only, at bottom of the page.
//
// 2015-12-23  ajax fetch deactivated because doesn't play well with escaping
// of bubbles text.
func ajaxOtherImplementations(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)

	excludedImplIDStr := r.FormValue("excludedImplId")
	excludedImplID := String2Int(excludedImplIDStr)

	// w.Write([]byte("123 456 789"))
	// return nil

	_, idiom, err := dao.getIdiom(ctx, idiomID)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "Could not find idiom %q", idiomIDStr)
	}

	myToggles := copyToggles(toggles)
	myToggles["editing"] = true
	userProfile := readUserProfile(r)

	if excludedImplID != -1 {
		excludeImpl(idiom, excludedImplID)
	}

	// This alters the idiom content in the Facade only
	implFavoriteLanguagesFirstWithOrder(idiom, userProfile.FavoriteLanguages, "", userProfile.SeeNonFavorite)

	type OtherImplFacade struct {
		PageMeta    PageMeta
		UserProfile UserProfile
		Idiom       *Idiom
		// ExcludedImplId int   not needed
	}

	data := &OtherImplFacade{
		PageMeta: PageMeta{
			Toggles: myToggles,
		},
		UserProfile: userProfile,
		Idiom:       idiom,
		// ExcludedImplId: excludedImplId,   not needed
	}

	return templates.ExecuteTemplate(w, "block-other-implementations", data)
}

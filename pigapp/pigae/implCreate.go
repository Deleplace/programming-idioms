package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"
	"golang.org/x/net/context"

	"github.com/gorilla/mux"
)

// ImplCreateFacade is the Facade for the New Implementation page.
type ImplCreateFacade struct {
	PageMeta               PageMeta
	UserProfile            UserProfile
	Idiom                  *Idiom
	LanguageSingleSelector LanguageSingleSelector
}

func implCreate(c context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	userProfile := readUserProfile(c, r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	preSelectedLanguage := NormLang(vars["lang"])

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
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
		},
		UserProfile: readUserProfile(c, r),
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
func ajaxOtherImplementations(c context.Context, w http.ResponseWriter, r *http.Request) error {

	idiomIDStr := r.FormValue("idiomId")
	idiomID := String2Int(idiomIDStr)

	excludedImplIDStr := r.FormValue("excludedImplId")
	excludedImplID := String2Int(excludedImplIDStr)

	// w.Write([]byte("123 456 789"))
	// return nil

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
	}

	myToggles := copyToggles(toggles)
	myToggles["editing"] = true
	userProfile := readUserProfile(c, r)

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

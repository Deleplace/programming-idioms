package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"appengine"
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

	c := appengine.NewContext(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	preSelectedLanguage := normLang(vars["lang"])

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
	}

	myToggles := copyToggles(toggles)
	myToggles["editing"] = true

	data := &ImplCreateFacade{
		PageMeta: PageMeta{
			PageTitle: fmt.Sprintf("Creating implementation for idiom %d : %s", idiom.Id, idiom.Title),
			Toggles:   myToggles,
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

// The block [Other implementations], read-only, at bottom of the page
func ajaxOtherImplementations(w http.ResponseWriter, r *http.Request) error {

	c := appengine.NewContext(r)

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
	userProfile := readUserProfile(r)

	if excludedImplID != -1 {
		// This alters the idiom content in the Facade only
		for i, impl := range idiom.Implementations {
			if impl.Id == excludedImplID {
				copy(idiom.Implementations[i:], idiom.Implementations[i+1:])
				idiom.Implementations = idiom.Implementations[:len(idiom.Implementations)-1]
				break
			}
		}
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

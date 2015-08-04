package pigae

import (
	"net/http"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"appengine"
)

// IdiomDetailFacade is the Facade for the Idiom Detail page.
type IdiomDetailFacade struct {
	PageMeta       PageMeta
	UserProfile    UserProfile
	Idiom          *Idiom
	SelectedImplID int
}

func idiomDetail(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	c := appengine.NewContext(r)

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)

	_, idiom, err := dao.getIdiom(c, idiomID)
	if err != nil {
		return PiError{"Could not find idiom " + idiomIDStr, http.StatusNotFound}
	}

	idiomTitleInURL := vars["idiomTitle"]
	if idiomTitleInURL != "" && uriNormalize(idiom.Title) != idiomTitleInURL {
		// Maybe the title has changed recently,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := NiceIdiomURL(idiom)
		http.Redirect(w, r, properURL, 301)
		return nil
	}

	selectedImplID := 0
	selectedImplLang := ""
	if selectedImplIDStr := vars["implId"]; len(selectedImplIDStr) > 0 {
		selectedImplID = String2Int(selectedImplIDStr)
		for _, impl := range idiom.Implementations {
			if selectedImplID == impl.Id {
				selectedImplLang = impl.LanguageName
				break
			}
		}
		if selectedImplLang == "" {
			// The requested implementation was not found.
			properURL := NiceIdiomURL(idiom)
			http.Redirect(w, r, properURL, 302)
			return nil
		}
	}

	favlangs := lookForFavoriteLanguages(r)
	includeNonFav := seeNonFavorite(r)
	implFavoriteLanguagesFirstWithOrder(idiom, favlangs, selectedImplLang, includeNonFav)

	implLangInURL := vars["implLang"]
	if implLangInURL != "" && strings.ToLower(selectedImplLang) != strings.ToLower(implLangInURL) {
		// Maybe an accident,
		// or someone is attempting a practical joke forging a funny URL ?
		properURL := NiceImplURL(idiom, selectedImplID, selectedImplLang)
		http.Redirect(w, r, properURL, 301)
		return nil
	}

	userProfile := readUserProfile(r)
	daoVotes.decorateIdiom(c, idiom, userProfile.Nickname)
	for i := range idiom.Implementations {
		daoVotes.decorateImpl(c, &(idiom.Implementations[i]), userProfile.Nickname)
	}

	myToggles := copyToggles(toggles)
	myToggles["actionEditIdiom"] = true
	myToggles["actionAddImpl"] = true
	data := &IdiomDetailFacade{
		PageMeta: PageMeta{
			PageTitle: idiom.Title,
			Toggles:   myToggles,
		},
		UserProfile:    userProfile,
		Idiom:          idiom,
		SelectedImplID: selectedImplID,
	}
	return templates.ExecuteTemplate(w, "page-idiom-detail", data)
}
